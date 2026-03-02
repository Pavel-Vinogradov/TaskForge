//go:build integration

package usecase

import (
	"TaskForge/internal/contextkeys"
	"TaskForge/internal/domain/repos"
	"TaskForge/internal/infrastructure/repository"
	"TaskForge/internal/interfaces/task"
	"context"
	"database/sql"
	"fmt"
	"testing"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/mysql"
	"github.com/testcontainers/testcontainers-go/wait"
)

type TaskIntegrationTestSuite struct {
	suite.Suite
	db        *sql.DB
	repo      repos.TaskRepository
	uc        *TaskUseCase
	container testcontainers.Container
}

func (suite *TaskIntegrationTestSuite) SetupSuite() {
	ctx := context.Background()

	mysqlContainer, err := mysql.RunContainer(ctx,
		testcontainers.WithImage("mysql:8.0"),
		mysql.WithDatabase("taskforge_test"),
		mysql.WithUsername("test"),
		mysql.WithPassword("test"),
		testcontainers.WithWaitStrategy(
			wait.ForLog("port: 3306  MySQL Community Server - GPL"),
		),
	)
	require.NoError(suite.T(), err)
	suite.container = mysqlContainer

	connectionString, err := mysqlContainer.ConnectionString(ctx, "parseTime=true")
	require.NoError(suite.T(), err)

	db, err := sql.Open("mysql", connectionString)
	require.NoError(suite.T(), err)

	require.Eventually(suite.T(), func() bool {
		err := db.Ping()
		return err == nil
	}, 30*time.Second, time.Second)

	suite.db = db

	suite.runMigrations()

	suite.repo = repository.NewTaskRepository(db)
	suite.uc = NewTaskUseCase(suite.repo, nil, nil)
}

func (suite *TaskIntegrationTestSuite) TearDownSuite() {
	if suite.db != nil {
		suite.db.Close()
	}
	if suite.container != nil {
		ctx := context.Background()
		_ = suite.container.Terminate(ctx)
	}
}

func (suite *TaskIntegrationTestSuite) SetupTest() {
	_, err := suite.db.Exec("DELETE FROM tasks")
	require.NoError(suite.T(), err)
}

func (suite *TaskIntegrationTestSuite) runMigrations() {
	// Create tasks table
	createTableSQL := `
	CREATE TABLE IF NOT EXISTS tasks (
		id INT AUTO_INCREMENT PRIMARY KEY,
		title VARCHAR(255) NOT NULL,
		description TEXT,
		status ENUM('todo', 'in_progress', 'done') NOT NULL DEFAULT 'todo',
		team_id INT NOT NULL,
		created_by INT NOT NULL,
		assignee_id INT NOT NULL,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
	);`

	_, err := suite.db.Exec(createTableSQL)
	require.NoError(suite.T(), err)
}

func (suite *TaskIntegrationTestSuite) TestCreateTask_Success() {
	ctx := context.WithValue(context.Background(), contextkeys.UserIDKey, 42)

	req := task.CreateTaskRequest{
		Title:       "Integration Test Task",
		Description: "This is an integration test task",
		TeamID:      10,
		AssigneeID:  5,
	}

	res, err := suite.uc.CreateTask(ctx, req)

	assert.NoError(suite.T(), err)
	assert.NotZero(suite.T(), res.ID)
	assert.Equal(suite.T(), "Integration Test Task", res.Title)
	assert.Equal(suite.T(), "This is an integration test task", res.Description)
	assert.Equal(suite.T(), "todo", res.Status)
	assert.Equal(suite.T(), 10, res.TeamID)
	assert.Equal(suite.T(), 42, res.CreatedBy)
	assert.Equal(suite.T(), 5, res.AssigneeID)
	assert.NotZero(suite.T(), res.CreatedAt)

	// Verify task exists in database
	dbTask, err := suite.repo.GetTaskByID(ctx, res.ID)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), res.ID, dbTask.Id)
	assert.Equal(suite.T(), "Integration Test Task", dbTask.Title)
}

func (suite *TaskIntegrationTestSuite) TestCreateTask_NoUser() {
	ctx := context.Background()

	req := task.CreateTaskRequest{
		Title:       "Test Task",
		Description: "Description",
		TeamID:      10,
		AssigneeID:  5,
	}

	_, err := suite.uc.CreateTask(ctx, req)

	assert.Error(suite.T(), err)
	assert.Contains(suite.T(), err.Error(), "user not authenticated")
}

func (suite *TaskIntegrationTestSuite) TestUpdateTask_Success() {
	ctx := context.WithValue(context.Background(), contextkeys.UserIDKey, 42)

	createReq := task.CreateTaskRequest{
		Title:       "Original Title",
		Description: "Original Description",
		TeamID:      10,
		AssigneeID:  5,
	}
	createdTask, err := suite.uc.CreateTask(ctx, createReq)
	require.NoError(suite.T(), err)

	newTitle := "Updated Title"
	newDesc := "Updated Description"
	updateReq := task.UpdateTaskRequest{
		Title:       &newTitle,
		Description: &newDesc,
	}

	updatedTask, err := suite.uc.UpdateTask(ctx, createdTask.ID, updateReq)

	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), createdTask.ID, updatedTask.ID)
	assert.Equal(suite.T(), "Updated Title", updatedTask.Title)
	assert.Equal(suite.T(), "Updated Description", updatedTask.Description)
	assert.Equal(suite.T(), "todo", updatedTask.Status) // Status unchanged
}

func (suite *TaskIntegrationTestSuite) TestUpdateTask_NotFound() {
	ctx := context.WithValue(context.Background(), contextkeys.UserIDKey, 42)

	req := task.UpdateTaskRequest{
		Title: stringPtr("Updated Title"),
	}

	_, err := suite.uc.UpdateTask(ctx, 999, req)

	assert.Error(suite.T(), err)
	assert.Contains(suite.T(), err.Error(), "task not found")
}

func (suite *TaskIntegrationTestSuite) TestUpdateTask_NoPermission() {
	ctx := context.WithValue(context.Background(), contextkeys.UserIDKey, 42)

	createReq := task.CreateTaskRequest{
		Title:       "Test Task",
		Description: "Description",
		TeamID:      10,
		AssigneeID:  5,
	}
	createdTask, err := suite.uc.CreateTask(ctx, createReq)
	require.NoError(suite.T(), err)

	ctx = context.WithValue(context.Background(), contextkeys.UserIDKey, 99)
	req := task.UpdateTaskRequest{
		Title: stringPtr("Updated Title"),
	}

	_, err = suite.uc.UpdateTask(ctx, createdTask.ID, req)

	assert.Error(suite.T(), err)
	assert.Contains(suite.T(), err.Error(), "insufficient permissions")
}

func (suite *TaskIntegrationTestSuite) TestListTasks_Success() {
	ctx := context.WithValue(context.Background(), contextkeys.UserIDKey, 42)

	for i := 0; i < 5; i++ {
		req := task.CreateTaskRequest{
			Title:       fmt.Sprintf("Task %d", i),
			Description: fmt.Sprintf("Description %d", i),
			TeamID:      10,
			AssigneeID:  5,
		}
		_, err := suite.uc.CreateTask(ctx, req)
		require.NoError(suite.T(), err)
	}

	listReq := task.TaskListRequest{
		TeamID: intPtr(10),
		Page:   1,
		Limit:  10,
	}

	result, err := suite.uc.ListTask(ctx, listReq)

	assert.NoError(suite.T(), err)
	assert.Len(suite.T(), result.Tasks, 5)
	assert.Equal(suite.T(), int64(5), result.Total)
}

func (suite *TaskIntegrationTestSuite) TestGetTaskByID_Success() {
	ctx := context.WithValue(context.Background(), contextkeys.UserIDKey, 42)

	createReq := task.CreateTaskRequest{
		Title:       "Test Task",
		Description: "Description",
		TeamID:      10,
		AssigneeID:  5,
	}
	createdTask, err := suite.uc.CreateTask(ctx, createReq)
	require.NoError(suite.T(), err)

	task, err := suite.repo.GetTaskByID(ctx, createdTask.ID)

	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), createdTask.ID, task.Id)
	assert.Equal(suite.T(), "Test Task", task.Title)
}

func TestTaskIntegrationSuite(t *testing.T) {
	suite.Run(t, new(TaskIntegrationTestSuite))
}

func stringPtr(s string) *string {
	return &s
}

func intPtr(i int) *int {
	return &i
}
