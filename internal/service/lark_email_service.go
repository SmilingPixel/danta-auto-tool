package service

import (
	"context"
	"dantaautotool/pkg/utils/http"
	"fmt"
	"net/mail"
	"os"

	lark "github.com/larksuite/oapi-sdk-go/v3"
	larkcore "github.com/larksuite/oapi-sdk-go/v3/core"
	"github.com/larksuite/oapi-sdk-go/v3/service/mail/v1"
	"github.com/rs/zerolog/log"
)

// LarkEmailServiceIntf defines the interface for LarkEmailService.
type LarkEmailServiceIntf interface {
	// SendEmailSimple sends an email to the specified email address.
	//
	// Parameters:
	// - me: The sender's email address.
	// - subject: The subject of the email.
	// - toMail: The recipient's email address.
	// - toName: The recipient's name.
	// - meName: The sender's name.
	// - bodyPlainText: The plain text body of the email.
	//
	// Returns:
	// - error: An error if the email could not be sent, otherwise nil.
	SendEmailSimple(me, subject, toMail, toName, meName, bodyPlainText string) error

	// SendEmail sends an email to the specified email address.
	//
	// Parameters:
	// - me: The sender's email address.
	// - subject: The subject of the email.
	// - to: A slice of recipient email addresses.
	// - cc: A slice of CC email addresses.
	// - bcc: A slice of BCC email addresses.
	// - headFrom: The sender's email address to be displayed in the email header.
	// - bodyHtml: The HTML body of the email.
	// - bodyPlainText: The plain text body of the email.
	//
	// Returns:
	// - error: An error if the email could not be sent, otherwise nil.
	//
	// See https://open.feishu.cn/document/server-docs/mail-v1/user_mailbox-message/send for more details.
	SendEmail(me, subject string, to, cc, bcc []*larkmail.MailAddress, headFrom *larkmail.MailAddress, bodyHtml, bodyPlainText string) error
}

// LarkEmailService provides methods to interact with Lark IM.
type LarkEmailService struct {
	client *lark.Client
}

// NewLarkEmailService creates a new instance of LarkEmailService.
func NewLarkEmailService() *LarkEmailService {
	return &LarkEmailService{
		client: http.LarkClient,
	}
}

// SendEmailSimple sends an email to the specified email address.
//
// Parameters:
// - me: The sender's email address.
// - subject: The subject of the email.
// - toMail: The recipient's email address.
// - toName: The recipient's name.
// - meName: The sender's name.
// - bodyPlainText: The plain text body of the email.
//
// Returns:
// - error: An error if the email could not be sent, otherwise nil.
func (s *LarkEmailService) SendEmailSimple(me, subject, toMail, toName, meName, bodyPlainText string) error {
	toMailAddr := larkmail.NewMailAddressBuilder().MailAddress(toMail).Name(toName).Build()
	meMailAddr := larkmail.NewMailAddressBuilder().MailAddress(me).Name(meName).Build()
	return s.SendEmail(me, subject, []*larkmail.MailAddress{toMailAddr}, nil, nil, meMailAddr, "", bodyPlainText)
}

// SendEmail sends an email to the specified email address.
//
// Parameters:
// - me: The sender's email address.
// - subject: The subject of the email.
// - to: A slice of recipient email addresses.
// - cc: A slice of CC email addresses.
// - bcc: A slice of BCC email addresses.
// - headFrom: The sender's email address to be displayed in the email header.
// - bodyHtml: The HTML body of the email.
// - bodyPlainText: The plain text body of the email.
//
// Returns:
// - error: An error if the email could not be sent, otherwise nil.
//
// See https://open.feishu.cn/document/server-docs/mail-v1/user_mailbox-message/send for more details.
func (s *LarkEmailService) SendEmail(me, subject string, to, cc, bcc []*larkmail.MailAddress, headFrom *larkmail.MailAddress, bodyHtml, bodyPlainText string) error {
	userAccessToken := os.Getenv("LARK_USER_ACCESS_TOKEN")
	if userAccessToken == "" {
		log.Warn().Msg("[LarkEmailService.SendEmail] LARK_USER_ACCESS_TOKEN is empty")
	}

	if me == "" {
		log.Warn().Msg("[LarkEmailService.SendEmail] me is empty, fallback to 'me'")
		me = "me"
	}
	msgBuilder := larkmail.NewMessageBuilder()
	msgBuilder.Subject(subject)
	if len(to) > 0 {
		msgBuilder.To(to)
	}
	if len(cc) > 0 {
		msgBuilder.Cc(cc)
	}
	if len(bcc) > 0 {
		msgBuilder.Bcc(bcc)
	}
	msgBuilder.HeadFrom(headFrom)
	msgBuilder.BodyHtml(bodyHtml)
	msgBuilder.BodyPlainText(bodyPlainText)

	req := larkmail.NewSendUserMailboxMessageReqBuilder().
		UserMailboxId(me).
		Message(msgBuilder.Build()).
		Build()

	resp, err := s.client.Mail.V1.UserMailboxMessage.Send(
		context.Background(),
		req,
		larkcore.WithUserAccessToken(userAccessToken),
	)

	if err != nil {
		log.Err(err).Msg("[LarkEmailService.SendEmail] failed to send email")
		return err
	}
	if !resp.Success() {
		log.Error().Msgf("[LarkEmailService.SendEmail] failed to send email: %s", resp.Msg)
		return fmt.Errorf("[LarkEmailService.SendEmail] failed to send email: %s", resp.Msg)
	}
	return nil
}

// validateEmail validates the email address.
func (s *LarkEmailService) validateEmail(email string) bool {
	_, err := mail.ParseAddress(email)
	return err == nil
}
