// Модуль хранилища
package store

import (
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"errors"
	"fmt"
	"github.com/EvgeniyBudaev/gophkeeper/internal/config"
	"github.com/EvgeniyBudaev/gophkeeper/internal/models"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgconn"
	_ "github.com/lib/pq"
)

// DBStore - хранилище данных
type DBStore struct {
	conn *sql.DB
}

// Store - интерфейс хранилища
type Store interface {
	CreateUser(c *gin.Context, user *models.User) (uint64, error)
	GetUser(c *gin.Context, u *models.User) (*models.User, error)
	PutDataRecord(c *gin.Context, data *models.DataRecord, userID uint64) error
	GetUserRecord(c *gin.Context, recordName string, userID uint64) (*models.DataRecord, error)
	GetUserRecords(c *gin.Context, userID uint64) ([]models.DataRecord, error)
}

var ErrLoginNotFound = errors.New("login not found")
var ErrDuplicateLogin = errors.New("login already registered")

// NewDBStore - создание хранилища данных
func NewStore(conn *sql.DB) Store {
	return &DBStore{conn: conn}
}

// NewPostgresConnection - создание подключения к PostgreSQL
func NewPostgresConnection(c *config.ServerConfig) (*sql.DB, error) {
	return sql.Open("postgres", c.DatabaseDSN)
}

func (db *DBStore) CreateUser(c *gin.Context, u *models.User) (uint64, error) {
	hashedPassword := hashPassword(u.Password)
	query := `INSERT INTO users (login, password) VALUES ($1, $2) RETURNING id`
	err := db.conn.QueryRowContext(c, query, &u.Login, hashedPassword).Scan(&u.ID)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			if pgErr.Code == "23505" {
				return 0, fmt.Errorf("error saving user to db: %w", err)
			}
		}
	}
	return u.ID, err
}

// GetUser - получение пользователя
func (db *DBStore) GetUser(c *gin.Context, u *models.User) (*models.User, error) {
	user := models.User{}
	query := `SELECT id, login FROM users WHERE login = $1`
	err := db.conn.QueryRowContext(c, query, u.Login).Scan(&user.ID, &user.Login)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("error user not found in db: %w", err)
		}
		return nil, err
	}
	return &user, nil
}

// PutDataRecord - сохранение данных
func (db *DBStore) PutDataRecord(c *gin.Context, data *models.DataRecord, userID uint64) error {
	query := `
		UPDATE data_records 
		SET uploaded_at=$1, type=$2, checksum=$3, data=$4, filepath=$5, name=$6, user_id=$7, blocked=$8
		WHERE user_id=$9
	`
	_, err := db.conn.ExecContext(c, query, &data.UploadedAt, &data.Type, &data.Checksum, &data.Data, &data.FilePath,
		&data.Name, &data.UserID, &data.Blocked, &userID)
	if err != nil {
		return fmt.Errorf("error saving data: %w", err)
	}
	return nil
}

// GetUserRecord- получение данных по названию записи и ID пользователя
func (db *DBStore) GetUserRecord(c *gin.Context, recordName string, userID uint64) (*models.DataRecord, error) {
	record := models.DataRecord{}
	query := `SELECT id, uploaded_at, type, checksum, data, filepath, name, user_id, blocked
              FROM data_records
              WHERE user_id=$1 AND name=$2`
	row := db.conn.QueryRowContext(c, query, userID, recordName)
	if row == nil {
		return nil, fmt.Errorf("no rows found")
	}
	err := row.Scan(&record.ID, &record.UploadedAt, &record.Type, &record.Checksum, &record.Data, &record.FilePath,
		&record.Name, &record.UserID, &record.Blocked)
	if err != nil {
		return nil, fmt.Errorf("error getting order: %w", err)
	}
	return &record, nil
}

// GetUserRecords - получение всех записей пользователя
func (db *DBStore) GetUserRecords(c *gin.Context, userID uint64) ([]models.DataRecord, error) {
	records := make([]models.DataRecord, 0)
	query := `SELECT id, uploaded_at, type, checksum, data, filepath, name, user_id, blocked
              FROM data_records
              WHERE user_id=$1`
	rows, err := db.conn.QueryContext(c, query, userID)
	if err != nil {
		return nil, fmt.Errorf("error getting all user records: %w", err)
	}
	defer rows.Close()
	for rows.Next() {
		record := models.DataRecord{}
		err := rows.Scan(&record.ID, &record.UploadedAt, &record.Type, &record.Checksum, &record.Data, &record.FilePath,
			&record.Name, &record.UserID, &record.Blocked)
		if err != nil {
			return nil, fmt.Errorf("error getting user record: %w", err)
		}
		records = append(records, record)
	}
	if len(records) == 0 {
		return nil, models.ErrNoData
	}
	return records, nil
}

// hashPassword — вспомогательная функция для хеширования пароля с использованием SHA-256.
func hashPassword(password string) string {
	hash := sha256.Sum256([]byte(password))
	return hex.EncodeToString(hash[:])
}
