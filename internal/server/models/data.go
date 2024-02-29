// Модуль данных
package models

import (
	"database/sql/driver"
	"errors"
	"fmt"
	"time"
)

type DataType string

const (
	PASS DataType = "PASS"
	TEXT DataType = "TEXT"
)

// Scan - реализация интерфейса sql.Scanner
func (s *DataType) Scan(value interface{}) error {
	sv, ok := value.(string)
	if !ok {
		return errors.New(fmt.Sprint("Failed to unmarshal DataType value: ", value))
	}
	*s = DataType(sv)
	return nil
}

// Value - реализация интерфейса driver.Valuer
func (s DataType) Value() (driver.Value, error) {
	return string(s), nil
}

// DataRecord - структура данных
type DataRecord struct {
	ID         uint64    `json:"id"`
	UploadedAt time.Time `json:"uploaded_at"`
	Type       DataType  `json:"type"`
	Checksum   string    `json:"checksum"`
	Data       string    `json:"data"`
	FilePath   string    `json:"filepath"`
	Name       string    `json:"name"`
	User       User      `json:"-"`
	UserID     uint64    `json:"-"`
	Key        string    `json:"key"`
}

// DataRecordRequest - структура данных запроса
type DataRecordRequest struct {
	Type     DataType `json:"type"`
	Checksum string   `json:"checksum"`
	Data     string   `json:"data"`
	Name     string   `json:"name"`
	ID       uint64   `json:"id"`
	Key      string   `json:"key"`
}
