package services_test

import (
	"context"
	"fmt"
	"log"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
	"github.com/wanliqun/go-wallet-app/models"
	"github.com/wanliqun/go-wallet-app/services"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var (
	db            *gorm.DB
	userGenerator models.FakeUserGenerator
)

// Setup PostgreSQL container for tests
func TestMain(m *testing.M) {
	ctx := context.Background()

	// Request to start a PostgreSQL container
	req := testcontainers.ContainerRequest{
		Image:        "postgres:13", // Specify PostgreSQL version
		ExposedPorts: []string{"5432/tcp"},
		Env: map[string]string{
			"POSTGRES_USER":     "username",
			"POSTGRES_PASSWORD": "password",
			"POSTGRES_DB":       "testdb",
		},
		WaitingFor: wait.ForListeningPort("5432/tcp").WithStartupTimeout(10 * time.Second),
	}

	// Start the container
	postgresC, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	if err != nil {
		log.Fatalf("failed to start container: %v", err)
	}
	defer postgresC.Terminate(ctx) // Ensure the container is terminated after tests

	// Get container's host and port
	host, _ := postgresC.Host(ctx)
	port, _ := postgresC.MappedPort(ctx, "5432")

	// Connect to the PostgreSQL container
	dsn := fmt.Sprintf("host=%s port=%s user=username password=password dbname=testdb sslmode=disable", host, port.Port())
	db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}

	// Run auto-migrations
	db.AutoMigrate(&models.User{}, &models.Vault{}, &models.Transaction{})

	// Run the tests
	code := m.Run()

	// Exit with the test result code
	os.Exit(code)
}

func TestGetUserByName(t *testing.T) {
	// Begin a transaction for each test and roll it back afterward to ensure isolation
	tx := db.Begin()
	defer tx.Rollback()

	testuser := userGenerator.Generate()
	tx.Create(testuser)

	userService := services.NewUserService(tx)

	t.Run("should return user when user exists", func(t *testing.T) {
		user, found, err := userService.GetUserByName(testuser.Name)
		assert.NoError(t, err)
		assert.True(t, found)
		assert.Equal(t, testuser.Name, user.Name)
		assert.Equal(t, testuser.Email, user.Email)
	})

	t.Run("should return not found when user does not exist", func(t *testing.T) {
		user, found, err := userService.GetUserByName("nonexistentuser")
		assert.NoError(t, err)
		assert.False(t, found)
		assert.Nil(t, user)
	})
}
