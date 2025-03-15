package config

import (
	"os"

	"github.com/rs/zerolog/log"
)

type GlobalConfig struct {

	// 飞书应用的 APP ID 和 APP Secret
    LarkAppID                       string
    LarkAppSecret                   string

	// 飞书应用中的审批卡片 ID
    LarkBannerApproveCardID         string

	// Banner 宣传位的多维表格的 APP Token 和 Table ID（包括申请表和使用记录表）
    LarkBannerBitableAppToken       string
    LarkBannerBitableApplicationTableID string
    LarkBannerBitableUsageTableID   string

	// Banner 宣传位的审批群 ID
    LarkBannerApproveGroupID        string

	// Danta 开发者邮箱（暂时没有用到）
    DantaDevEmail                   string

	// Github 个人访问令牌
    GithubPersonalAccessToken       string

	// Github 仓库的 owner、name 和 Banner 配置文件路径
    GithubDanxiRepoOwner            string
    GithubDanxiRepoName             string
    GithubDanxiRepoAppConfigPath    string
}

var Config GlobalConfig

// LoadConfig loads the configuration from environment variables.
func LoadConfig() {
    Config = GlobalConfig{
        LarkAppID:                       os.Getenv("LARK_APP_ID"),
        LarkAppSecret:                   os.Getenv("LARK_APP_SECRET"),
        LarkBannerApproveCardID:         os.Getenv("LARK_BANNER_APPROVE_CARD_ID"),
        LarkBannerBitableAppToken:       os.Getenv("LARK_BANNER_BITABLE_APP_TOKEN"),
        LarkBannerBitableApplicationTableID: os.Getenv("LARK_BANNER_BITABLE_APPLICATION_TABLE_ID"),
        LarkBannerBitableUsageTableID:   os.Getenv("LARK_BANNER_BITABLE_USAGE_TABLE_ID"),
        LarkBannerApproveGroupID:        os.Getenv("LARK_BANNER_APPROVE_GROUP_ID"),
        DantaDevEmail:                   os.Getenv("DANTA_DEV_EMAIL"),
        GithubPersonalAccessToken:       os.Getenv("GITHUB_PERSONAL_ACCESS_TOKEN"),
        GithubDanxiRepoOwner:            os.Getenv("GITHUB_DANXI_REPO_OWNER"),
        GithubDanxiRepoName:             os.Getenv("GITHUB_DANXI_REPO_NAME"),
        GithubDanxiRepoAppConfigPath:    os.Getenv("GITHUB_DANXI_REPO_APP_CONFIG_PATH"),
    }

	// Check if any of the required environment variables are missing
	if Config.LarkAppID == "" {
		log.Error().Msg("LARK_APP_ID is empty")
	}
	if Config.LarkAppSecret == "" {
		log.Error().Msg("LARK_APP_SECRET is empty")
	}
	if Config.LarkBannerApproveCardID == "" {
		log.Error().Msg("LARK_BANNER_APPROVE_CARD_ID is empty")
	}
	if Config.LarkBannerBitableAppToken == "" {
		log.Error().Msg("LARK_BANNER_BITABLE_APP_TOKEN is empty")
	}
	if Config.LarkBannerBitableApplicationTableID == "" {
		log.Error().Msg("LARK_BANNER_BITABLE_APPLICATION_TABLE_ID is empty")
	}
	if Config.LarkBannerBitableUsageTableID == "" {
		log.Error().Msg("LARK_BANNER_BITABLE_USAGE_TABLE_ID is empty")
	}
	if Config.LarkBannerApproveGroupID == "" {
		log.Error().Msg("LARK_BANNER_APPROVE_GROUP_ID is empty")
	}
	if Config.DantaDevEmail == "" {
		log.Error().Msg("DANTA_DEV_EMAIL is empty")
	}
	if Config.GithubPersonalAccessToken == "" {
		log.Error().Msg("GITHUB_PERSONAL_ACCESS_TOKEN is empty")
	}
	if Config.GithubDanxiRepoOwner == "" {
		log.Error().Msg("GITHUB_DANXI_REPO_OWNER is empty")
	}
	if Config.GithubDanxiRepoName == "" {
		log.Error().Msg("GITHUB_DANXI_REPO_NAME is empty")
	}
	if Config.GithubDanxiRepoAppConfigPath == "" {
		log.Error().Msg("GITHUB_DANXI_REPO_APP_CONFIG_PATH is empty")
	}
}