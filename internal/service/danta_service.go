package service

import (
	"dantaautotool/internal/model"

	"github.com/rs/zerolog/log"

	"github.com/larksuite/oapi-sdk-go/v3/service/bitable/v1"
)

// DantaServiceIntf defines the interface for DantaService.
type DantaServiceIntf interface {
	// UpdateBannerAndNotify updates the banner and notifies the applicants.
	UpdateBannerAndNotify(content string, toEmailList []string) error

	// ConvertBitableRecord2Banner converts a BitableRecord to a Banner.
	ConvertBitableRecord2Banner(record *larkbitable.AppTableRecord) *model.Banner
}


// DantaService provides methods to handle business logic related to Danta.
type DantaService struct {
	// larkDocService is used to interact with Lark documents
	larkDocService LarkDocServiceIntf

	// larkEmailService is used to interact with Lark emails
	larkEmailService LarkEmailServiceIntf
}


// NewDantaService creates a new instance of DantaService.
// It takes a LarkDocServiceIntf as a parameter and returns a pointer to DantaService.
func NewDantaService(
	larkDocService LarkDocServiceIntf,
	larkEmailService LarkEmailServiceIntf,
) *DantaService {
	return &DantaService{
		larkDocService: larkDocService,
		larkEmailService: larkEmailService,
	}
}

// UpdateBannerAndNotify do the following things:
//  1. Edit banner config file (in Github repo)
//  2. Send email to applicants
func (s *DantaService) UpdateBannerAndNotify(
	content string,
	toEmailList []string,
) error {
	log.Info().Msgf("[UpdateBannerAndNotify] Start updating banner and notifying applicants, content: %+v, toEmailList: %+v", content, toEmailList)



	for _, email := range toEmailList {
		log.Info().Msgf("[UpdateBannerAndNotify] Sending email to: %s", email)
		// TODO: Send email @xunzhou24
	}

	return nil
}


// ConvertBitableRecord2Banner converts a BitableRecord to a Banner.
// It returns a pointer to Banner.
func (s *DantaService) ConvertBitableRecord2Banner(record *larkbitable.AppTableRecord) *model.Banner {
	// A bitable record structure: https://open.feishu.cn/document/server-docs/docs/bitable-v1/bitable-structure
	tmp := record.Fields["Banner 内容"]
	log.Info().Msgf("[ConvertBitableRecord2Banner] tmp: %+v", tmp)
	log.Info().Msgf("[ConvertBitableRecord2Banner] type of tmp: %T", tmp)
	bannerContentFieldArray, ok := record.Fields["Banner 内容"].([]any)
	if !ok {
		log.Error().Msg("[ConvertBitableRecord2Banner] Invalid format for 'Banner 内容'")
		return nil
	}
	bannerContentField, ok := bannerContentFieldArray[0].(map[string]any)
	if !ok {
		log.Error().Msg("[ConvertBitableRecord2Banner] Invalid format for 'Banner 内容'")
		return nil
	}
	bannerContent, ok := bannerContentField["text"].(string)
	if !ok {
		log.Error().Msg("[ConvertBitableRecord2Banner] Invalid format for 'Banner 内容'")
		return nil
	}

	applicantEmailFieldArray, ok := record.Fields["邮箱"].([]any)
	if !ok {
		log.Error().Msg("[ConvertBitableRecord2Banner] Invalid format for '邮箱'")
		return nil
	}
	applicantEmailField, ok := applicantEmailFieldArray[0].(map[string]any)
	if !ok {
		log.Error().Msg("[ConvertBitableRecord2Banner] Invalid format for '邮箱'")
		return nil
	}
	applicantEmail, ok := applicantEmailField["text"].(string)
	if !ok {
		log.Error().Msg("[ConvertBitableRecord2Banner] Invalid format for '邮箱'")
		return nil
	}
	return &model.Banner{
		Content: bannerContent,
		ApplicantEmail: applicantEmail,
	}
}
