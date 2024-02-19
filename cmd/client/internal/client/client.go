// Модуль клиента
package client

import (
	"crypto/tls"
	"github.com/EvgeniyBudaev/gophkeeper/internal/logger"
	"github.com/spf13/viper"
	"log"
	"net/http"
	"sync"
)

// httpClientInstance - синглетон клиента
type httpClientInstance struct {
	*http.Client
	APIURL string
}

var (
	httpClient *httpClientInstance
	once       sync.Once
)

// GetHTTPClient - возвращает клиент
func GetHTTPClient() *httpClientInstance {
	once.Do(
		func() {
			l, err := logger.NewLogger()
			if err != nil {
				log.Fatal(err)
			}
			apiURL := viper.GetString("api")
			if apiURL == "" {
				l.Errorln("empty API URL")
				httpClient = nil
				return
			}
			httpClient = &httpClientInstance{
				Client: &http.Client{
					Transport: &http.Transport{
						TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
					}},
				APIURL: apiURL,
			}
		})
	return httpClient
}
