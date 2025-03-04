package service

import (

    lark "github.com/larksuite/oapi-sdk-go/v3"

)

// LarkIMServiceIntf defines the interface for LarkIMService.
type LarkIMServiceIntf interface {
    // SendMessage sends a message to a chat given its ID.
    // It returns an error if any occurs.
    SendMessage(chatID, message string) error
}

// LarkIMService provides methods to interact with Lark IM.
type LarkIMService struct {
    client *lark.Client
}

// NewLarkIMService creates a new instance of LarkIMService.
// It takes a Lark client as a parameter and returns a pointer to LarkIMService.
func NewLarkIMService(client *lark.Client) *LarkIMService {
    return &LarkIMService{
        client: client,
    }
}
