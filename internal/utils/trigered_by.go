package utils

import (
	"os"
	"strconv"
	"strings"
)

/* Catatan penting: kalau argumen mengandung spasi, harus kamu quote saat jalankan program.
go run main.go "hello world" 123 */

func TrigeredBy() *int {
	for _, arg := range os.Args[1:] { // skip arg pertama (path program)
		if strings.HasPrefix(strings.ToLower(arg), "--user=") {
			parts := strings.SplitN(arg, "=", 2)
			if len(parts) == 2 {
				if v, err := strconv.Atoi(parts[1]); err == nil {
					return &v
				}
			}
		}
	}

	// fallback ke default
	defaultVal := 1
	return &defaultVal
}
