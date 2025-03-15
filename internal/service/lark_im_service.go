package service

import (
	"context"
	"dantaautotool/pkg/utils/http"

	"github.com/bytedance/sonic"
	lark "github.com/larksuite/oapi-sdk-go/v3"
	"github.com/larksuite/oapi-sdk-go/v3/event/dispatcher/callback"
	larkim "github.com/larksuite/oapi-sdk-go/v3/service/im/v1"
	"github.com/rs/zerolog/log"
)

// LarkIMServiceIntf defines the interface for LarkIMService.
type LarkIMServiceIntf interface {

	// SendCardMessageByTemplate sends a card message to a chat given its ID.
	// It returns an error if any occurs.
	SendCardMessageByTemplate(receiveIdType, receiveID string, templateCardID string, templateVariables map[string]interface{}) error

	// SendMessage sends a message to a chat given its ID.
	// It returns an error if any occurs.
	SendMessage(receiveIdType, receiveID, content string) error
}

// LarkIMService provides methods to interact with Lark IM.
type LarkIMService struct {
	client *lark.Client
}

// NewLarkIMService creates a new instance of LarkIMService.
func NewLarkIMService() *LarkIMService {
	return &LarkIMService{
		client: http.LarkClient,
	}
}

// SendCardMessageByTemplate sends a card message to a chat given its ID.
// It returns an error if any occurs.
func (s *LarkIMService) SendCardMessageByTemplate(receiveIdType, receiveID string, templateCardID string, templateVariables map[string]interface{}) error {
	card := &callback.Card{
		Type: "template",
		Data: &callback.TemplateCard{
			TemplateID:       templateCardID,
			TemplateVariable: templateVariables,
		},
	}

	content, err := sonic.MarshalString(card)
	if err != nil {
		log.Err(err).Msg("[LarkIMService] Failed to marshal card")
		return err
	}

	err = s.SendMessage(receiveIdType, receiveID, content)
	if err != nil {
		log.Err(err).Msg("[LarkIMService] Failed to send card message")
		return err
	}

	return nil
}

// SendMessage sends a message to a chat given its ID.
// It returns an error if any occurs.
// See https://open.feishu.cn/document/server-docs/im-v1/message/create for more details.
func (s *LarkIMService) SendMessage(receiveIdType, receiveID, content string) error {
	resp, err := s.client.Im.Message.Create(context.Background(), larkim.NewCreateMessageReqBuilder().
		ReceiveIdType(receiveIdType).
		Body(larkim.NewCreateMessageReqBodyBuilder().
			MsgType(larkim.MsgTypeInteractive).
			ReceiveId(receiveID).
			Content(content).
			Build()).
		Build())
	if err != nil {
		log.Err(err).Msg("[LarkIMService] Failed to send message")
		return err
	}
	log.Info().Msgf("[LarkIMService] Send message response: %v", resp)
	if !resp.Success() {
		log.Error().Msgf("[LarkIMService] Failed to send message: %s", resp.Error())
		return err
	}
	return nil
}
