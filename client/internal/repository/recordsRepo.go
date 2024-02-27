package repository

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/EvgeniyBudaev/gophkeeper/internal/models"
	"github.com/EvgeniyBudaev/gophkeeper/internal/utils"
	"os"
	"path/filepath"
)

type DataRecordRepository interface {
	Add(data *models.DataRecord) error
}

type PassRepository struct {
	login string
}

func NewPassRepository(login string) *PassRepository {
	return &PassRepository{login: login}
}

func (r *PassRepository) Add(data *models.DataRecord) error {
	ext := utils.GetExtension(data.Type)
	filename := fmt.Sprintf("%s%s", data.Name, ext)
	filepath := filepath.Join(".", r.login, filename)
	localFile, err := os.OpenFile(filepath, os.O_RDONLY, 0600)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			wrLocalFile, err := os.OpenFile(filepath, os.O_CREATE|os.O_RDWR|os.O_APPEND, 0600)
			if err != nil {
				return err
			}
			defer wrLocalFile.Close()
			if err := json.NewEncoder(wrLocalFile).Encode(data); err != nil {
				return err
			}
			return nil
		}
		return err
	}
	var localData models.DataRecord
	if err := json.NewDecoder(localFile).Decode(&localData); err != nil {
		return err
	}
	if err := localFile.Close(); err != nil {
		return err
	}
	wrLocalFile, err := os.OpenFile(filepath, os.O_RDWR, 0600)
	if data.ID != 0 && localData.ID == 0 {
		wrLocalFile.Truncate(0)
		if err := json.NewEncoder(wrLocalFile).Encode(data); err != nil {
			return err
		}
	}
	return nil
}
