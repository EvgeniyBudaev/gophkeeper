// Модуль расширяет функциональными возможностями
package utils

import "github.com/EvgeniyBudaev/gophkeeper/internal/models"

// GetExtension - возвращает расширение
func GetExtension(dataType models.DataType) string {
	switch dataType {
	case models.PASS:
		return ".json"
	case models.TEXT:
		return ".json"
	default:
		return ".json"
	}
}
