package listener

import (
	"context"
	"dantaautotool/internal/service"
	"dantaautotool/pkg"
	"fmt"
	"os"

	larkcore "github.com/larksuite/oapi-sdk-go/v3/core"
	"github.com/larksuite/oapi-sdk-go/v3/event/dispatcher"
	larkdrive "github.com/larksuite/oapi-sdk-go/v3/service/drive/v1"
	larkws "github.com/larksuite/oapi-sdk-go/v3/ws"
	"github.com/rs/zerolog/log"
)

// LarkListener listens to Lark events
// See https://open.feishu.cn/document/server-docs/event-subscription-guide/overview for more details
type LarkListener struct {
	// The dispatcher is used to dispatch events to corresponding event handlers
	client *larkws.Client

	// larkDocService is used to interact with Lark documents
	larkDocService service.LarkDocServiceIntf
}

// NewLarkListener creates a new LarkListener
func NewLarkListener() *LarkListener {
	
	return &LarkListener{
		client: nil,
	}
}

func (l *LarkListener) Start() error {
	eventHandler := dispatcher.
		NewEventDispatcher("", ""). // the 2 parameters must be empty strings
		OnP2FileEditV1(func(ctx context.Context, event *larkdrive.P2FileEditV1) error {
			return l.handleDocEditEvent(ctx, event)
		})

	// Create a client
	appID := os.Getenv("LARK_APP_ID")
	appSecret := os.Getenv("LARK_APP_SECRET")
	if appID == "" || appSecret == "" {
		log.Error().Msg("[LarkListener] LARK_APP_ID or LARK_APP_SECRET is empty")
		return fmt.Errorf("LARK_APP_ID or LARK_APP_SECRET is empty")
	}
	cli := larkws.NewClient(
		appID,
		appSecret,
		larkws.WithEventHandler(eventHandler),
		larkws.WithLogLevel(larkcore.LogLevelDebug),
	)

	// Start the client and handle errors properly
	go func() {
		if err := cli.Start(context.Background()); err != nil {
			log.Error().Err(err).Msg("[LarkListener] Failed to start client")
		}
	}()

	return nil
}

func (l *LarkListener) handleDocEditEvent(ctx context.Context, event *larkdrive.P2FileEditV1) error {
	fileToken := event.Event.FileToken
	if fileToken == nil || *fileToken == "" {
		log.Error().Msg("[LarkListener] fileToken is empty")
		return fmt.Errorf("fileToken is empty")
	}
	log.Info().Msgf("[LarkListener] Received doc edit event, fileToken: %s", *fileToken)

	// Match by document title
	// TODO: Maybe we can use a more sophisticated way to match the document @xunzhou
	title, err := l.larkDocService.GetDocumentTitle(*fileToken)
	if err != nil {
		log.Error().Err(err).Msg("[LarkListener] Failed to get document title")
		return err
	}
	log.Info().Msgf("[LarkListener] Document title retrieved: %s", title)
	
	// Apply for banner, send a vote card to Lark
	if title == pkg.LARK_DOC_TITLE_BANNER_QUESTIONAIRE {
	}


	return nil
}

