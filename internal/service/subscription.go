package service

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"regexp"
	"strings"
	"time"

	"git.uhomes.net/uhs-go/go-bisub/internal/config"
	"git.uhomes.net/uhs-go/go-bisub/internal/models"
	"git.uhomes.net/uhs-go/go-bisub/internal/repository"
	"gorm.io/gorm"
)

type SubscriptionService struct {
	repo        *repository.SubscriptionRepository
	statsRepo   *repository.StatsRepository
	dataSources map[string]*gorm.DB
	config      *config.Config
}

func NewSubscriptionService(repo *repository.SubscriptionRepository, statsRepo *repository.StatsRepository, dataSources map[string]*gorm.DB, cfg *config.Config) *SubscriptionService {
	return &SubscriptionService{
		repo:        repo,
		statsRepo:   statsRepo,
		dataSources: dataSources,
		config:      cfg,
	}
}

func (s *SubscriptionService) CreateSubscription(ctx context.Context, req *models.CreateSubscriptionRequest, creatorID uint64) (*models.Subscription, error) {
	// 解析并校验extra_config
	var extraConfig models.ExtraConfig
	if err := json.Unmarshal(req.ExtraConfig, &extraConfig); err != nil {
		return nil, fmt.Errorf("invalid extra_config: %w", err)
	}

	// SQL安全校验
	if err := s.validateSQL(extraConfig.SQLContent); err != nil {
		return nil, fmt.Errorf("SQL validation failed: %w", err)
	}

	subscription := &models.Subscription{
		Type:        req.Type,
		SubKey:      req.SubKey,
		Version:     req.Version,
		Title:       req.Title,
		Abstract:    req.Abstract,
		Status:      req.Status,
		CreatedBy:   creatorID,
		ExtraConfig: req.ExtraConfig,
	}

	if err := s.repo.Create(ctx, subscription); err != nil {
		return nil, err
	}

	return subscription, nil
}

func (s *SubscriptionService) ExecuteSubscription(ctx context.Context, subType, key string, version *uint8, req *models.ExecuteSubscriptionRequest, clientIP, apiURL string) ([]map[string]interface{}, error) {
	// 获取订阅
	var subscription *models.Subscription
	var err error

	if version != nil {
		subscription, err = s.repo.GetByKeyAndVersion(ctx, subType, key, *version)
		if err != nil || subscription.Status == models.StatusExpired {
			// 如果指定版本不存在或已失效，获取活跃的最高版本
			subscription, err = s.repo.GetActiveByKey(ctx, subType, key)
		}
	} else {
		subscription, err = s.repo.GetActiveByKey(ctx, subType, key)
	}

	if err != nil {
		return nil, fmt.Errorf("subscription not found: %w", err)
	}

	// 解析extra_config
	var extraConfig models.ExtraConfig
	if err := json.Unmarshal(subscription.ExtraConfig, &extraConfig); err != nil {
		return nil, fmt.Errorf("invalid extra_config: %w", err)
	}

	// 替换SQL变量
	executedSQL, err := s.replaceVariables(extraConfig.SQLContent, req.Variables, extraConfig.SQLReplace)
	if err != nil {
		return nil, err
	}

	// 选择数据源
	dataSource := req.DataSource
	if dataSource == "" {
		dataSource = "default"
	}

	db, exists := s.dataSources[dataSource]
	if !exists {
		return nil, fmt.Errorf("data source %s not found", dataSource)
	}

	// 设置超时
	timeout := time.Duration(req.Timeout) * time.Millisecond
	if timeout == 0 {
		timeout = s.config.Server.Timeout
	}

	execCtx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	// 执行SQL
	startTime := time.Now()
	rows, err := db.WithContext(execCtx).Raw(executedSQL).Rows()
	if err != nil {
		return nil, fmt.Errorf("SQL execution failed: %w", err)
	}
	defer rows.Close()

	// 处理结果
	results, err := s.processRows(rows)
	if err != nil {
		return nil, err
	}

	executionDuration := uint32(time.Since(startTime).Milliseconds())

	// 异步记录统计
	requestResponse := models.RequestResponse{
		Params:         req.Variables,
		InstanceSQL:    executedSQL,
		InstanceSource: dataSource,
		RequestIP:      clientIP,
		Version:        subscription.Version,
	}
	requestResponseJSON, _ := json.Marshal(requestResponse)

	go s.recordStats(context.Background(), &models.SubscriptionStats{
		SubKey:            subscription.SubKey,
		Version:           subscription.Version,
		ExecutionDuration: executionDuration,
		RequestURL:        apiURL,
		RequestResponse:   requestResponseJSON,
		InstanceSource:    dataSource,
	})

	return results, nil
}

func (s *SubscriptionService) validateSQL(sqlContent string) error {
	// 移除注释和多余空格
	cleaned := strings.TrimSpace(regexp.MustCompile(`--.*|/\*[\s\S]*?\*/`).ReplaceAllString(sqlContent, ""))

	// 检查是否为允许的SQL类型
	for _, allowedType := range s.config.Security.AllowedSQLTypes {
		if strings.HasPrefix(strings.ToUpper(cleaned), strings.ToUpper(allowedType)) {
			return nil
		}
	}

	return fmt.Errorf("SQL type not allowed, only %v are permitted", s.config.Security.AllowedSQLTypes)
}

func (s *SubscriptionService) replaceVariables(sqlContent string, variables map[string]interface{}, sqlReplace map[string]string) (string, error) {

	result := sqlContent

	// 查找所有变量占位符
	re := regexp.MustCompile(`(\w+_replace)`)
	matches := re.FindAllString(sqlContent, -1)

	for _, match := range matches {
		value, exists := variables[match]
		if !exists {
			return "", fmt.Errorf("missing required variable: %s", match)
		}

		// 简单的SQL注入防护
		valueStr := fmt.Sprintf("%v", value)
		if strings.Contains(valueStr, "'") || strings.Contains(valueStr, ";") || strings.Contains(valueStr, "--") {
			return "", fmt.Errorf("invalid variable value: %s", match)
		}

		result = strings.ReplaceAll(result, match, valueStr)
	}

	return result, nil
}

func (s *SubscriptionService) processRows(rows *sql.Rows) ([]map[string]interface{}, error) {
	columns, err := rows.Columns()
	if err != nil {
		return nil, err
	}

	var results []map[string]interface{}

	for rows.Next() {
		values := make([]interface{}, len(columns))
		valuePtrs := make([]interface{}, len(columns))

		for i := range values {
			valuePtrs[i] = &values[i]
		}

		if err := rows.Scan(valuePtrs...); err != nil {
			return nil, err
		}

		row := make(map[string]interface{})
		for i, col := range columns {
			val := values[i]
			if b, ok := val.([]byte); ok {
				row[col] = string(b)
			} else {
				row[col] = val
			}
		}

		results = append(results, row)
	}

	return results, rows.Err()
}

func (s *SubscriptionService) recordStats(ctx context.Context, stats *models.SubscriptionStats) {
	if err := s.statsRepo.Create(ctx, stats); err != nil {
		// 记录日志但不影响主流程
		fmt.Printf("Failed to record stats: %v\n", err)
	}
}

func (s *SubscriptionService) marshalRequestParams(req *models.ExecuteSubscriptionRequest) json.RawMessage {
	data, _ := json.Marshal(req)
	return data
}

func (s *SubscriptionService) GetStats(ctx context.Context, req *models.StatsQueryRequest) ([]*models.StatsResponse, error) {
	// 默认查询最近7天
	endTime := time.Now()
	startTime := endTime.AddDate(0, 0, -7)

	if req.StartTime != "" {
		if t, err := time.Parse("2006-01-02", req.StartTime); err == nil {
			startTime = t
		}
	}

	if req.EndTime != "" {
		if t, err := time.Parse("2006-01-02", req.EndTime); err == nil {
			endTime = t.Add(24 * time.Hour) // 包含结束日期的全天
		}
	}

	limit := req.Limit
	if limit <= 0 || limit > 100 {
		limit = 20
	}

	offset := req.Offset
	if offset < 0 {
		offset = 0
	}

	return s.statsRepo.GetStats(ctx, startTime, endTime, limit, offset)
}

func (s *SubscriptionService) GetSubscriptions(ctx context.Context, limit, offset int, subKey, title, status string) ([]*models.Subscription, int64, error) {
	if limit <= 0 || limit > 100 {
		limit = 20
	}
	if offset < 0 {
		offset = 0
	}

	return s.repo.List(ctx, limit, offset, subKey, title, status)
}

func (s *SubscriptionService) GetSubscription(ctx context.Context, subType, key string, version *uint8) (*models.Subscription, error) {
	if version != nil {
		return s.repo.GetByKeyAndVersion(ctx, subType, key, *version)
	}
	return s.repo.GetActiveByKey(ctx, subType, key)
}

func (s *SubscriptionService) UpdateSubscription(ctx context.Context, subType, key string, version uint8, req *models.UpdateSubscriptionRequest) (*models.Subscription, error) {
	// 获取现有订阅
	subscription, err := s.repo.GetByKeyAndVersion(ctx, subType, key, version)
	if err != nil {
		return nil, fmt.Errorf("subscription not found: %w", err)
	}

	// 如果更新了extra_config，需要验证SQL
	if len(req.ExtraConfig) > 0 {
		var extraConfig models.ExtraConfig
		if err := json.Unmarshal(req.ExtraConfig, &extraConfig); err != nil {
			return nil, fmt.Errorf("invalid extra_config: %w", err)
		}
		if err := s.validateSQL(extraConfig.SQLContent); err != nil {
			return nil, fmt.Errorf("SQL validation failed: %w", err)
		}
		subscription.ExtraConfig = req.ExtraConfig
	}

	// 更新字段
	if req.Title != "" {
		subscription.Title = req.Title
	}
	if req.Abstract != "" {
		subscription.Abstract = req.Abstract
	}
	if req.Status != "" {
		subscription.Status = req.Status
	}

	if err := s.repo.Update(ctx, subscription); err != nil {
		return nil, err
	}

	return subscription, nil
}

func (s *SubscriptionService) UpdateStatus(ctx context.Context, subType, key string, version uint8, status string) error {
	return s.repo.UpdateStatus(ctx, subType, key, version, status)
}

func (s *SubscriptionService) DeleteSubscription(ctx context.Context, subType, key string, version uint8) error {
	return s.repo.Delete(ctx, subType, key, version)
}
