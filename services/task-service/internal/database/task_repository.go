package database

import (
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/williamschweitzer/task-management-app/services/task-service/internal/model"
	"gorm.io/gorm"
)

func CreateTask(task *model.Task) error {
	if task == nil {
		return fmt.Errorf("task cannot be nil")
	}

	if err := DB.Create(task).Error; err != nil {
		return err
	}

	return nil
}

func GetTasksByUserID(userID uuid.UUID) ([]model.Task, error) {
	var tasks []model.Task

	result := DB.Where("user_id = ?", userID).Find(&tasks)
	if result.Error != nil {
		return nil, result.Error
	}

	if len(tasks) == 0 {
		return []model.Task{}, nil // Return empty slice, not error
	}

	return tasks, nil
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

func UpdateTask(taskID uuid.UUID, updates map[string]interface{}) (*model.Task, error) {
	result := DB.Model(&model.Task{}).Where("id = ?", taskID).Updates(updates)

	if result.Error != nil {
		return nil, fmt.Errorf("failed to update task: %w", result.Error)
	}

	if result.RowsAffected == 0 {
		return nil, fmt.Errorf("task not found")
	}

	var task model.Task
	if err := DB.First(&task, taskID).Error; err != nil {
		return nil, fmt.Errorf("failed to fetch updated task: %w", err)
	}

	return &task, nil
}

func DeleteTask(taskID uuid.UUID) error {
	result := DB.Delete(&model.Task{}, taskID)

	if result.Error != nil {
		return fmt.Errorf("failed to delete task: %w", result.Error)
	}

	if result.RowsAffected == 0 {
		return fmt.Errorf("task not found")
	}

	return nil
}

func CompleteTask(taskID uuid.UUID) (*model.Task, error) {
	updates := map[string]interface{}{
		"status":       "done",
		"completed_at": time.Now(),
		"updated_at":   time.Now(),
	}

	result := DB.Model(&model.Task{}).Where("id = ?", taskID).Updates(updates)

	if result.Error != nil {
		return nil, fmt.Errorf("failed to complete task: %w", result.Error)
	}

	if result.RowsAffected == 0 {
		return nil, fmt.Errorf("task not found")
	}

	var task model.Task
	if err := DB.First(&task, taskID).Error; err != nil {
		return nil, fmt.Errorf("failed to fetch completed task: %w", err)
	}

	return &task, nil
}
