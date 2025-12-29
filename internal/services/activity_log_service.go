package services

import (
	"encoding/json"

	"project-bulky-be/internal/models"
	"project-bulky-be/internal/repositories"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type ActivityLogService interface {
	// Logging
	Log(ctx *gin.Context, action models.ActivityAction, modul, deskripsi string, options ...LogOption)
	LogCreate(ctx *gin.Context, modul, entityType string, entityID uuid.UUID, deskripsi string, newData interface{})
	LogUpdate(ctx *gin.Context, modul, entityType string, entityID uuid.UUID, deskripsi string, oldData, newData interface{})
	LogDelete(ctx *gin.Context, modul, entityType string, entityID uuid.UUID, deskripsi string, oldData interface{})

	// Query
	GetLogs(filter repositories.ActivityLogFilter) ([]models.ActivityLog, int64, error)
	GetLogByID(id uuid.UUID) (*models.ActivityLog, error)
	GetLogsByEntity(entityType string, entityID uuid.UUID) ([]models.ActivityLog, error)
}

type LogOption func(*models.ActivityLog)

type activityLogService struct {
	repo repositories.ActivityLogRepository
}

func NewActivityLogService(repo repositories.ActivityLogRepository) ActivityLogService {
	return &activityLogService{repo: repo}
}

// Log creates a new activity log entry
func (s *activityLogService) Log(ctx *gin.Context, action models.ActivityAction, modul, deskripsi string, options ...LogOption) {
	log := &models.ActivityLog{
		Action:    action,
		Modul:     modul,
		Deskripsi: deskripsi,
	}

	// Extract user info from context
	if userID, exists := ctx.Get("user_id"); exists {
		if uid, err := uuid.Parse(userID.(string)); err == nil {
			log.UserID = &uid
		}
	}

	if userType, exists := ctx.Get("user_type"); exists {
		log.UserType = userType.(string)
	} else {
		log.UserType = "SYSTEM"
	}

	// Extract IP and User-Agent
	ipAddress := ctx.ClientIP()
	log.IPAddress = &ipAddress

	userAgent := ctx.GetHeader("User-Agent")
	log.UserAgent = &userAgent

	// Apply options
	for _, opt := range options {
		opt(log)
	}

	// Save asynchronously
	go s.repo.Create(log)
}

// LogCreate logs a CREATE action
func (s *activityLogService) LogCreate(ctx *gin.Context, modul, entityType string, entityID uuid.UUID, deskripsi string, newData interface{}) {
	newDataJSON, _ := json.Marshal(newData)

	s.Log(ctx, models.ActionCreate, modul, deskripsi,
		WithEntity(entityType, entityID),
		WithNewData(newDataJSON),
	)
}

// LogUpdate logs an UPDATE action
func (s *activityLogService) LogUpdate(ctx *gin.Context, modul, entityType string, entityID uuid.UUID, deskripsi string, oldData, newData interface{}) {
	oldDataJSON, _ := json.Marshal(oldData)
	newDataJSON, _ := json.Marshal(newData)

	s.Log(ctx, models.ActionUpdate, modul, deskripsi,
		WithEntity(entityType, entityID),
		WithOldData(oldDataJSON),
		WithNewData(newDataJSON),
	)
}

// LogDelete logs a DELETE action
func (s *activityLogService) LogDelete(ctx *gin.Context, modul, entityType string, entityID uuid.UUID, deskripsi string, oldData interface{}) {
	oldDataJSON, _ := json.Marshal(oldData)

	s.Log(ctx, models.ActionDelete, modul, deskripsi,
		WithEntity(entityType, entityID),
		WithOldData(oldDataJSON),
	)
}

// GetLogs retrieves activity logs with filters
func (s *activityLogService) GetLogs(filter repositories.ActivityLogFilter) ([]models.ActivityLog, int64, error) {
	return s.repo.FindAll(filter)
}

// GetLogByID retrieves a specific log by ID
func (s *activityLogService) GetLogByID(id uuid.UUID) (*models.ActivityLog, error) {
	return s.repo.FindByID(id)
}

// GetLogsByEntity retrieves all logs for a specific entity
func (s *activityLogService) GetLogsByEntity(entityType string, entityID uuid.UUID) ([]models.ActivityLog, error) {
	return s.repo.FindByEntity(entityType, entityID)
}

// Log options
func WithEntity(entityType string, entityID uuid.UUID) LogOption {
	return func(log *models.ActivityLog) {
		log.EntityType = &entityType
		log.EntityID = &entityID
	}
}

func WithOldData(data json.RawMessage) LogOption {
	return func(log *models.ActivityLog) {
		log.OldData = data
	}
}

func WithNewData(data json.RawMessage) LogOption {
	return func(log *models.ActivityLog) {
		log.NewData = data
	}
}

func WithUserInfo(userID uuid.UUID, userType string) LogOption {
	return func(log *models.ActivityLog) {
		log.UserID = &userID
		log.UserType = userType
	}
}
