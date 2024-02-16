package env

import (
	"os"
	"strconv"
)

func Getenv(key, fallback string) string {
	value, exists := os.LookupEnv(key)
	if !exists {
		return fallback
	}
	return value
}

func GetenvInt(key string, fallback int) int {
	value, exists := os.LookupEnv(key)
	if !exists {
		return fallback
	}
	valueInt, err := strconv.Atoi(value)
	if err != nil {
		return fallback
	}
	return valueInt
}
