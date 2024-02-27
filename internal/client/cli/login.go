// Модуль логина
package cli

import (
	"context"
	"errors"
	"fmt"
	"github.com/EvgeniyBudaev/gophkeeper/internal/client/httpClient"
	"github.com/EvgeniyBudaev/gophkeeper/internal/client/logic"
	"github.com/EvgeniyBudaev/gophkeeper/internal/client/utils"
	"github.com/EvgeniyBudaev/gophkeeper/internal/server/logger"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"log"
	"net"
	"time"
)

// init - создаем команду логина
func init() {
	rootCmd.AddCommand(loginCmd)
}

var loginCmd = &cobra.Command{
	Use:   "login",
	Short: "Login to gophkeeper",
	Run: func(cmd *cobra.Command, args []string) {
		logger, err := logger.NewLogger()
		if err != nil {
			log.Fatal(err)
		}
		Login(context.Background(), logger.Named("login"))
	},
}

// Login - запуск логина
func Login(ctx context.Context, logger *zap.SugaredLogger) {
	for {
		token := viper.GetString("token")
		if token == "" {
			logger.Infoln("Login:")
			var login string
			fmt.Scanln(&login)
			logger.Infoln("Password:")
			var password string
			fmt.Scanln(&password)
			httpclient := httpClient.GetHTTPClient()
			creds, err := logic.Login(ctx, httpclient, login, password)
			if err != nil {
				var target *net.OpError
				if errors.As(err, &target) {
					if err := utils.CreateUsersDir(login); err != nil {
						logger.Errorf("err: %w", err)
					}
					logger.Infof("created local dir for user: %s\n", login)
				}
				return
			}
			viper.Set("login", login)
			viper.Set("token", creds.Token)
			viper.Set("expires_at", time.Now().Add(time.Duration(creds.ExpiresIn)*time.Second))
			if err := viper.WriteConfigAs("./gophkeeper.json"); err != nil {
				logger.Errorf("err saving config: %w", err)
			}
			if err := utils.CreateUsersDir(login); err != nil {
				logger.Errorf("err: %w", err)
			}
			return
		}
		viper.Set("token", "")
	}
}
