package api

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"TODO/todo"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockTodoService is a mock implementation of the TodoService interface
type MockTodoService struct {
	mock.Mock
}

func (m *MockTodoService) CreateTodoService(body io.Reader) (*todo.Todo, error) {
	args := m.Called(body)
	return args.Get(0).(*todo.Todo), args.Error(1)
}

func (m *MockTodoService) GetTodosService() ([]todo.Todo, error) {
	args := m.Called()
	return args.Get(0).([]todo.Todo), args.Error(1)
}

func (m *MockTodoService) UpdateTodoService(id int, body io.Reader) (*todo.Todo, error) {
	args := m.Called(id, body)
	return args.Get(0).(*todo.Todo), args.Error(1)
}

func (m *MockTodoService) DeleteTodoService(id int, body io.Reader) error {
	args := m.Called(id, body)
	return args.Error(0)
}

func TestCreate(t *testing.T) {
	// Arrange
	mockService := new(MockTodoService)
	handler := TodoHandler{Service: mockService}

	newTodo := &todo.Todo{ID: 1, Title: "New Todo", Status: "Pending"}
	mockService.On("CreateTodoService", mock.Anything).Return(newTodo, nil)

	reqBody := `{"title":"New Todo","status":"Pending"}`
	req, err := http.NewRequest("POST", "/todo", bytes.NewBuffer([]byte(reqBody)))
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()

	// Act
	handler.Create(rr, req)

	// Assert
	assert.Equal(t, http.StatusCreated, rr.Code)
	var responseTodo todo.Todo
	err = json.NewDecoder(rr.Body).Decode(&responseTodo)
	assert.NoError(t, err)
	assert.Equal(t, *newTodo, responseTodo)

	// Verify the expectations
	mockService.AssertExpectations(t)
}
func TestCreateFailure(t *testing.T) {
	// Arrange
	mockService := new(MockTodoService)
	handler := TodoHandler{Service: mockService}

	// Simulate an error from the service
	newTodo := &todo.Todo{ID: 1, Title: "New Todo", Status: "Pending"}
	mockService.On("CreateTodoService", mock.Anything).Return(newTodo, errors.New("service error"))

	// Provide invalid input that should trigger a bad request (e.g., incorrect data type)
	reqBody := `{"title"::New Todo,"status":"Pending"}`
	req, err := http.NewRequest("POST", "/todo", bytes.NewBuffer([]byte(reqBody)))
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()

	// Act
	handler.Create(rr, req)

	// Assert
	assert.Equal(t, http.StatusBadRequest, rr.Code)
	mockService.AssertExpectations(t)
}

func TestGet(t *testing.T) {
	// Arrange
	mockService := new(MockTodoService)
	handler := TodoHandler{Service: mockService}

	todos := []todo.Todo{
		{ID: 1, Title: "Todo 1", Status: "Pending"},
		{ID: 2, Title: "Todo 2", Status: "Completed"},
	}
	mockService.On("GetTodosService").Return(todos, nil)

	req, err := http.NewRequest("GET", "/todo", nil)
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()

	// Act
	handler.Get(rr, req)

	// Assert
	assert.Equal(t, http.StatusOK, rr.Code)
	var responseTodos []todo.Todo
	err = json.NewDecoder(rr.Body).Decode(&responseTodos)
	assert.NoError(t, err)
	assert.Equal(t, todos, responseTodos)

	// Verify the expectations
	mockService.AssertExpectations(t)
}
func TestGetFailure(t *testing.T) {
	// Arrange
	mockService := new(MockTodoService)
	handler := TodoHandler{Service: mockService}

	todos := []todo.Todo{}
	mockService.On("GetTodosService").Return(todos, errors.New("failed to get"))

	req, err := http.NewRequest("GET", "/todo", nil)
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()

	// Act
	handler.Get(rr, req)

	// Assert
	assert.Equal(t, http.StatusInternalServerError, rr.Code)

	// Verify the expectations
	mockService.AssertExpectations(t)
}

func TestUpdate(t *testing.T) {
	// Arrange
	mockService := new(MockTodoService)
	handler := TodoHandler{Service: mockService}

	// Define the expected result from the mock
	updatedTodo := todo.Todo{
		ID:     2,
		Title:  "Updated Todo",
		Status: "Completed",
	}
	mockService.On("UpdateTodoService", 2, mock.Anything).Return(&updatedTodo, nil)

	// Prepare the request with valid JSON data
	reqBody := `{"title":"Updated Todo","status":"Completed"}`
	req, err := http.NewRequest("PATCH", "/todo/2", bytes.NewBuffer([]byte(reqBody)))
	if err != nil {
		t.Fatal(err)
	}

	// Use the router to handle the request
	router := SetupRoutes(&handler)
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	// Debugging: Print the response body for inspection
	// t.Log("Response Body:", rr.Body.String())

	// Assert
	assert.Equal(t, http.StatusOK, rr.Code) // Expecting 200 OK

	var responseTodo todo.Todo
	err = json.NewDecoder(rr.Body).Decode(&responseTodo)
	if err != nil {
		t.Fatal("Failed to decode response body:", err)
	}
	assert.Equal(t, updatedTodo, responseTodo)

	// Verify the expectations
	mockService.AssertExpectations(t)
}
func TestUpdateFailure(t *testing.T) {
	// Arrange
	mockService := new(MockTodoService)
	handler := TodoHandler{Service: mockService}

	// Simulate the service returning an error, causing the handler to return a BadRequest
	mockService.On("UpdateTodoService", 1, mock.Anything).Return(nil, errors.New("update failed"))

	// Prepare a request with valid JSON data
	reqBody := `{"title":"Updated Todo","status":"Completed"}`
	req, err := http.NewRequest("PATCH", "/todo/1", bytes.NewBuffer([]byte(reqBody)))
	if err != nil {
		t.Fatal(err)
		t.Log("error creating HTTP request")
	}
	rr := httptest.NewRecorder()

	// Act
	handler.Update(rr, req)

	// Assert
	assert.Equal(t, http.StatusBadRequest, rr.Code) // Expecting BadRequest due to the simulated service error

	// mockService.AssertExpectations(t)
}

func TestDelete(t *testing.T) {
	// Arrange
	mockService := new(MockTodoService)
	handler := TodoHandler{Service: mockService}

	// Define the expected result from the mock
	mockService.On("DeleteTodoService", 1, mock.Anything).Return(nil)

	// Prepare the request
	req, err := http.NewRequest("DELETE", "/todo/1", nil)
	if err != nil {
		t.Fatal(err)
	}

	// Use the router to handle the request
	router := SetupRoutes(&handler)
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	// Assert
	assert.Equal(t, http.StatusNoContent, rr.Code) // Expecting 204 No Content

	// Verify the expectations
	mockService.AssertExpectations(t)
}
func TestDeleteFailure(t *testing.T) {
	// Arrange
	mockService := new(MockTodoService)
	handler := TodoHandler{Service: mockService}

	// Define the expected result from the mock
	mockService.On("DeleteTodoService", 1, mock.Anything).Return(errors.New("couldn't delete"))

	// Prepare a request with valid JSON data
	reqBody := `{"title":"Updated Todo","status":"Completed"}`
	req, err := http.NewRequest("DELETE", "/todo/1", bytes.NewBuffer([]byte(reqBody)))
	if err != nil {
		t.Fatal(err)
		t.Log("error creating HTTP request")
	}
	rr := httptest.NewRecorder()

	// Act
	handler.Delete(rr, req)

	// Assert
	assert.Equal(t, http.StatusBadRequest, rr.Code) // Expecting BadRequest due to the simulated service error

	// mockService.AssertExpectations(t)
}
