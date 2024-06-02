package handler

import (
	"fmt"
	"net/url"
	"reflect"
	"strconv"
	"strings"
)

func prop(q url.Values, key string) string {
	values := q[key]
	value := ""
	if len(values) > 0 {
		value = values[0]
		// log.Printf("%s %v", key, value)
	}
	return value
}

func intProp(q url.Values, key string, def int) (int, error) {
	s := prop(q, key)
	v := def
	if s != "" {
		value, err := strconv.ParseInt(s, 10, 64)
		if err != nil {
			return 0, err
		}
		v = int(value)
	}
	return v, nil
}

func floatProp(q url.Values, key string, def float64) (float64, error) {
	s := prop(q, key)
	v := def
	if s != "" {
		value, err := strconv.ParseFloat(s, 64)
		if err != nil {
			return 0, err
		}
		v = value
	}
	return v, nil
}

func populateStructFromQuery(query url.Values, dest interface{}) error {
	destVal := reflect.ValueOf(dest).Elem()

	// Collect all the expected property names in a set
	expectedProperties := make(map[string]struct{})
	for i := 0; i < destVal.NumField(); i++ {
		field := destVal.Type().Field(i)
		expectedProperties[strings.ToLower(field.Name)] = struct{}{}
	}

	// Check for unknown properties
	for key := range query {
		if _, found := expectedProperties[key]; !found {
			return fmt.Errorf("unknown property: %s", key)
		}
	}

	for i := 0; i < destVal.NumField(); i++ {
		field := destVal.Type().Field(i)
		queryValue := query.Get(strings.ToLower(field.Name))

		if queryValue == "" {
			continue
		}

		fieldVal := destVal.Field(i)
		switch fieldVal.Kind() {
		case reflect.Bool:
			if queryValue == "true" {
				fieldVal.SetBool(true)
			} else if queryValue == "false" {
				fieldVal.SetBool(false)
			} else {
				return fmt.Errorf("%s: expected true or false got %v", field.Name, queryValue)
			}
		case reflect.String:
			fieldVal.SetString(queryValue)
		case reflect.Int:
			if value, err := strconv.Atoi(queryValue); err == nil {
				fieldVal.SetInt(int64(value))
			} else {
				return fmt.Errorf("%s: %v", field.Name, err)
			}
		case reflect.Float64:
			if value, err := strconv.ParseFloat(queryValue, 64); err == nil {
				fieldVal.SetFloat(value)
			} else {
				return fmt.Errorf("%s: %v", field.Name, err)
			}
		}
	}
	return nil
}
