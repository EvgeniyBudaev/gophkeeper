// Модуль регистрации
package cli

import (
	"context"
	"errors"
	"fmt"
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

// init представляет команду инициализации
func init() {
	rootCmd.AddCommand(registerCmd)
}

var registerCmd = &cobra.Command{
	Use:   "register",
	Short: "Register to gophkeeper",
	Run: func(cmd *cobra.Command, args []string) {
		logger, err := logger.NewLogger()
		if err != nil {
			log.Fatal(err)
		}
		Register(context.Background(), logger)
	},
}

// Register запускает процесс регистрации
func Register(ctx context.Context, logger *zap.SugaredLogger) {
	logger.Infoln("Login:")
	var login string
	fmt.Scanln(&login)
	logger.Infoln("Password:")
	var password string
	fmt.Scanln(&password)
	creds, err := logic.Register(logger, login, password)
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
	if err := utils.CreateUsersDir(login); err != nil {
		logger.Errorf("err: %w", err)
	}
	if err := viper.WriteConfigAs("./gophkeeper.json"); err != nil {
		logger.Errorf("err saving config: %w", err)
	}
}
