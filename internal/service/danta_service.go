package service

import (
	"dantaautotool/config"
	"dantaautotool/internal/entity"
	"fmt"

	"github.com/pelletier/go-toml/v2"
	"github.com/rs/zerolog/log"

	larkbitable "github.com/larksuite/oapi-sdk-go/v3/service/bitable/v1"
)

// DantaServiceIntf defines the interface for DantaService.
type DantaServiceIntf interface {
	// UpdateBannerAndNotify updates the banner and notifies the applicants.
	UpdateBannerAndNotify(newBanner entity.Banner, toEmailList []string) error

	// UpdateBanner edits banner config file (in Github repo)
	UpdateBanner(newBanner entity.Banner) error

	// NotifyBannerUpdate send email to applicants when banner is updated
	NotifyBannerUpdate(newBanner entity.Banner, toEmailList []string) error

	// ConvertBitableRecord2BannerApplication converts a BitableRecord to a Banner.
	ConvertBitableRecord2BannerApplication(record *larkbitable.AppTableRecord) *entity.BannerApplication
}

// DantaService provides methods to handle business logic related to Danta.
type DantaService struct {
	// larkDocService is used to interact with Lark documents
	larkDocService LarkDocServiceIntf

	// larkEmailService is used to interact with Lark emails
	larkEmailService LarkEmailServiceIntf

	// githubService is used to interact with Github
	githubService GithubServiceIntf
}

// NewDantaService creates a new instance of DantaService.
// It takes a LarkDocServiceIntf as a parameter and returns a pointer to DantaService.
func NewDantaService(
	larkDocService LarkDocServiceIntf,
	larkEmailService LarkEmailServiceIntf,
	githubService GithubServiceIntf,
) *DantaService {
	return &DantaService{
		larkDocService:   larkDocService,
		larkEmailService: larkEmailService,
		githubService:    githubService,
	}
}

// UpdateBannerAndNotify do the following things:
//  1. Edit banner config file (in Github repo)
//  2. Send email to applicants
func (s *DantaService) UpdateBannerAndNotify(newBanner entity.Banner, toEmailList []string) error {
	err := s.UpdateBanner(newBanner)
	if err != nil {
		log.Err(err).Msg("[DantaService.UpdateBannerAndNotify] Failed to update banner")
		return err
	}

	err = s.NotifyBannerUpdate(newBanner, toEmailList)
	if err != nil {
		log.Err(err).Msg("[DantaService.UpdateBannerAndNotify] Failed to notify applicants")
		return err
	}

	return nil
}

// UpdateBanner edits banner config file (in Github repo)
func (s *DantaService) UpdateBanner(newBanner entity.Banner) error {
	log.Info().Msgf("[DantaService.UpdateBanner] Start updating banner and notifying applicants, newBanner: %+v", newBanner)

	bannerRepoOwner := config.Config.GithubDanxiRepoOwner
	if bannerRepoOwner == "" {
		log.Error().Msg("[DantaService.UpdateBanner] GITHUB_DANXI_REPO_OWNER is empty")
		return fmt.Errorf("GITHUB_DANXI_REPO_OWNER is empty")
	}
	bannerRepoName := config.Config.GithubDanxiRepoName
	if bannerRepoName == "" {
		log.Error().Msg("[DantaService.UpdateBanner] GITHUB_DANXI_REPO_NAME is empty")
		return fmt.Errorf("GITHUB_DANXI_REPO_NAME is empty")
	}
	bannerRepoAppConfigPath := config.Config.GithubDanxiRepoAppConfigPath
	if bannerRepoAppConfigPath == "" {
		log.Error().Msg("[DantaService.UpdateBanner] GITHUB_DANXI_REPO_APP_CONFIG_PATH is empty")
		return fmt.Errorf("GITHUB_DANXI_REPO_APP_CONFIG_PATH is empty")
	}

	repoContent, err := s.githubService.GetFileContent(
		bannerRepoOwner,
		bannerRepoName,
		bannerRepoAppConfigPath,
	)
	if err != nil {
		log.Err(err).Msg("[DantaService.UpdateBanner] Failed to get banner config file content")
		return err
	}

	// for the file structure, see:
	// https://github.com/SmilingPixel/DanXi-Backend/blob/main/public/tmp_wait_for_json_editor.toml
	configContent := repoContent.DecodedContent
	sha := repoContent.SHA

	// parse and update the config content
	dantaAppContentConfig := entity.DantaAppContentConfig{}
	err = toml.Unmarshal([]byte(configContent), &dantaAppContentConfig)
	if err != nil {
		log.Err(err).Msg("[DantaService.UpdateBanner] Failed to unmarshal config content")
		return err
	}
	dantaAppContentConfig.Banners = append(dantaAppContentConfig.Banners, newBanner)
	updatedConfigContentBytes, err := toml.Marshal(dantaAppContentConfig)
	if err != nil {
		log.Err(err).Msg("[DantaService.UpdateBanner] Failed to marshal updated config content")
		return err
	}
	updatedConfigContent := string(updatedConfigContentBytes)

	// update file in Github
	err = s.githubService.CreateOrUpdateFileContent(
		bannerRepoOwner,
		bannerRepoName,
		bannerRepoAppConfigPath,
		"commit message", // TODO: use a meaningful commit message @xunzhou24
		updatedConfigContent,
		sha,
		"main", // commit to default branch
		&entity.Committer{
			Name:  "Danta Auto Tool",
			Email: "danta@example.com",
		}, // TODO: use the real committer @xunzhou24
	)
	if err != nil {
		log.Err(err).Msg("[DantaService.UpdateBanner] Failed to update file content in Github")
		return err
	}

	return nil
}

// NotifyBannerUpdate send email to applicants when banner is updated
func (s *DantaService) NotifyBannerUpdate(newBanner entity.Banner, toEmailList []string) error {

	for _, email := range toEmailList {
		log.Info().Msgf("[DantaService.NotifyBannerUpdate] Sending email to: %s", email)
		err := s.larkEmailService.SendEmailSimple(
			config.Config.DantaDevEmail,
			"您提交的置顶申请已经通过",
			email,
			"测试toname",
			"旦挞的小菜鸡周迅",
			"您提交的置顶申请已经通过",
		)
		if err != nil {
			log.Error().Err(err).Msg("[DantaService.NotifyBannerUpdate] Failed to send email")
			return err
		}
		log.Info().Msgf("[DantaService.NotifyBannerUpdate] Email sent to: %s", email)
	}

	return nil
}

// ConvertBitableRecord2BannerApplication converts a BitableRecord to a Banner application.
// It returns a pointer to Banner.
func (s *DantaService) ConvertBitableRecord2BannerApplication(record *larkbitable.AppTableRecord) *entity.BannerApplication {
	// A bitable record structure: https://open.feishu.cn/document/server-docs/docs/bitable-v1/bitable-structure
	bannerTitle, err := getFieldTextValueFromRecord(record, "标题")
	if err != nil {
		log.Err(err).Msg("[DantaService.ConvertBitableRecord2BannerApplication] Failed to get banner title")
		return nil
	}
	bannerAction, err := getFieldTextValueFromRecord(record, "操作")
	if err != nil {
		log.Err(err).Msg("[DantaService.ConvertBitableRecord2BannerApplication] Failed to get banner action")
		return nil
	}
	bannerButton, err := getFieldTextValueFromRecord(record, "操作提示")
	if err != nil {
		log.Err(err).Msg("[DantaService.ConvertBitableRecord2BannerApplication] Failed to get banner button")
		return nil
	}
	applicantEmail, err := getFieldTextValueFromRecord(record, "邮箱")
	if err != nil {
		log.Err(err).Msg("[DantaService.ConvertBitableRecord2BannerApplication] Failed to get applicant email")
		return nil
	}
	return &entity.BannerApplication{
		Banner: entity.Banner{
			Title:  bannerTitle,
			Action: bannerAction,
			Button: bannerButton,
		},
		ApplicantEmail: applicantEmail,
	}
}
