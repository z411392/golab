package utils

import (
	"os"
)

func OnDevelopment() bool {
	return os.Getenv("ENV") == "development"
}
