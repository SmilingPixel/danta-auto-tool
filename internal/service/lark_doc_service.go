package service

import (
	"context"
	"fmt"

	lark "github.com/larksuite/oapi-sdk-go/v3"
	larkdocx "github.com/larksuite/oapi-sdk-go/v3/service/docx/v1"
	"github.com/rs/zerolog/log"
)

// LarkDocServiceIntf defines the interface for LarkDocService.
type LarkDocServiceIntf interface {
	// GetDocumentTitle retrieves the title of a document given its ID.
	// It returns the title as a string and an error if any occurs.
	GetDocumentTitle(documentID string) (string, error)
}

// LarkDocService provides methods to interact with Lark documents.
type LarkDocService struct {
	client *lark.Client
}

// NewLarkDocService creates a new instance of LarkDocService.
// It takes a Lark client as a parameter and returns a pointer to LarkDocService.
func NewLarkDocService(client *lark.Client) *LarkDocService {
	return &LarkDocService{
		client: client,
	}
}

// getBasicInfo retrieves the basic information of a document given its ID.
// It returns a pointer to larkdocx.GetDocumentRespData and an error if any occurs.
func (s *LarkDocService) getBasicInfo(documentID string) (*larkdocx.GetDocumentRespData, error) {
	// 创建请求对象
	req := larkdocx.NewGetDocumentReqBuilder().
		DocumentId(documentID).
		Build()
	resp, err := s.client.Docx.V1.Document.Get(context.Background(), req)
	if err != nil {
		log.Error().Err(err).Msg("[getBasicInfo] Failed to get document")
		return nil, err
	}
	if !resp.Success() {
		log.Error().Msgf("[getBasicInfo] Failed to get document: %s", resp.Msg)
		return nil, err
	}
	return resp.Data, nil
}

// GetDocumentTitle retrieves the title of a document given its ID.
// It returns the title as a string and an error if any occurs.
func (s *LarkDocService) GetDocumentTitle(documentID string) (string, error) {
	data, err := s.getBasicInfo(documentID)
	if err != nil {
		return "", err
	}
	if data == nil || data.Document == nil || data.Document.Title == nil {
		return "", fmt.Errorf("document data is nil for documentID: %s", documentID)
	}
	return *data.Document.Title, nil
}
