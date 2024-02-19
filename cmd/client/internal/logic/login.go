// Модуль логина
package logic

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/EvgeniyBudaev/gophkeeper/cmd/client/internal/client"
	"github.com/EvgeniyBudaev/gophkeeper/internal/models"
	"go.uber.org/zap"
	"net/http"
	"net/url"
)

// LoginReq - модель запроса логина
type LoginReq struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

// Login - логин
func Login(ctx context.Context, logger *zap.SugaredLogger, login string, password string) (creds *models.TokenResponse, err error) {
	httpclient := client.GetHTTPClient()
	if httpclient == nil {
		return nil, fmt.Errorf("configuration error")
	}
	endpoint, _ := url.JoinPath(httpclient.APIURL, "api/user/login")
	b, _ := json.Marshal(LoginReq{Login: login, Password: password})
	request, err := http.NewRequestWithContext(ctx, http.MethodPost, endpoint, bytes.NewBuffer(b))
	if err != nil {
		return nil, err
	}
	request.Header.Add("Content-Type", "application/json")
	response, err := httpclient.Do(request)
	if err != nil {
		return nil, err
	}
	defer func() {
		if err := response.Body.Close(); err != nil {
			logger.Errorf("error: %w", err)
		}
	}()
	if response.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Error in Login")
	}
	creds = &models.TokenResponse{}
	if err = json.NewDecoder(response.Body).Decode(creds); err != nil {
		return nil, fmt.Errorf("error decode body: %w", err)
	}
	return creds, nil
}
