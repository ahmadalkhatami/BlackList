package utils

import (
	"os"
	"strings"
)

func TrigeredBy() string {
	for _, arg := range os.Args {

		if strings.HasPrefix(strings.ToLower(arg), "--user=") {

			parts := strings.SplitN(arg, "=", 2)
			if len(parts) == 2 {
				return parts[1]
			}
		}
	}
	return "1"
}
