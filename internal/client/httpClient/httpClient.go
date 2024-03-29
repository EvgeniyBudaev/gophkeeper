// Модуль клиента
package httpClient

import (
	"crypto/tls"
	"github.com/EvgeniyBudaev/gophkeeper/internal/server/logger"
	"github.com/spf13/viper"
	"log"
	"net/http"
	"sync"
)

// httpClientInstance - синглетон клиента
type HttpClientInstance struct {
	*http.Client
	APIURL string
}

var (
	httpClient *HttpClientInstance
	once       sync.Once
)

// GetHTTPClient - возвращает клиент
func GetHTTPClient() *HttpClientInstance {
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
			httpClient = &HttpClientInstance{
				Client: &http.Client{
					Transport: &http.Transport{
						TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
					}},
				APIURL: apiURL,
			}
		})
	return httpClient
}
