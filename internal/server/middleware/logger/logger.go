// Модуль логирования
package logger

import (
	"bytes"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"io"
	"time"
)

// Logger - логгер
func Logger(logger *zap.SugaredLogger) (gin.HandlerFunc, error) {
	return func(c *gin.Context) {
		uri := c.Request.RequestURI
		method := c.Request.Method
		body, err := io.ReadAll(c.Request.Body)
		if err != nil {
			return
		}
		c.Request.Body = io.NopCloser(bytes.NewReader(body))
		t := time.Now()
		c.Next()
		duration := time.Since(t)
		logger.Infoln(
			"URI", uri,
			"Method", method,
			"Duration", duration,
			"Status", c.Writer.Status(),
			"Size", c.Writer.Size(),
		)
		logger.Infoln("Data", string(body))
	}, nil
}
