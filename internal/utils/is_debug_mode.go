package utils

import (
	"os"
	"strings"
)

func IsDebugMode() bool {

	if strings.ToLower(os.Getenv("DEBUG")) == "true" {
		return true
	}

	for _, arg := range os.Args {
		if strings.ToLower(arg) == "--debug" {
			return true
		}
	}

	return false
}
