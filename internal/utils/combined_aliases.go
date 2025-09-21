package utils

import (
	"reflect"
	"strings"
)

func CombineAliases(v interface{}, prefixField string) []string {
	val := reflect.ValueOf(v)
	if !val.IsValid() {
		return nil
	}

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

		if field.PkgPath != "" {
			continue
		}

		methodName := "Get" + field.Name
		method := reflect.ValueOf(v).MethodByName(methodName)

		var s string
		if method.IsValid() &&
			method.Type().NumIn() == 0 &&
			method.Type().NumOut() == 1 &&
			method.Type().Out(0).Kind() == reflect.String {

			out := method.Call(nil)
			s = out[0].String()
		} else if strings.HasPrefix(field.Name, prefixField) {
			f := val.Field(i)
			switch f.Kind() {
			case reflect.String:
				s = f.String()
			case reflect.Ptr:
				if !f.IsNil() && f.Elem().Kind() == reflect.String {
					s = f.Elem().String()
				}
			}
		}

		if strings.TrimSpace(s) != "" {
			aliases = append(aliases, s)
		}
	}

	return aliases
}
