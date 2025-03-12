package listener

import (
	"context"
	"dantaautotool/internal/entity"
	"dantaautotool/internal/service"
	"dantaautotool/pkg"
	"dantaautotool/pkg/utils/http"
	"fmt"
	"os"

	"github.com/bytedance/sonic"
	larkcore "github.com/larksuite/oapi-sdk-go/v3/core"
	"github.com/larksuite/oapi-sdk-go/v3/event/dispatcher"
	"github.com/larksuite/oapi-sdk-go/v3/event/dispatcher/callback"
	larkdrive "github.com/larksuite/oapi-sdk-go/v3/service/drive/v1"
	larkim "github.com/larksuite/oapi-sdk-go/v3/service/im/v1"
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
		OnP2FileBitableRecordChangedV1(func(ctx context.Context, event *larkdrive.P2FileBitableRecordChangedV1) error {
			log.Debug().Msgf("[LarkListener] Received bitable record chanded event: %s", larkcore.Prettify(event))
			return l.handleBitableRecordChangeEvent(ctx, event)
		}).
		OnP2CardActionTrigger(func(ctx context.Context, event *callback.CardActionTriggerEvent) (*callback.CardActionTriggerResponse, error) {
			log.Debug().Msgf("[LarkListener] Received card action trigger event: %s", larkcore.Prettify(event))
			return l.handleCardActionTriggerEvent(ctx, event)
		}).
		OnP2MessageReceiveV1(func(ctx context.Context, event *larkim.P2MessageReceiveV1) error {
			log.Debug().Msgf("[LarkListener] Received message receive event: %s", larkcore.Prettify(event))
			return l.handleMessageReceiveEvent(ctx, event)
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

// handleBitableRecordChangeEvent handles bitable record changed events
func (l *LarkListener) handleBitableRecordChangeEvent(_ context.Context, event *larkdrive.P2FileBitableRecordChangedV1) error {
	fileToken := event.Event.FileToken
	if fileToken == nil || *fileToken == "" {
		log.Error().Msg("[LarkListener.handleBitableRecordChangeEvent] fileToken is empty")
		return fmt.Errorf("fileToken is empty")
	}
	log.Info().Msgf("[LarkListener.handleBitableRecordChangeEvent] Received bitable record changed event, fileToken: %s", *fileToken)

	// Match by file token
	bannerAnalysisDocToken := os.Getenv("LARK_BANNER_BITABLE_APP_TOKEN")
	bannerAnalysisTableID := os.Getenv("LARK_BANNER_BITABLE_TABLE_ID")
	if bannerAnalysisDocToken == "" || bannerAnalysisTableID == "" {
		log.Error().Msg("[LarkListener.handleBitableRecordChangeEvent] LARK_BANNER_BITABLE_APP_TOKEN or LARK_BANNER_BITABLE_TABLE_ID is empty")
		return fmt.Errorf("LARK_BANNER_BITABLE_APP_TOKEN or LARK_BANNER_BITABLE_TABLE_ID is empty")
	}
	bannerVoteCardID := os.Getenv("LARK_BANNER_APPROVE_CARD_ID")
	if bannerVoteCardID == "" {
		log.Error().Msg("[LarkListener.handleBitableRecordChangeEvent] LARK_BANNER_APPROVE_CARD_ID is empty")
		return fmt.Errorf("LARK_BANNER_APPROVE_CARD_ID is empty")
	}
	if *fileToken == bannerAnalysisDocToken {
		addedRecordIds := make([]string, 0)
		for _, action := range event.Event.ActionList {
			// Only handle add action
			if *action.Action != pkg.LARK_BITABLE_RECORD_ACTION_ADD {
				continue
			}
			addedRecordIds = append(addedRecordIds, *action.RecordId)
		}
		if len(addedRecordIds) == 0 {
			log.Info().Msg("[LarkListener.handleBitableRecordChangeEvent] No added record found")
			return nil
		}
		// Batch query bitable records
		addedRecords, err := l.larkDocService.BatchQueryBitableRecords(bannerAnalysisDocToken, bannerAnalysisTableID, addedRecordIds)
		if err != nil {
			log.Error().Err(err).Msg("[LarkListener.handleBitableRecordChangeEvent] Failed to batch query bitable records")
			return err
		}
		// For each added record, send a banner vote card
		for _, addedRecord := range addedRecords {
			bannerApplication := l.dantaService.ConvertBitableRecord2BannerApplication(addedRecord)
			if bannerApplication == nil {
				log.Error().Msg("[LarkListener.handleBitableRecordChangeEvent] Failed to convert bitable record to banner application")
				return fmt.Errorf("failed to convert bitable record to banner application")
			}
			err = l.larkIMService.SendCardMessageByTemplate(*event.Event.OperatorId.OpenId, bannerVoteCardID, map[string]interface{}{
				"banner_title":  bannerApplication.Title,
				"banner_action": bannerApplication.Action,
				"banner_button": bannerApplication.Button,
				"applicant_email": bannerApplication.ApplicantEmail,
			})
			if err != nil {
				log.Error().Err(err).Msg("[LarkListener] Failed to send banner vote card")
				return err
			}
			log.Info().Msg("[LarkListener] Banner vote card sent")
		}
	}
	return nil
}

// handleCardActionTriggerEvent handles card action trigger events
// Note: The event value must have a field named "action" to distinguish different buttons
func (l *LarkListener) handleCardActionTriggerEvent(_ context.Context, event *callback.CardActionTriggerEvent) (*callback.CardActionTriggerResponse, error) {
	// handle card button click callback
	// https://open.feishu.cn/document/uAjLw4CM/ukzMukzMukzM/feishu-cards/card-callback-communication
	log.Info().Msgf("[handleCardActionTriggerEvent], data: %s\n", larkcore.Prettify(event))

	// Use action to distinguish different buttons. You can configure the action of the button in the card building tool.
	actionDetail := event.Event.Action.Value
	// An action is a map with following fields:
	//  {
	//    "action": "...",
	//    "banner_content": "...",
	//    "applicant_email": "...",
	//  }
	var ok bool
	actionType, ok := actionDetail["action"].(string)
	if !ok {
		log.Error().Msgf("[handleCardActionTriggerEvent] Failed to parse action type, actionDetail: %v", actionDetail)
		return nil, fmt.Errorf("failed to parse action")
	}
	bannerTitle, ok := actionDetail["banner_title"].(string)
	if !ok {
		log.Error().Msgf("[handleCardActionTriggerEvent] Failed to parse banner content, actionDetail: %v", actionDetail)
		return nil, fmt.Errorf("failed to parse action")
	}
	bannerAction, ok := actionDetail["banner_action"].(string)
	if !ok {
		log.Error().Msgf("[handleCardActionTriggerEvent] Failed to parse action, actionDetail: %v", actionDetail)
		return nil, fmt.Errorf("failed to parse action")
	}
	bannerButton, ok := actionDetail["banner_button"].(string)
	if !ok {
		log.Error().Msgf("[handleCardActionTriggerEvent] Failed to parse action, actionDetail: %v", actionDetail)
		return nil, fmt.Errorf("failed to parse action")
	}
	applicantEmail, ok := actionDetail["applicant_email"].(string)
	if !ok {
		log.Error().Msgf("[handleCardActionTriggerEvent] Failed to parse action email, actionDetail: %v", actionDetail)
		return nil, fmt.Errorf("failed to parse action")
	}

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
		newBanner := entity.Banner{
			Title:  bannerTitle,
			Action: bannerAction,
			Button: bannerButton,
		}
		err := l.dantaService.UpdateBannerAndNotify(newBanner, []string{applicantEmail})
		if err != nil {
			log.Error().Err(err).Msg("[handleCardActionTriggerEvent] Failed to update banner and notify")
			return nil, err
		}
		log.Info().Msg("[handleCardActionTriggerEvent] Banner updated and notification sent")
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
	return nil, fmt.Errorf("unknown action type: %s", actionType)
}

func (l *LarkListener) handleMessageReceiveEvent(_ context.Context, event *larkim.P2MessageReceiveV1) error {
	fmt.Printf("[OnP2MessageReceiveV1 access], data: %s\n", larkcore.Prettify(event))
	/**
	* 解析用户发送的消息。
	* Parse the message sent by the user.
	 */
	var respContent map[string]string
	err := sonic.Unmarshal([]byte(*event.Event.Message.Content), &respContent)
	/**
	* 检查消息类型是否为文本
	* Check if the message type is text
	 */
	if err != nil || *event.Event.Message.MessageType != "text" {
		respContent = map[string]string{
			"text": "解析消息失败，请发送文本消息\nparse message failed, please send text message",
		}
	}

	/**
	* 构建回复消息
	* Build reply message
	 */
	content := larkim.NewTextMsgBuilder().
		TextLine("收到你发送的消息: " + respContent["text"]).
		TextLine("Received message: " + respContent["text"]).
		Build()

	if *event.Event.Message.ChatType == "p2p" {
		/**
		* 使用SDK调用发送消息接口。 Use SDK to call send message interface.
		* https://open.feishu.cn/document/uAjLw4CM/ukTMukTMukTM/reference/im-v1/message/create
		 */
		resp, err := http.LarkClient.Im.Message.Create(context.Background(), larkim.NewCreateMessageReqBuilder().
			ReceiveIdType(larkim.ReceiveIdTypeChatId). // 消息接收者的 ID 类型，设置为会话ID。 ID type of the message receiver, set to chat ID.
			Body(larkim.NewCreateMessageReqBodyBuilder().
				MsgType(larkim.MsgTypeText).            // 设置消息类型为文本消息。 Set message type to text message.
				ReceiveId(*event.Event.Message.ChatId). // 消息接收者的 ID 为消息发送的会话ID。 ID of the message receiver is the chat ID of the message sending.
				Content(content).
				Build()).
			Build())

		if err != nil || !resp.Success() {
			fmt.Println(err)
			fmt.Println(resp.Code, resp.Msg, resp.RequestId())
			return nil
		}

	} else {
		/**
		* 使用SDK调用回复消息接口。 Use SDK to call send message interface.
		* https://open.feishu.cn/document/server-docs/im-v1/message/reply
		 */
		resp, err := http.LarkClient.Im.Message.Reply(context.Background(), larkim.NewReplyMessageReqBuilder().
			MessageId(*event.Event.Message.MessageId).
			Body(larkim.NewReplyMessageReqBodyBuilder().
				MsgType(larkim.MsgTypeText). // 设置消息类型为文本消息。 Set message type to text message.
				Content(content).
				Build()).
			Build())
		if err != nil || !resp.Success() {
			fmt.Printf("logId: %s, error response: \n%s", resp.RequestId(), larkcore.Prettify(resp.CodeError))
			return nil
		}
	}

	return nil
}
