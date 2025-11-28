package service

import (
	"context"
	"encoding/json"

	"git.uhomes.net/uhs-go/go-bisub/internal/models"
	"git.uhomes.net/uhs-go/go-bisub/internal/repository"
)

type OperationLogService struct {
	repo *repository.OperationLogRepository
}

func NewOperationLogService(repo *repository.OperationLogRepository) *OperationLogService {
	return &OperationLogService{repo: repo}
}

// LogOperation 记录操作日志
func (s *OperationLogService) LogOperation(ctx context.Context, log *models.OperationLog) {
	// 异步记录日志，不影响主业务流程
	go func() {
		if err := s.repo.Create(context.Background(), log); err != nil {
			// 记录到系统日志，但不抛出错误
			// TODO: 使用结构化日志记录
		}
	}()
}

// CreateOperationLog 创建操作日志
func (s *OperationLogService) CreateOperationLog(userID uint64, username, operation, resource, resourceID, status, clientIP, userAgent, requestURL, method string, duration uint32, errorMsg string, requestData, responseData interface{}) *models.OperationLog {
	var reqData, respData json.RawMessage

	if requestData != nil {
		if data, err := json.Marshal(requestData); err == nil {
			reqData = data
		}
	}

	if responseData != nil {
		if data, err := json.Marshal(responseData); err == nil {
			respData = data
		}
	}

	return &models.OperationLog{
		UserID:       userID,
		Username:     username,
		Operation:    operation,
		Resource:     resource,
		ResourceID:   resourceID,
		Status:       status,
		ClientIP:     clientIP,
		UserAgent:    userAgent,
		RequestURL:   requestURL,
		Method:       method,
		Duration:     duration,
		ErrorMsg:     errorMsg,
		RequestData:  reqData,
		ResponseData: respData,
	}
}

// GetOperationLogs 获取操作日志列表
func (s *OperationLogService) GetOperationLogs(ctx context.Context, req *models.OperationLogRequest) ([]*models.OperationLog, int64, error) {
	return s.repo.List(ctx, req)
}
