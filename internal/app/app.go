// Модуль приложения
package app

import (
	"bytes"
	"crypto/md5"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/EvgeniyBudaev/gophkeeper/internal/adapters/store"
	"github.com/EvgeniyBudaev/gophkeeper/internal/config"
	"github.com/EvgeniyBudaev/gophkeeper/internal/middleware/auth"
	"github.com/EvgeniyBudaev/gophkeeper/internal/models"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"net/http"
	"os"
	"time"
)

// App - структура приложения
type App struct {
	config *config.ServerConfig
	store  store.Store
	logger *zap.SugaredLogger
}

const (
	bcryptCost   = 7
	maxCookieAge = 3600 * 24 * 30
)

// NewApp - конструктор приложения
func NewApp(config *config.ServerConfig, store store.Store, logger *zap.SugaredLogger) *App {
	return &App{
		config: config,
		store:  store,
		logger: logger,
	}
}

// NewServer - конструктор сервера
func (a *App) NewServer() (*http.Server, error) {
	r, err := a.SetupRouter()
	if err != nil {
		return nil, fmt.Errorf("error init router: %w", err)
	}
	return &http.Server{
		Addr:    a.config.RunAddr,
		Handler: r,
	}, nil
}

// RecordNotFoundError - ошибка запись не найдена в БД
type RecordNotFoundError struct {
	Message string
}

// Error - возвращает ошибку
func (e *RecordNotFoundError) Error() string {
	return e.Message
}

// Login - логин пользователя
func (a *App) Login(c *gin.Context) {
	a.logger.Info("/api/user/login")
	req := c.Request
	res := c.Writer
	userCreds := models.User{}
	if err := json.NewDecoder(req.Body).Decode(&userCreds); err != nil {
		a.logger.Debug("user credentials cannot be decoded: %v", zap.Error(err))
		res.WriteHeader(http.StatusBadRequest)
		return
	}
	userReq := models.User{
		Login:    userCreds.Login,
		Password: userCreds.Password,
	}
	u, err := a.store.GetUser(c, &models.User{Login: userReq.Login})
	if err != nil {
		if errors.Is(err, store.ErrLoginNotFound) {
			a.logger.Debug("login not found: %v", zap.Error(err))
			res.WriteHeader(http.StatusUnauthorized)
			return
		} else {
			a.logger.Debug("cannot get user: %v", zap.Error(err))
			res.WriteHeader(http.StatusInternalServerError)
			return
		}
	}
	ok := verifyPassword(userReq.Password, u.Password)
	if !ok {
		a.logger.Debug("cannot verifyPassword: %v", zap.Error(err))
		res.WriteHeader(http.StatusUnauthorized)
		return
	}
	userReq.ID = u.ID
	jwt, err := auth.BuildJWTString(userReq.ID)
	if err != nil {
		a.logger.Debug("cannot build jwt string for authorized user: %v", zap.Error(err))
		res.WriteHeader(http.StatusInternalServerError)
		return
	}
	c.JSON(http.StatusOK, models.TokenResponse{
		Token:     jwt,
		ExpiresIn: maxCookieAge,
	})
}

// Register - регистрация пользователя
func (a *App) Register(c *gin.Context) {
	a.logger.Info("/api/user/register")
	req := c.Request
	res := c.Writer
	userCreds := models.User{}
	if err := json.NewDecoder(req.Body).Decode(&userCreds); err != nil {
		a.logger.Debug("body cannot be decoded: %v", zap.Error(err))
		res.WriteHeader(http.StatusBadRequest)
		return
	}
	userReq := models.User{
		Login:    userCreds.Login,
		Password: userCreds.Password,
	}
	if _, err := a.store.CreateUser(c, &userReq); err != nil {
		if errors.Is(err, store.ErrDuplicateLogin) {
			a.logger.Debug("login already taken: %v", zap.Error(err))
			res.WriteHeader(http.StatusConflict)
			return
		} else {
			a.logger.Debug("cannot operate user creds: %v", zap.Error(err))
			res.WriteHeader(http.StatusInternalServerError)
			return
		}
	}
	if err := os.MkdirAll(fmt.Sprintf("./userdata/%s-%d/", userReq.Login, userReq.ID), 0700); err != nil {
		a.logger.Debug("cannot create user folder: %v", zap.Error(err))
		res.WriteHeader(http.StatusInternalServerError)
		return
	}
	jwt, err := auth.BuildJWTString(userReq.ID)
	if err != nil {
		a.logger.Debug("cannot build jwt string: %v", zap.Error(err))
		res.WriteHeader(http.StatusInternalServerError)
		return
	}
	c.JSON(http.StatusCreated, models.TokenResponse{
		Token:     jwt,
		ExpiresIn: maxCookieAge,
	})
}

// PutDataRecord - запись данных
func (a *App) PutDataRecord(c *gin.Context) {
	a.logger.Info("/")
	userID := c.GetUint64(auth.UserIDKey.ToString())
	req := c.Request
	res := c.Writer
	if userID == 0 {
		a.logger.Debug("user unauthorized")
		res.WriteHeader(http.StatusUnauthorized)
		return
	}
	var record models.DataRecordRequest
	if err := json.NewDecoder(req.Body).Decode(&record); err != nil {
		a.logger.Debug("cannot decode body: %w", zap.Error(err))
		res.WriteHeader(http.StatusBadRequest)
		return
	}
	parts := bytes.Split([]byte(record.Data), []byte(":"))
	if len(parts) <= 1 {
		a.logger.Debug("cannot parts <= 1")
		res.WriteHeader(http.StatusBadRequest)
		return
	}
	if record.Type == models.PASS || record.Type == models.TEXT {
		checksum := fmt.Sprintf("%x", md5.Sum([]byte(record.Data)))
		if record.Checksum != checksum {
			a.logger.Debug("wrong checksum from request, corrupted data")
			res.WriteHeader(http.StatusBadRequest)
			return
		}
	}
	data := &models.DataRecord{
		UploadedAt: time.Now(),
		Type:       record.Type,
		Checksum:   fmt.Sprintf("%x", md5.Sum([]byte(record.Data))),
		Data:       record.Data,
		FilePath:   "",
		UserID:     userID,
		Name:       record.Name,
		Blocked:    false,
	}
	if record.ID != 0 {
		data.ID = record.ID
	}
	if err := a.store.PutDataRecord(c, data); err != nil {
		a.logger.Debug("unhandled error: %v", zap.Error(err))
		res.WriteHeader(http.StatusInternalServerError)
		return
	}
	c.JSON(http.StatusCreated, data)
}

// GetDataRecord - получение записи
func (a *App) GetDataRecord(c *gin.Context) {
	a.logger.Info("/:name")
	res := c.Writer
	recordName := c.Param("name")
	userID := c.GetUint64(auth.UserIDKey.ToString())

	if userID == 0 {
		a.logger.Debug("user unauthorized")
		res.WriteHeader(http.StatusUnauthorized)
		return
	}
	fmt.Println("userID:", userID)
	fmt.Println("recordName:", recordName)
	record, err := a.store.GetUserRecord(c, recordName, userID)
	if err != nil {
		var recordNotFoundError *RecordNotFoundError
		if errors.As(err, &recordNotFoundError) {
			res.WriteHeader(http.StatusNotFound)
			return
		}
		a.logger.Debug("error getting user record: %v", zap.Error(err))
		res.WriteHeader(http.StatusInternalServerError)
		return
	}
	c.JSON(http.StatusOK, record)
}

// GetDataRecords - получение записей пользователя
func (a *App) GetDataRecords(c *gin.Context) {
	a.logger.Info("/list")
	userID := c.GetUint64(auth.UserIDKey.ToString())
	res := c.Writer
	if userID == 0 {
		res.WriteHeader(http.StatusUnauthorized)
		return
	}
	records, err := a.store.GetUserRecords(c, userID)
	if err != nil {
		if errors.Is(err, models.ErrNoData) {
			res.WriteHeader(http.StatusNoContent)
			return
		}
		a.logger.Debug("error getting user records: %v", zap.Error(err))
		res.WriteHeader(http.StatusInternalServerError)
		return
	}
	c.JSON(http.StatusOK, records)
}

func verifyPassword(inputPassword, storedHash string) bool {
	hash := sha256.Sum256([]byte(inputPassword))
	hashedInput := hex.EncodeToString(hash[:])
	return hashedInput == storedHash
}
