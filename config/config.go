package config

import (
	"reflect"
	"strings"
)

func getFactory(key string) reflect.Type {
	factories := map[string]interface{}{
		"App":      App{},
		"Database": Database{},
	}

	return reflect.TypeOf(factories[key])
}

func Get(key string, defaultval ...string) interface{} {
	split := strings.Split(key, ".")
	rtype := getFactory(split[0])

	dfval := ""
	if len(defaultval) > 0 {
		dfval = defaultval[0]
	}

	field, fieldfound := rtype.FieldByName(split[1])
	if !fieldfound {
		return dfval
	}

	value, valuefound := field.Tag.Lookup("default")
	if !valuefound {
		return dfval
	}

	return value
}
