// Модуль для работы с роутером
package app

import (
	"fmt"
	"github.com/EvgeniyBudaev/gophkeeper/internal/middleware/auth"
	ginLogger "github.com/EvgeniyBudaev/gophkeeper/internal/middleware/logger"
	"github.com/gin-gonic/gin"
)

const (
	rootRoute    = "/"
	userAPIRoute = "/api/user"
)

// SetupRouter Инициализация роутера
func (a *App) SetupRouter() (*gin.Engine, error) {
	r := gin.New()
	ginLoggerMiddleware, err := ginLogger.Logger(a.logger)
	if err != nil {
		return nil, fmt.Errorf("error creating middleware logger func: %w", err)
	}
	r.Use(ginLoggerMiddleware)
	userAPI := r.Group(userAPIRoute)
	{
		userAPI.POST("register", a.Register)
		userAPI.POST("login", a.Login)
		recordsAPI := userAPI.Group("records")
		recordsAPI.Use(auth.AuthMiddleware(a.logger))
		{
			recordsAPI.POST(rootRoute, a.PutDataRecord)
			recordsAPI.GET("list", a.GetDataRecords)
			recordsAPI.GET(":name", a.GetDataRecord)
		}
	}
	return r, nil
}
