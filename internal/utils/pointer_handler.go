package utils

import "time"

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
