package service

import (
	"context"
	"dantaautotool/pkg/utils/http"
	"fmt"

	lark "github.com/larksuite/oapi-sdk-go/v3"
	larkbitable "github.com/larksuite/oapi-sdk-go/v3/service/bitable/v1"
	larkdocx "github.com/larksuite/oapi-sdk-go/v3/service/docx/v1"
	"github.com/rs/zerolog/log"
)

// LarkDocServiceIntf defines the interface for LarkDocService.
type LarkDocServiceIntf interface {
	// GetDocumentTitle retrieves the title of a document given its ID.
	// It returns the title as a string and an error if any occurs.
	GetDocumentTitle(documentID string) (string, error)

	// BatchQueryBitableRecords retrieves records from Bitable given a list of record IDs.
	// It returns a slice of pointers to larkbitable.AppTableRecord and an error if any occurs.
	BatchQueryBitableRecords(appToken, tableID string, recordIDs []string) ([]*larkbitable.AppTableRecord, error)

	// AddBitableRecord adds a record to a Bitable.
	AddBitableRecord(appToken, tableID string, fields map[string]interface{}) error
}

// LarkDocService provides methods to interact with Lark documents.
type LarkDocService struct {
	client *lark.Client
}

// NewLarkDocService creates a new instance of LarkDocService.
func NewLarkDocService() *LarkDocService {
	return &LarkDocService{
		client: http.LarkClient,
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

// BatchQueryBitableRecords retrieves records from Bitable given a list of record IDs.
// It returns a slice of pointers to larkbitable.AppTableRecord and an error if any occurs.
func (s *LarkDocService) BatchQueryBitableRecords(appToken, tableID string, recordIDs []string) ([]*larkbitable.AppTableRecord, error) {
	if len(recordIDs) == 0 {
		return nil, nil
	}
	req := larkbitable.NewBatchGetAppTableRecordReqBuilder().
		AppToken(appToken).
		TableId(tableID).
		Body(larkbitable.NewBatchGetAppTableRecordReqBodyBuilder().
			RecordIds(recordIDs).
			UserIdType(`open_id`).
			WithSharedUrl(false).
			AutomaticFields(false).
			Build()).
		Build()

	resp, err := s.client.Bitable.V1.AppTableRecord.BatchGet(context.Background(), req)

	if err != nil {
		log.Error().Err(err).Msg("[BatchQueryBitableRecords] Failed to batch query records")
		return nil, err
	}

	if !resp.Success() {
		log.Error().Msgf("[BatchQueryBitableRecords] Failed to batch query records: %s", resp.Msg)
		return nil, fmt.Errorf("failed to batch query records: %s", resp.Msg)
	}

	return resp.Data.Records, nil
}

// getFieldTextValueFromRecord retrieves the text value of a field from a record.
// It returns the text value as a string and an error if any occurs.
func getFieldTextValueFromRecord(record *larkbitable.AppTableRecord, fieldKey string) (string, error) {
	fieldArray, ok := record.Fields[fieldKey].([]any)
	if !ok {
		log.Error().Msgf("[GetFieldTextValueFromRecord] Invalid format for field key: %s", fieldKey)
		return "", fmt.Errorf("invalid format for field key: %s", fieldKey)
	}
	field, ok := fieldArray[0].(map[string]any)
	if !ok {
		log.Error().Msgf("[GetFieldTextValueFromRecord] Invalid format for field key: %s", fieldKey)
		return "", fmt.Errorf("invalid format for field key: %s", fieldKey)
	}
	fieldValue, ok := field["text"].(string)
	if !ok {
		log.Error().Msgf("[GetFieldTextValueFromRecord] Invalid format for field key: %s", fieldKey)
		return "", fmt.Errorf("invalid format for field key: %s", fieldKey)
	}
	return fieldValue, nil
}

// AddBitableRecord adds a record to a Bitable.
func (s *LarkDocService) AddBitableRecord(appToken, tableID string, fields map[string]interface{}) error {
	req := larkbitable.NewCreateAppTableRecordReqBuilder().
		AppToken(appToken).
		TableId(tableID).
			AppTableRecord(larkbitable.NewAppTableRecordBuilder().
			Fields(fields).
			Build()).
		Build()
	resp, err := s.client.Bitable.V1.AppTableRecord.Create(context.Background(), req)

	if err != nil {
		log.Err(err).Msg("[LarkDocService.AddBitableRecord] Failed to create record")
		return err
	}

	if !resp.Success() {
		log.Error().Msgf("[LarkDocService.AddBitableRecord] Failed to add record: %s", resp.CodeError)
		return fmt.Errorf("failed to add record: %s", resp.CodeError)
	}
	return nil
}
