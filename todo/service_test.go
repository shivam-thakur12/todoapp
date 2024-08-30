package todo

import (
	"TODO/todo/config"
	"bytes"
	"io"
	"testing"
	"time"

	"github.com/contribsys/faktory/client"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestCreateTodoService(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := NewMockTodoRepository(ctrl)
	mockCache := NewMockRedisCache(ctrl)
	mockFaktoryClient := NewMockFaktoryClient(ctrl)

	todoService := NewTodoService(mockRepo, mockCache, &mockFaktoryClient.Client)

	todoJSON := `{"title": "Test Todo", "status": "pending"}`
	body := io.NopCloser(bytes.NewBufferString(todoJSON))

	todo, err := todoService.CreateTodoService(body)

	assert.NoError(t, err)
	assert.NotNil(t, todo)
	assert.Equal(t, "Test Todo", todo.Title)
	assert.Equal(t, "pending", todo.Status)
}

func TestGetTodosService(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := NewMockTodoRepository(ctrl)
	mockCache := NewMockRedisCache(ctrl)
	mockFaktoryClient := NewMockFaktoryClient(ctrl)

	todoService := NewTodoService(mockRepo, mockCache, &mockFaktoryClient.Client)

	todos, err := todoService.GetTodosService()

	assert.NoError(t, err)
	assert.NotNil(t, todos)
	assert.Len(t, todos, 1)
	assert.Equal(t, "Test Todo", todos[0].Title)
}

func TestUpdateTodoService(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := NewMockTodoRepository(ctrl)
	mockCache := NewMockRedisCache(ctrl)
	mockFaktoryClient := NewMockFaktoryClient(ctrl)

	todoService := NewTodoService(mockRepo, mockCache, &mockFaktoryClient.Client)

	todoJSON := `{"title": "Updated Todo", "status": "completed"}`
	body := io.NopCloser(bytes.NewBufferString(todoJSON))

	updatedTodo, err := todoService.UpdateTodoService(1, body)

	assert.NoError(t, err)
	assert.NotNil(t, updatedTodo)
	assert.Equal(t, "Updated Todo", updatedTodo.Title)
	assert.Equal(t, "completed", updatedTodo.Status)
}
func TestUpdateTodoServicFailure(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := NewMockTodoRepository(ctrl)
	mockCache := NewMockRedisCache(ctrl)
	mockFaktoryClient := NewMockFaktoryClient(ctrl)

	todoService := NewTodoService(mockRepo, mockCache, &mockFaktoryClient.Client)

	todoJSON := `{"title": "Updated Todo", "status": "completed"}`
	body := io.NopCloser(bytes.NewBufferString(todoJSON))

	updatedTodo, err := todoService.UpdateTodoService(2, body)

	assert.NoError(t, err)
	assert.NotNil(t, updatedTodo)
	assert.Equal(t, "Updated Todo", updatedTodo.Title)
	assert.Equal(t, "completed", updatedTodo.Status)
}
func TestDeleteTodoService(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	mockRepo := NewMockTodoRepository(mockCtrl)
	mockCache := NewMockRedisCache(mockCtrl)

	configg := config.InitConfig()
	faktoryClient := InitFaktory(configg)

	defer faktoryClient.Close()

	todoService := NewTodoService(mockRepo, mockCache, faktoryClient)

	err := todoService.DeleteTodoService(1, nil)

	assert.NoError(t, err)
}

type MockTodoRepository struct {
	mock *gomock.Controller
}

func NewMockTodoRepository(ctrl *gomock.Controller) *MockTodoRepository {
	return &MockTodoRepository{
		mock: ctrl,
	}
}

func (m *MockTodoRepository) CreateTodoRepo(todo *Todo) error {
	return nil
}

func (m *MockTodoRepository) GetTodosRepo() ([]Todo, error) {
	return []Todo{
		{ID: 1, Title: "Test Todo", Status: "pending"},
	}, nil
}

func (m *MockTodoRepository) UpdateTodoRepo(todo *Todo) (int64, error) {
	return 1, nil
}

func (m *MockTodoRepository) DeleteTodoRepo(id int) (int64, error) {
	return 1, nil
}

type MockRedisCache struct {
	mock *gomock.Controller
}

func NewMockRedisCache(ctrl *gomock.Controller) *MockRedisCache {
	return &MockRedisCache{
		mock: ctrl,
	}
}

func (m *MockRedisCache) Set(key string, value interface{}, expiration time.Duration) error {
	// Your mock implementation here
	return nil
}

func (m *MockRedisCache) Get(key string) (string, error) {
	return "", nil
}

func (m *MockRedisCache) Del(key string) error {
	return nil
}

type FaktoryClient interface {
	Push(job *client.Job) error
}

type MockFaktoryClient struct {
	client.Client
	mock *gomock.Controller
}

func NewMockFaktoryClient(ctrl *gomock.Controller) *MockFaktoryClient {
	return &MockFaktoryClient{
		Client: client.Client{}, // Initialize with a real or mocked client
		mock:   ctrl,
	}
}

func (m *MockFaktoryClient) Push(job *client.Job) error {
	return nil
}
