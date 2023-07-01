package utils

import (
	"os"
)

func Get(name string, defaultValue string) string {
	value, ok := os.LookupEnv(name)
	if !ok {
		value = defaultValue
	}
	return value
}
