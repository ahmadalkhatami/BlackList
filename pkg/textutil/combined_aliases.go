package textutil

import (
	"reflect"
	"strings"
)

func CombineAliases(v interface{}, prefixField string) []string {
	val := reflect.ValueOf(v)
	if !val.IsValid() {
		return nil
	}

	// Kalau pointer, dereference
	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}

	if val.Kind() != reflect.Struct {
		return nil
	}

	typ := val.Type()
	aliases := make([]string, 0)

	for i := 0; i < val.NumField(); i++ {
		field := typ.Field(i)

		// Skip field unexported
		if field.PkgPath != "" {
			continue
		}

		// Cek method getter dulu: Get + field.Name
		methodName := "Get" + field.Name
		method := reflect.ValueOf(v).MethodByName(methodName)

		var s string
		if method.IsValid() && method.Type().NumIn() == 0 && method.Type().NumOut() == 1 && method.Type().Out(0).Kind() == reflect.String {
			// Ambil dari getter
			out := method.Call(nil)
			s = out[0].String()
		} else if strings.HasPrefix(field.Name, prefixField) && val.Field(i).Kind() == reflect.String {
			// Ambil langsung dari field
			s = val.Field(i).String()
		}

		if s != "" {
			aliases = append(aliases, s)
		}
	}

	return aliases
}
