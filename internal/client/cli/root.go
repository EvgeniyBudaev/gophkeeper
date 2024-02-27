// Модуль root
package cli

import (
	"github.com/EvgeniyBudaev/gophkeeper/internal/server/logger"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"log"
	"os"
)

var (
	// Used for flags.
	cfgFile string
	apiURL  string
	rootCmd = &cobra.Command{
		Use:   "gclient",
		Short: "A generator for Cobra based Applications",
		Long: `Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	}
)

// Execute - запускает выполнение командной строки (CLI) приложения, созданного с использованием Cobra.
func Execute() error {
	return rootCmd.Execute()
}

// init - инициализация
func init() {
	cobra.OnInitialize(initConfig)
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.cobra.yaml)")
	rootCmd.PersistentFlags().StringVar(&apiURL, "api", "http://localhost:8080", "API URL")
	rootCmd.PersistentFlags().StringP("login", "l", "", "author name for copyright attribution")
	rootCmd.PersistentFlags().StringP("token", "t", "", "author name for copyright attribution")
	viper.BindPFlag("api", rootCmd.PersistentFlags().Lookup("api"))
	viper.BindPFlag("login", rootCmd.PersistentFlags().Lookup("login"))
	viper.BindPFlag("token", rootCmd.PersistentFlags().Lookup("token"))
	viper.BindPFlag("expires_at", rootCmd.PersistentFlags().Lookup("expires_at"))
}

// initConfig- инициализация конфига
func initConfig() {
	l, err := logger.NewLogger()
	if err != nil {
		log.Fatal(err)
	}
	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	} else {
		home, err := os.UserHomeDir()
		cobra.CheckErr(err)
		viper.AddConfigPath(home)
		viper.AddConfigPath(".")
		viper.SetConfigType("json")
		viper.SetConfigName("gophkeeper")
	}
	viper.AutomaticEnv()
	if err := viper.ReadInConfig(); err == nil {
		l.Info("Using config file:", viper.ConfigFileUsed())
	}
}
