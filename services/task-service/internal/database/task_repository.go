package database

import (
	"github.com/williamschweitzer/task-management-app/services/task-service/internal/model"
)

func CreateTask(task model.Task) error {
	if err := DB.Create(&task).Error; err != nil {
		return err
	}

	return nil
}
