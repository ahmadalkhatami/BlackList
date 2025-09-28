package utils

import "time"

func Ptr[T any](v T) *T {
	return &v
}

// StringPtr mengembalikan pointer ke string
func StringPtr(s string) *string {
	return &s
}

// IntPtr mengembalikan pointer ke int
func IntPtr(i int) *int {
	return &i
}

// Float64Ptr mengembalikan pointer ke float64
func Float64Ptr(f float64) *float64 {
	return &f
}

// TimePtr mengembalikan pointer ke time.Time
func TimePtr(t time.Time) *time.Time {
	return &t
}

func ValStr(ptr *string) string {
	if ptr == nil {
		return "-"
	}
	if *ptr == "" {
		return "-"
	}
	return *ptr
}

func ValInt64(ptr *int64) int64 {
	if ptr == nil {
		return 0
	}
	return *ptr
}

func ValFloat64(ptr *float64) float64 {
	if ptr == nil {
		return 0
	}
	return *ptr
}

func ValTime(ptr *time.Time) string {
	if ptr == nil {
		return "-"
	}
	// format Indonesia DD/MM/YYYY
	return ptr.Format("02/01/2006")
}

var layouts = []string{
	"2006-01-02",      // 1996-12-22
	"02/01/2006",      // 22/12/1996
	"02-01-2006",      // 22-12-1996
	"2006/01/02",      // 1996/12/22
	"02 Jan 2006",     // 22 Dec 1996
	"January 2, 2006", // December 22, 1996
	"02-Jan-2006",     // 22-Dec-1996
}

func ParseDatePtrAny(value string) *time.Time {
	for _, layout := range layouts {
		if t, err := time.Parse(layout, value); err == nil {
			return &t
		}
	}
	return nil // gagal parse semua layout
}

func PtrOrEmpty(s string) *string {
	if s == "" {
		return nil
	}
	return &s
}

func Max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
