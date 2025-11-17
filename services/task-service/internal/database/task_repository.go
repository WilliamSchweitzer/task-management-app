package database

import (
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/williamschweitzer/task-management-app/services/task-service/internal/model"
	"gorm.io/gorm"
)

func CreateTask(task model.Task) error {
	if err := DB.Create(&task).Error; err != nil {
		return err
	}

	return nil
}

func GetTask(taskID uuid.UUID) (*model.Task, error) {
	var task model.Task
	if err := DB.First(&task, taskID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("task not found with id: %s", taskID)
		}
		return nil, err
	}

	return &task, nil
}
