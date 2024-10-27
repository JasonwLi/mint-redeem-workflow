package activities

import (
	"context"
	"fmt"
	"mint-redeem-workflow/db"
	"mint-redeem-workflow/models"

	"gorm.io/gorm"
)

func UpdateStatusActivity(ctx context.Context, requestID string, status string) error {
	var request models.Request
	if err := db.Db.First(&request, "id = ?", requestID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return fmt.Errorf("request with ID %s not found", requestID)
		}
		return err
	}

	request.Status = status
	if err := db.Db.Save(&request).Error; err != nil {
		return err
	}

	return nil
}
