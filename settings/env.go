package settings

import (
	"fmt"
	"os"
)

func GETENV(key string) string {
	env, isPresent := os.LookupEnv(key)

	if !isPresent {
		panic(fmt.Sprintf("%s is not present in environment variables", key))
	}
	return env
}

func GetEnvOrDefault(key string, defaultValue string) string {
	env, isPresent := os.LookupEnv(key)

	if !isPresent {
		return defaultValue
	}
	return env
}
