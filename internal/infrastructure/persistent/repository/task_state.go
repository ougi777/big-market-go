package repository

import (
	"context"
	"time"

	"bm-go/internal/infrastructure/persistent/po"

	"gorm.io/gorm"
)

func setTaskState(ctx context.Context, db *gorm.DB, userID string, messageID string, state string) error {
	return db.WithContext(ctx).
		Model(&po.Task{}).
		Where("user_id = ? and message_id = ?", userID, messageID).
		Updates(map[string]any{
			"state":       state,
			"update_time": time.Now(),
		}).
		Error
}
