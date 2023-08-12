package main

import (
	"strconv"
	"time"
)

func ToString[T any](t interface{}) string {
	switch v := t.(type) {
	case string:
		return v
	case int:
		return strconv.Itoa(v)
	case float64:
		return strconv.FormatFloat(v, 'E', -1, 64)
	case bool:
		return strconv.FormatBool(v)
	case time.Time:
		return v.Format(time.RFC3339)
	case nil:
		return "NULL"
	case []byte:
		return string(v)
	case float32:
		return strconv.FormatFloat(float64(v), 'f', -1, 32)
	}

	return ""
}
