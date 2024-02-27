// Модуль хранения записей
package logic

import (
	"bytes"
	"context"
	"crypto/md5"
	"encoding/json"
	"fmt"
	"github.com/EvgeniyBudaev/gophkeeper/client/internal/httpClient"
	"github.com/EvgeniyBudaev/gophkeeper/client/internal/repository"
	"github.com/EvgeniyBudaev/gophkeeper/internal/models"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"golang.org/x/sync/errgroup"
	"net/http"
	"net/url"
	"os"
	"strings"
)

// SaveOrUpdateData - запись иди обновление данных
func SaveOrUpdateData(logger *zap.SugaredLogger, data *models.DataRecord) error {
	login := viper.GetString("login")
	if login == "" {
		err := fmt.Errorf("not logged in")
		logger.Error(err)
		return err
	}
	var repo repository.DataRecordRepository
	switch data.Type {
	case models.PASS:
		repo = repository.NewPassRepository(login)
	default:
		return fmt.Errorf("unsupported data type")
	}
	if err := repo.Add(data); err != nil {
		return err
	}
	return nil
}

// GetRecords - получение записей
func GetRecord(ctx context.Context, name string) (*models.DataRecord, error) {
	token := viper.GetString("token")
	if token == "" {
		return nil, fmt.Errorf("No auth data, login first")
	}
	httpclient := httpClient.GetHTTPClient()
	if httpclient == nil {
		return nil, fmt.Errorf("configuration error")
	}
	endpoint, _ := url.JoinPath(httpclient.APIURL, "api/user/records", name)
	request, err := http.NewRequestWithContext(ctx, http.MethodGet, endpoint, nil)
	if err != nil {
		return nil, err
	}
	request.Header.Add("Content-Type", "application/json")
	request.Header.Add("Authorization", fmt.Sprintf("Bearer %s", token))
	response, err := httpclient.Do(request)
	if err != nil {
		return nil, fmt.Errorf("error: %w", err)
	}
	if response.StatusCode == http.StatusNotFound {
		return nil, fmt.Errorf("record not found")
	}
	var record models.DataRecord
	if err = json.NewDecoder(response.Body).Decode(&record); err != nil {
		return nil, fmt.Errorf("error decode body: %w", err)
	}
	return &record, nil
}

// PutRecord - создание записи
func PutRecord(ctx context.Context, args []string) (*models.DataRecord, error) {
	if len(args) != 3 {
		return nil, fmt.Errorf("bad request")
	}
	dataType := strings.ToLower(args[0])
	var data string
	switch dataType {
	case "pass":
		data = args[1]
	default:
		path := args[1]
		fi, err := os.Stat(path)
		if err != nil {
			return nil, err
		}
		data = fi.Name()
	}
	token := viper.GetString("token")
	if token == "" {
		return nil, fmt.Errorf("No auth data, login first")
	}
	httpclient := httpClient.GetHTTPClient()
	if httpclient == nil {
		return nil, fmt.Errorf("configuration error")
	}
	endpoint, _ := url.JoinPath(httpclient.APIURL, "api/user/records")
	checksum := fmt.Sprintf("%x", md5.Sum([]byte(data)))
	dataObj := models.DataRecordRequest{
		Type:     models.DataType(dataType),
		Name:     args[2],
		Data:     data,
		Checksum: checksum,
	}
	dataObjB, err := json.Marshal(dataObj)
	if err != nil {
		return nil, err
	}
	request, err := http.NewRequestWithContext(ctx, http.MethodPost, endpoint, bytes.NewBuffer(dataObjB))
	if err != nil {
		return nil, err
	}
	request.Header.Add("Content-Type", "application/json")
	request.Header.Add("Authorization", fmt.Sprintf("Bearer %s", token))
	response, err := httpclient.Do(request)
	if err != nil {
		return &models.DataRecord{
			Data:     dataObj.Data,
			Checksum: dataObj.Checksum,
			Type:     dataObj.Type,
			Name:     dataObj.Name,
		}, err
	}
	if response.StatusCode != http.StatusCreated {
		return nil, fmt.Errorf("error in Post data")
	}
	var record models.DataRecord
	if err = json.NewDecoder(response.Body).Decode(&record); err != nil {
		return nil, fmt.Errorf("error decode body: %w", err)
	}
	return &record, nil
}

// ListRecords - получение списка записей
func ListRecords(ctx context.Context, logger *zap.SugaredLogger) ([]models.DataRecord, error) {
	token := viper.GetString("token")
	if token == "" {
		err := fmt.Errorf("no auth data, login first")
		logger.Error(err)
		return nil, err
	}
	httpclient := httpClient.GetHTTPClient()
	if httpclient == nil {
		err := fmt.Errorf("configuration error")
		logger.Error(err)
		return nil, err
	}
	endpoint, _ := url.JoinPath(httpclient.APIURL, "api/user/records/list")
	request, err := http.NewRequestWithContext(ctx, http.MethodGet, endpoint, nil)
	if err != nil {
		logger.Error(err)
		return nil, err
	}
	request.Header.Add("Content-Type", "application/json")
	request.Header.Add("Authorization", fmt.Sprintf("Bearer %s", token))
	response, err := httpclient.Do(request)
	if err != nil {
		return nil, err
	}
	if response.StatusCode == http.StatusNoContent {
		logger.Infoln("no records found")
		return nil, nil
	}
	if response.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("error in listrecords\n")
	}
	records := make([]models.DataRecord, 0)
	if err = json.NewDecoder(response.Body).Decode(&records); err != nil {
		return nil, fmt.Errorf("error decode body: %w\n", err)
	}
	return records, nil
}

// SyncDataRecords - синхронизация данных
func SyncDataRecords(ctx context.Context, logger *zap.SugaredLogger) error {
	records, err := ListRecords(ctx, logger)
	if err != nil {
		return err
	}
	g := new(errgroup.Group)
	for _, r := range records {
		data := r
		g.Go(func() error {
			if err := SaveOrUpdateData(logger, &data); err != nil {
				return err
			}

			return nil
		})
	}
	if err := g.Wait(); err != nil {
		return err
	}
	return nil
}
