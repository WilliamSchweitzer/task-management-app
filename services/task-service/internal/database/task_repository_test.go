package database

import (
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/williamschweitzer/task-management-app/services/task-service/internal/model"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func setupMockDB(t *testing.T) (*gorm.DB, sqlmock.Sqlmock) {
	sqlDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Failed to create mock: %v", err)
	}

	gormDB, err := gorm.Open(postgres.New(postgres.Config{
		Conn: sqlDB,
	}), &gorm.Config{})
	if err != nil {
		t.Fatalf("Failed to open gorm connection: %v", err)
	}

	return gormDB, mock
}

func TestCreateTask(t *testing.T) {
	gormDB, mock := setupMockDB(t)
	DB = gormDB

	userID := uuid.New()
	taskID := uuid.New()
	now := time.Now()

	task := model.Task{
		ID:        taskID,
		UserID:    userID,
		Title:     "Test Task",
		Status:    "todo",
		CreatedAt: now,
		UpdatedAt: now,
	}

	mock.ExpectBegin()
	mock.ExpectQuery(`INSERT INTO "tasks"."tasks"`).
		WithArgs(
			sqlmock.AnyArg(), // id
			sqlmock.AnyArg(), // user_id
			sqlmock.AnyArg(), // title
			sqlmock.AnyArg(), // description
			sqlmock.AnyArg(), // status
			sqlmock.AnyArg(), // priority
			sqlmock.AnyArg(), // due_date
			sqlmock.AnyArg(), // completed_at
			sqlmock.AnyArg(), // created_at
			sqlmock.AnyArg(), // updated_at
		).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(taskID))
	mock.ExpectCommit()

	err := CreateTask(&task)

	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGetTask(t *testing.T) {
	gormDB, mock := setupMockDB(t)
	DB = gormDB

	taskID := uuid.New()
	userID := uuid.New()
	now := time.Now()
	title := "Test Task"
	status := "todo"

	t.Run("task found", func(t *testing.T) {
		rows := sqlmock.NewRows([]string{
			"id", "user_id", "title", "description", "status", "priority",
			"due_date", "completed_at", "created_at", "updated_at",
		}).AddRow(
			taskID, userID, title, nil, status, nil,
			nil, nil, now, now,
		)

		mock.ExpectQuery(`SELECT \* FROM "tasks"."tasks"`).
			WithArgs(taskID, 1).
			WillReturnRows(rows)

		task, err := GetTask(taskID)

		assert.NoError(t, err)
		assert.NotNil(t, task)
		assert.Equal(t, taskID, task.ID)
		assert.Equal(t, userID, task.UserID)
		assert.Equal(t, title, task.Title)
		assert.Equal(t, status, task.Status)
	})

	t.Run("task not found", func(t *testing.T) {
		mock.ExpectQuery(`SELECT \* FROM "tasks"."tasks"`).
			WithArgs(taskID, 1).
			WillReturnError(gorm.ErrRecordNotFound)

		task, err := GetTask(taskID)

		assert.Error(t, err)
		assert.Nil(t, task)
		assert.Contains(t, err.Error(), "task not found")
	})
}

func TestUpdateTask(t *testing.T) {
	gormDB, mock := setupMockDB(t)
	DB = gormDB

	taskID := uuid.New()
	userID := uuid.New()
	now := time.Now()

	updates := map[string]interface{}{
		"title":      "Updated Title",
		"status":     "in-progress",
		"updated_at": now,
	}

	t.Run("successful update", func(t *testing.T) {
		mock.ExpectBegin()
		mock.ExpectExec(`UPDATE "tasks"."tasks"`).
			WillReturnResult(sqlmock.NewResult(0, 1))
		mock.ExpectCommit()

		rows := sqlmock.NewRows([]string{
			"id", "user_id", "title", "description", "status", "priority",
			"due_date", "completed_at", "created_at", "updated_at",
		}).AddRow(
			taskID, userID, "Updated Title", nil, "in-progress", nil,
			nil, nil, now, now,
		)

		mock.ExpectQuery(`SELECT \* FROM "tasks"."tasks"`).
			WithArgs(taskID, 1).
			WillReturnRows(rows)

		task, err := UpdateTask(taskID, updates)

		assert.NoError(t, err)
		assert.NotNil(t, task)
		assert.Equal(t, "Updated Title", task.Title)
		assert.Equal(t, "in-progress", task.Status)
	})

	t.Run("task not found", func(t *testing.T) {
		mock.ExpectBegin()
		mock.ExpectExec(`UPDATE "tasks"."tasks"`).
			WillReturnResult(sqlmock.NewResult(0, 0))
		mock.ExpectCommit()

		task, err := UpdateTask(taskID, updates)

		assert.Error(t, err)
		assert.Nil(t, task)
		assert.Contains(t, err.Error(), "task not found")
	})
}

func TestDeleteTask(t *testing.T) {
	gormDB, mock := setupMockDB(t)
	DB = gormDB

	taskID := uuid.New()

	t.Run("successful delete", func(t *testing.T) {
		mock.ExpectBegin()
		mock.ExpectExec(`DELETE FROM "tasks"."tasks"`).
			WithArgs(taskID).
			WillReturnResult(sqlmock.NewResult(0, 1))
		mock.ExpectCommit()

		err := DeleteTask(taskID)

		assert.NoError(t, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("task not found", func(t *testing.T) {
		mock.ExpectBegin()
		mock.ExpectExec(`DELETE FROM "tasks"."tasks"`).
			WithArgs(taskID).
			WillReturnResult(sqlmock.NewResult(0, 0))
		mock.ExpectCommit()

		err := DeleteTask(taskID)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "task not found")
	})
}

func TestCompleteTask(t *testing.T) {
	gormDB, mock := setupMockDB(t)
	DB = gormDB

	taskID := uuid.New()
	userID := uuid.New()
	now := time.Now()

	t.Run("successful completion", func(t *testing.T) {
		mock.ExpectBegin()
		mock.ExpectExec(`UPDATE "tasks"."tasks"`).
			WillReturnResult(sqlmock.NewResult(0, 1))
		mock.ExpectCommit()

		rows := sqlmock.NewRows([]string{
			"id", "user_id", "title", "description", "status", "priority",
			"due_date", "completed_at", "created_at", "updated_at",
		}).AddRow(
			taskID, userID, "Test Task", nil, "done", nil,
			nil, now, now, now,
		)

		mock.ExpectQuery(`SELECT \* FROM "tasks"."tasks"`).
			WithArgs(taskID, 1).
			WillReturnRows(rows)

		task, err := CompleteTask(taskID)

		assert.NoError(t, err)
		assert.NotNil(t, task)
		assert.Equal(t, "done", task.Status)
		assert.NotNil(t, task.CompletedAt)
	})

	t.Run("task not found", func(t *testing.T) {
		mock.ExpectBegin()
		mock.ExpectExec(`UPDATE "tasks"."tasks"`).
			WillReturnResult(sqlmock.NewResult(0, 0))
		mock.ExpectCommit()

		task, err := CompleteTask(taskID)

		assert.Error(t, err)
		assert.Nil(t, task)
		assert.Contains(t, err.Error(), "task not found")
	})
}
