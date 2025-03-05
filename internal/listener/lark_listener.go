package listener

import (
	"context"
	"dantaautotool/internal/service"
	"dantaautotool/pkg"
	"fmt"
	"os"

	larkcore "github.com/larksuite/oapi-sdk-go/v3/core"
	"github.com/larksuite/oapi-sdk-go/v3/event/dispatcher"
	"github.com/larksuite/oapi-sdk-go/v3/event/dispatcher/callback"
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

	// larkIMService is used to interact with Lark IM
	larkIMService service.LarkIMServiceIntf

	// dantaService is used to handle business logic related to Danta
	dantaService service.DantaServiceIntf
}

// NewLarkListener creates a new LarkListener
func NewLarkListener(
	larkDocService service.LarkDocServiceIntf,
	larkIMService service.LarkIMServiceIntf,
	dantaService service.DantaServiceIntf,
) *LarkListener {
	return &LarkListener{
		client:         nil,
		larkDocService: larkDocService,
		larkIMService:  larkIMService,
		dantaService:   dantaService,
	}
}

func (l *LarkListener) Start() error {
	eventHandler := dispatcher.
		NewEventDispatcher("", ""). // the 2 parameters must be empty strings
		OnP2FileEditV1(func(ctx context.Context, event *larkdrive.P2FileEditV1) error {
			return l.handleDocEditEvent(ctx, event)
		}).
		OnP2CardActionTrigger(func(ctx context.Context, event *callback.CardActionTriggerEvent) (*callback.CardActionTriggerResponse, error) {
			return l.handleCardActionTriggerEvent(ctx, event)
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

// handleDocEditEvent handles document edit events
func (l *LarkListener) handleDocEditEvent(_ context.Context, event *larkdrive.P2FileEditV1) error {
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
		bannerVoteCardID := os.Getenv("LARK_BANNER_VOTE_CARD_ID")
		if bannerVoteCardID == "" {
			log.Error().Msg("[LarkListener] LARK_BANNER_VOTE_CARD_ID is empty")
			return fmt.Errorf("LARK_BANNER_VOTE_CARD_ID is empty")
		}
		// TODO: extract the banner content from the document @xunzhou24
		// TODO: openID should be the first operator's openID @xunzhou24
		bannerContent := "This is a banner content" // this needs to be dynamically extracted later
		applicantEmail := "example@example.com" // consider making this dynamic as well
		err = l.larkIMService.SendCardMessageByTemplate(*event.Event.OperatorIdList[0].OpenId, bannerVoteCardID, map[string]interface{}{
			"banner_content": bannerContent,
			"applicant_email": applicantEmail,
		})
		if err != nil {
			log.Error().Err(err).Msg("[LarkListener] Failed to send banner vote card")
			return err
		}
		log.Info().Msg("[LarkListener] Banner vote card sent")
	}


	return nil
}


// handleCardActionTriggerEvent handles card action trigger events
// Note: The event value must have a field named "action" to distinguish different buttons
func (l *LarkListener) handleCardActionTriggerEvent(ctx context.Context, event *callback.CardActionTriggerEvent) (*callback.CardActionTriggerResponse, error) {
	// handle card button click callback
	// https://open.feishu.cn/document/uAjLw4CM/ukzMukzMukzM/feishu-cards/card-callback-communication
	log.Info().Msgf("[handleCardActionTriggerEvent], data: %s\n", larkcore.Prettify(event))

	// Use action to distinguish different buttons. You can configure the action of the button in the card building tool.
	action, ok := event.Event.Action.Value["action"].(map[string]string)
	if !ok {
		log.Error().Msg("[handleCardActionTriggerEvent] action not found")
		return nil, fmt.Errorf("action not found")
	}
	// An action is a map with following fields:
	// {
	//   "action": "...",
	//   "banner_content": "...",
	//   "applicant_email": "...",
	// }
	actionType := action["action"]
	bannerContent := action["banner_content"]
	applicantEmail := action["applicant_email"]

	if actionType == pkg.LARK_IM_CARD_ACTION_APPROVE {
		card := callback.CardActionTriggerResponse{
			Toast: &callback.Toast{ 
				Type:    "success",
				Content: "Approved!",
				I18nContent: map[string]string{
					"zh_cn": "已通过",
					"en_us": "Approved!",
				},
			},
			// Card: &callback.Card{
			// 	Type: "template",
			// 	Data: &callback.TemplateCard{
			// 		TemplateID: APPROVED_CARD_ID,
			// 		TemplateVariable: map[string]interface{}{
			// 			"user_ids": []string{event.Event.Operator.OpenID},
			// 			"notes":    event.Event.Action.FormValue["notes_input"],
			// 		},
			// 	},
			// },
		}
		err := l.dantaService.UpdateBannerAndNotify(bannerContent, []string{applicantEmail})
		if err != nil {
			log.Error().Err(err).Msg("[handleCardActionTriggerEvent] Failed to update banner and notify")
			return nil, err
		}
		return &card, nil
	} else if actionType == pkg.LARK_IM_CARD_ACTION_DISAPPROVE {
		card := callback.CardActionTriggerResponse{
			Toast: &callback.Toast{
				Type:    "info",
				Content: "Disapproved!",
				I18nContent: map[string]string{
					"zh_cn": "已驳回",
					"en_us": "Disapproved!",
				},
			},
		}
		return &card, nil
	}

	log.Warn().Msg("[handleCardActionTriggerEvent] Unknown action received")
	return nil, fmt.Errorf("unknown action: %s", action)
}

