package env

import (
	"fmt"
	"os"
	"strconv"
	"time"
)

func GetInt(key string, defaultValue int) int {
	value, exists := os.LookupEnv(key)
	if !exists {
		return defaultValue
	}

	intValue, err := strconv.Atoi(value)
	if err != nil {
		panic(fmt.Errorf("environment variable %s=%q cannot be converted to an int", key, value))
	}
	return intValue
}

func GetString(key, defaultVal string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return defaultVal
}

func GetFloat(key string, defaultVal float64) float64 {
	if value, ok := os.LookupEnv(key); ok {
		f, err := strconv.ParseFloat(value, 64)
		if err != nil {
			panic(fmt.Errorf("environment variable %s=%q cannot be converted to a float", key, value))
		}
		return f
	}
	return defaultVal
}

func GetBool(key string, defaultVal bool) bool {
	if value, ok := os.LookupEnv(key); ok {
		b, err := strconv.ParseBool(value)
		if err != nil {
			panic(fmt.Errorf("environment variable %s=%q cannot be converted to a bool", key, value))
		}
		return b
	}
	return defaultVal
}

func GetDuration(key string, defaultValue time.Duration) time.Duration {
	value, exists := os.LookupEnv(key)
	if !exists {
		return defaultValue
	}

	durationValue, err := time.ParseDuration(value)
	if err != nil {
		panic(fmt.Errorf("environment variable %s=%q cannot be converted to a time.Duration", key, value))
	}
	return durationValue
}
