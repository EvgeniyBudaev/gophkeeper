// Модуль регистрации
package logic

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/EvgeniyBudaev/gophkeeper/cmd/client/internal/client"
	"github.com/EvgeniyBudaev/gophkeeper/internal/models"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"net/http"
	"net/url"
)

// Register - регистрация нового пользователя
func Register(logger *zap.SugaredLogger, login string, password string) (creds *models.TokenResponse, err error) {
	httpclient := client.GetHTTPClient()
	if httpclient == nil {
		return nil, fmt.Errorf("configuration error")
	}
	endpoint, _ := url.JoinPath(httpclient.APIURL, "api/user/register")
	b, _ := json.Marshal(LoginReq{Login: login, Password: password})
	request, err := http.NewRequest(http.MethodPost, endpoint, bytes.NewBuffer(b))
	if err != nil {
		return nil, fmt.Errorf("error: %w\n", err)
	}
	request.Header.Add("Content-Type", "application/json")
	response, err := httpclient.Do(request)
	if err != nil {
		return nil, fmt.Errorf("error: %w\n", err)
	}
	defer func() {
		if err := response.Body.Close(); err != nil {
			logger.Errorf("error: %w\n", err)
		}
	}()
	if response.StatusCode != http.StatusCreated {
		return nil, fmt.Errorf("Error in Register")
	}
	creds = &models.TokenResponse{}
	if err = json.NewDecoder(response.Body).Decode(creds); err != nil {
		return nil, fmt.Errorf("error decode body: %w\n", err)
	}
	viper.Set("api", httpclient.APIURL)
	return creds, nil
}
