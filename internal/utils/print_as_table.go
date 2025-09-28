package utils

import (
	"fmt"
	"os"
	"reflect"
	"text/tabwriter"
	"time"
)

// PrintAsTable mencetak slice of struct dalam format tabel
func PrintAsTable(results interface{}) {
	val := reflect.ValueOf(results)

	if val.Kind() != reflect.Slice {
		fmt.Println("PrintAsTable hanya untuk slice of struct")
		return
	}
	if val.Len() == 0 {
		fmt.Println("Data kosong")
		return
	}

	elemType := val.Index(0).Type()

	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)

	// header
	for i := 0; i < elemType.NumField(); i++ {
		fmt.Fprintf(w, "%s\t", elemType.Field(i).Name)
	}
	fmt.Fprintln(w)

	// separator
	for i := 0; i < elemType.NumField(); i++ {
		fmt.Fprintf(w, "--------\t")
	}
	fmt.Fprintln(w)

	// isi
	for i := 0; i < val.Len(); i++ {
		row := val.Index(i)
		for j := 0; j < row.NumField(); j++ {
			field := row.Field(j)
			fieldValue := formatField(field)
			fmt.Fprintf(w, "%v\t", fieldValue)
		}
		fmt.Fprintln(w)
	}

	w.Flush()
}

func PrintStruct(s interface{}) {
	val := reflect.ValueOf(s)
	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}
	if val.Kind() != reflect.Struct {
		fmt.Println("PrintStruct hanya untuk struct")
		return
	}

	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)

	for i := 0; i < val.NumField(); i++ {
		fmt.Fprintf(w, "%s\t", val.Type().Field(i).Name)
	}
	fmt.Fprintln(w)

	for i := 0; i < val.NumField(); i++ {
		fmt.Fprintf(w, "--------\t")
	}
	fmt.Fprintln(w)

	for i := 0; i < val.NumField(); i++ {
		field := val.Field(i)
		fmt.Fprintf(w, "%v\t", formatField(field))
	}
	fmt.Fprintln(w)
	w.Flush()
}

func PrintStringSliceAsTable(slice []string) {
	for i, v := range slice {
		fmt.Printf("%02d\t%s\n", i, v)
	}
}

func formatField(v reflect.Value) interface{} {
	if v.Kind() == reflect.Ptr {
		if v.IsNil() {
			return ""
		}
		// Ambil nilai yang ditunjuk pointer
		val := v.Elem().Interface()

		// Jika time.Time → format lebih rapi
		if t, ok := val.(time.Time); ok {
			return t.Format("02/01/2006 15:04:05")
		}
		return val
	}

	// Kalau bukan pointer
	if v.Type().String() == "time.Time" {
		t := v.Interface().(time.Time)
		return t.Format("02/01/2006 15:04:05")
	}

	return v.Interface()
}
