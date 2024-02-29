// Модуль утилиты для работы с пользователями
package utils

import (
	"fmt"
	"os"
	"path"
)

// CreateUsersDir - создание директории пользователей
func CreateUsersDir(username string) error {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("error getting user's home directory: %v", err)
	}
	userDir := path.Join(homeDir, username, "."+"gophkeeper")
	if err := os.MkdirAll(userDir, 0750); err != nil {
		return fmt.Errorf("error creating user's dir: %v", err)
	}
	return nil
}
