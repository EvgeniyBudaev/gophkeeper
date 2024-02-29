// Модуль логгера
package utils

const (
	err int = iota + 1
	warn
	info
	debug
)

// ConvertLogLevelToInt - конвертирует логгер в числовое значение
func ConvertLogLevelToInt(logLevel string) int {
	switch logLevel {
	case "debug":
		return debug
	case "info":
		return info
	case "warn":
		return warn
	case "error":
		return err
	default:
		return 0
	}
}
