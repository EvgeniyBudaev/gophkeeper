package app

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/EvgeniyBudaev/gophkeeper/internal/adapters/store"
	"github.com/EvgeniyBudaev/gophkeeper/internal/config"
	"github.com/EvgeniyBudaev/gophkeeper/internal/logger"
	"github.com/EvgeniyBudaev/gophkeeper/internal/models"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestRegister(t *testing.T) {
	l, _ := logger.NewLogger()
	c, err := config.Load(l)
	if err != nil {
		fmt.Errorf("failed config load: %w", err)
		return
	}
	conn, err := store.NewPostgresConnection(c)
	if err != nil {
		fmt.Errorf("failed new postgres connection: %w", err)
		return
	}
	app := NewApp(&config.ServerConfig{}, store.NewStore(conn), l)
	testCases := []struct {
		name           string
		requestBody    models.User
		expectedStatus int
	}{
		{
			name: "Successful user registration",
			requestBody: models.User{
				Login:    "testuser",
				Password: "testpassword",
			},
			expectedStatus: http.StatusCreated,
		},
		{
			name: "Duplicate login error",
			requestBody: models.User{
				Login:    "duplicateuser",
				Password: "testpassword",
			},
			expectedStatus: http.StatusConflict,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			reqBody, _ := json.Marshal(tc.requestBody)
			req, _ := http.NewRequest("POST", "/api/user/register", bytes.NewBuffer(reqBody))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)
			c.Request = req
			app.Register(c)
			assert.Equal(t, tc.expectedStatus, w.Code)
		})
	}
}

func TestLogin(t *testing.T) {
	l, _ := logger.NewLogger()
	c, err := config.Load(l)
	if err != nil {
		fmt.Errorf("failed config load: %w", err)
		return
	}
	conn, err := store.NewPostgresConnection(c)
	if err != nil {
		fmt.Errorf("failed new postgres connection: %w", err)
		return
	}
	app := NewApp(&config.ServerConfig{}, store.NewStore(conn), l)
	testCases := []struct {
		name           string
		requestBody    models.User
		expectedStatus int
	}{
		{
			name: "Successful login",
			requestBody: models.User{
				Login:    "testuser",
				Password: "testpassword",
			},
			expectedStatus: http.StatusOK,
		},
		{
			name: "Incorrect credentials",
			requestBody: models.User{
				Login:    "testuser",
				Password: "wrongpassword",
			},
			expectedStatus: http.StatusUnauthorized,
		},
		{
			name: "Login not found",
			requestBody: models.User{
				Login:    "nonexistentuser",
				Password: "testpassword",
			},
			expectedStatus: http.StatusUnauthorized,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			reqBody, _ := json.Marshal(tc.requestBody)
			req, _ := http.NewRequest("POST", "/api/user/login", bytes.NewBuffer(reqBody))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)
			c.Request = req
			app.Login(c)
			assert.Equal(t, tc.expectedStatus, w.Code)
		})
	}
}

func TestPutDataRecord(t *testing.T) {
	l, _ := logger.NewLogger()
	c, err := config.Load(l)
	if err != nil {
		fmt.Errorf("failed config load: %w", err)
		return
	}
	conn, err := store.NewPostgresConnection(c)
	if err != nil {
		fmt.Errorf("failed new postgres connection: %w", err)
		return
	}
	app := NewApp(&config.ServerConfig{}, store.NewStore(conn), l)
	testCases := []struct {
		name           string
		requestBody    models.DataRecordRequest
		expectedStatus int
	}{
		{
			name: "Successful data record creation",
			requestBody: models.DataRecordRequest{
				Type:     models.TEXT,
				Data:     "test:data",
				Checksum: "94ee059335e587e501cc4bf90613e081",
				Name:     "Test Record",
			},
			expectedStatus: http.StatusCreated,
		},
		{
			name: "Invalid checksum",
			requestBody: models.DataRecordRequest{
				Type:     models.TEXT,
				Data:     "test:data",
				Checksum: "invalid_checksum",
				Name:     "Test Record",
			},
			expectedStatus: http.StatusBadRequest,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			reqBody, _ := json.Marshal(tc.requestBody)
			req, _ := http.NewRequest("POST", "/data-record", bytes.NewBuffer(reqBody))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)
			c.Request = req
			c.Set("userID", uint64(1))
			app.PutDataRecord(c)
			assert.Equal(t, tc.expectedStatus, w.Code)
		})
	}
}

func TestGetDataRecord(t *testing.T) {
	l, _ := logger.NewLogger()
	c, err := config.Load(l)
	if err != nil {
		fmt.Errorf("failed config load: %w", err)
		return
	}
	conn, err := store.NewPostgresConnection(c)
	if err != nil {
		fmt.Errorf("failed new postgres connection: %w", err)
		return
	}
	app := NewApp(&config.ServerConfig{}, store.NewStore(conn), l)
	testCases := []struct {
		name           string
		requestPath    string
		userID         uint64
		expectedStatus int
	}{
		{
			name:           "Successful data record retrieval",
			requestPath:    "/testrecord",
			userID:         1,
			expectedStatus: http.StatusOK,
		},
		{
			name:           "Unauthorized access",
			requestPath:    "/testrecord",
			userID:         0,
			expectedStatus: http.StatusUnauthorized,
		},
		{
			name:           "Record not found",
			requestPath:    "/nonexistentrecord",
			userID:         1,
			expectedStatus: http.StatusNotFound,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			req, _ := http.NewRequest("GET", tc.requestPath, nil)
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)
			c.Request = req
			c.Set("userID", tc.userID)
			app.GetDataRecord(c)
			assert.Equal(t, tc.expectedStatus, w.Code)
		})
	}
}

func TestGetDataRecords(t *testing.T) {
	l, _ := logger.NewLogger()
	c, err := config.Load(l)
	if err != nil {
		fmt.Errorf("failed config load: %w", err)
		return
	}
	conn, err := store.NewPostgresConnection(c)
	if err != nil {
		fmt.Errorf("failed new postgres connection: %w", err)
		return
	}
	app := NewApp(&config.ServerConfig{}, store.NewStore(conn), l)
	testCases := []struct {
		name           string
		userID         uint64
		expectedStatus int
	}{
		{
			name:           "Successful data records retrieval",
			userID:         1,
			expectedStatus: http.StatusOK,
		},
		{
			name:           "Unauthorized access",
			userID:         0,
			expectedStatus: http.StatusUnauthorized,
		},
		{
			name:           "No data found",
			userID:         2,
			expectedStatus: http.StatusNoContent,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			req, _ := http.NewRequest("GET", "/list", nil)
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)
			c.Request = req
			c.Set("userID", tc.userID)
			app.GetDataRecords(c)
			assert.Equal(t, tc.expectedStatus, w.Code)
		})
	}
}
