package todo

import (
	"TODO/todo/config"

	"bytes"
	"log"
	"os"
	"testing"

	"github.com/golang/mock/gomock"
	_ "github.com/lib/pq"
	"github.com/stretchr/testify/assert"
)

// MockDB is a mock implementation of sql.DB
type MockDB struct {
	mockCtrl *gomock.Controller
}

func NewMockDB(ctrl *gomock.Controller) *MockDB {
	return &MockDB{mockCtrl: ctrl}
}

func (db *MockDB) Ping() error {
	// Mocked Ping behavior
	return nil
}
func TestInitDB(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// // Set environment variables to control the behavior during the test
	// os.Setenv("DB_USER", "testuser")
	// os.Setenv("DB_PASSWORD", "testpass")
	// os.Setenv("DB_NAME", "testdb")
	// os.Setenv("DB_HOST", "localhost")
	// os.Setenv("DB_SSLMODE", "disable")

	// defer os.Clearenv()

	// Load configuration
	configg := config.InitConfig()

	// Call the InitDB function
	InitDB(configg)

	// Check that the DB variable has been set
	assert.NotNil(t, DB)
}

func TestInitDBFailure(t *testing.T) {
	// Create a custom invalid configuration
	configg := config.Config{
		Database: config.DatabaseConfig{
			User:     "invaliduser",
			Password: "invalidpass",
			Dbname:   "invaliddb",
			Host:     "localhost",
			Sslmode:  "disable",
		},
		// Other fields can be filled as needed but are not critical for this test
	}

	// Capture logs to check for error messages
	var logBuffer bytes.Buffer
	log.SetOutput(&logBuffer)
	defer log.SetOutput(os.Stderr) // Restore default log output

	// Initialize the database with invalid configuration
	InitDB(configg)

	// Check that the DB variable is nil (indicating that the connection failed)
	if DB != nil {
		t.Errorf("DB should not be initialized with invalid configuration")
	}

	// Check for specific error messages in the logs
	logOutput := logBuffer.String()
	if assert.Contains(t, logOutput, "Failed to ping database:", "Expected error message not found in logs") {
		t.Logf("Log output: %s", logOutput) // Log the actual output for debugging
	} else {
		t.Errorf("Expected error message not found in logs: %s", logOutput)
	}
}

func TestRunMigrations(t *testing.T) {
	// Set environment variables to simulate a test environment
	// os.Setenv("DB_USER", "testuser")
	// os.Setenv("DB_PASSWORD", "testpass")
	// os.Setenv("DB_NAME", "testdb")
	// os.Setenv("DB_HOST", "localhost")
	// os.Setenv("DB_SSLMODE", "disable")
	// os.Setenv("MIGRATIONS_PATH", "../db/migrations")

	// defer os.Clearenv()

	// Load configuration
	configg := config.InitConfig()

	// Simulate running migrations
	defer func() {
		if r := recover(); r != nil {
			assert.Contains(t, r, "no change")
		}
	}()

	RunMigrations(configg)

	// No specific assertions needed since migration output is logged
}

func TestNewRedisClient(t *testing.T) {
	configg := config.Config{
		Redis: config.RedisConfig{
			Address:  "localhost:6379",
			Password: "",
			DB:       0,
		},
	}

	redisClient := NewRedisClient(configg)
	assert.NotNil(t, redisClient)

	pong, err := redisClient.Ping(redisClient.Context()).Result()
	assert.NoError(t, err)
	assert.Equal(t, "PONG", pong)
}

func TestInitFaktory(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockFaktoryClient := NewMockFaktoryClient(ctrl)

	configg := config.Config{
		Faktory: config.FaktoryConfig{
			URL: "tcp://localhost:7419",
		},
	}

	InitFaktory(configg)

	assert.NotNil(t, mockFaktoryClient.Client)
}

func TestInitConfig(t *testing.T) {
	config := config.InitConfig()
	assert.NotNil(t, config)
	assert.NotEmpty(t, config.Database.User)
	assert.NotEmpty(t, config.Redis.Address)
	assert.NotEmpty(t, config.Faktory.URL)
}
