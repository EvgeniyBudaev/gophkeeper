// Модуль пользователя
package models

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
)

var (
	ErrNoData = errors.New("no data")
)

// User - модель пользователя
type User struct {
	Login    string `json:"login"`
	Password string `json:"-"`
	ID       uint64 `json:"id,omitempty"`
}

// UserCredentialsSchema - структура для хранения данных пользователя
type UserCredentialsSchema struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

// TokenResponse - ответ сервера
type TokenResponse struct {
	Token     string `json:"token"`
	ExpiresIn int    `json:"expires_in"`
}

// GetUserFolder - получить путь к папке пользователя
func (u *User) GetUserFolder() ([]fs.DirEntry, error) {
	return os.ReadDir(fmt.Sprintf("./userdata/%s-%d", u.Login, u.ID))
}
