package utils

import (
	"os"
	"strings"
)

func GetEnvAsBool(key string, defaultVal bool) bool {
	valStr := os.Getenv(key)
	if valStr == "" {
		return defaultVal
	}

	valStr = strings.ToLower(valStr)
	return valStr == "true" || valStr == "1" || valStr == "yes" || valStr == "on"
}