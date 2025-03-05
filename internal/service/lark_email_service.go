package service

import (
	"dantaautotool/pkg/utils/http"

	lark "github.com/larksuite/oapi-sdk-go/v3"
)

// LarkEmailServiceIntf defines the interface for LarkEmailService.
type LarkEmailServiceIntf interface {
    
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
