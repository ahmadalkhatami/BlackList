package utils

import "reflect"

// MapOne memetakan 1 objek ke objek lain
func MapOne[T any, U any](in T, mapFn func(T) U) U {
	return mapFn(in)
}

// MapSlice memetakan slice objek ke slice objek lain
func MapSlice[T any, U any](in []T, mapFn func(T) U) []U {
	out := make([]U, len(in))
	for i, v := range in {
		out[i] = mapFn(v)
	}
	return out
}

// MapSlice2 menggabungkan 2 slice menjadi slice DTO baru
// Bisa dipakai untuk hasil join antar tabel/data
func MapSlice2[T any, U any, R any](list []T, lookup map[string]U, getKey func(T) string, mapper func(T, U) R) []R {
	result := make([]R, 0, len(list)) // Pre-allocate kapasitas agar performa lebih baik

	for _, item := range list {
		key := getKey(item)
		if val, ok := lookup[key]; ok {
			result = append(result, mapper(item, val))
		}
	}

	return result
}

// GetStructKeys mengembalikan slice berisi nama field struct yang public uppercase prefixes
func GetStructKeys(s interface{}) []string {
	if s == nil {
		return nil
	}

	t := reflect.TypeOf(s)

	// Kalau pointer, dereference
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}

	// Pastikan struct
	if t.Kind() != reflect.Struct {
		return nil
	}

	var keys []string
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		// Skip unexported field
		if field.PkgPath != "" {
			continue
		}
		keys = append(keys, field.Name)
	}

	return keys
}

func StrSliceToInterface(slice []string) []interface{} {
	out := make([]interface{}, len(slice))
	for i, v := range slice {
		out[i] = v
	}
	return out
}
