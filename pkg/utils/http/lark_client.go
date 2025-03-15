package http

import (
	"dantaautotool/config"
	"fmt"

	lark "github.com/larksuite/oapi-sdk-go/v3"
	"github.com/rs/zerolog/log"
)

var (
	// LarkClient is the client used to interact with Lark.
	LarkClient *lark.Client
)

// InitLarkClient initializes the Lark client.
func InitLarkClient() error {
	// Create a client
	appID := config.Config.LarkAppID
	appSecret := config.Config.LarkAppSecret
	if appID == "" || appSecret == "" {
		log.Error().Msg("[InitLarkClient] LARK_APP_ID or LARK_APP_SECRET is empty")
		return fmt.Errorf("LARK_APP_ID or LARK_APP_SECRET is empty")
	}
	LarkClient = lark.NewClient(appID, appSecret)
	log.Info().Msg("[InitLarkClient] Lark client initialized successfully")
	return nil
}
