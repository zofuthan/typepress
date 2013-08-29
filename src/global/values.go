package global

import (
	"net/url"
	"reflect"
	"strings"

	"github.com/gosexy/to"
)

// fetch value of keys from url.Values
func FetchValue(values url.Values, keys ...string) (r map[string]interface{}) {
	r = map[string]interface{}{}
	for _, k := range keys {
		r[k] = values.Get(k)
	}
	return
}

// fetch value of keys from url.Values, keys must be exists
func FetchValueExists(values url.Values, keys ...string) (r map[string]interface{}, allExists bool) {
	r = map[string]interface{}{}
	for _, k := range keys {
		if _, ok := values[k]; ok {
			r[k] = values.Get(k)
		}
	}
	allExists = len(r) == len(keys)
	return
}

// check keys of url.Values must be exists
func ValuesExists(values url.Values, keys ...string) bool {
	for _, k := range keys {
		if _, ok := values[k]; !ok {
			return false
		}
	}
	return true
}

// filter url.Values
func FilterValues(values url.Values, keys ...string) {
LOOP:
	for key, _ := range values {
		for _, k := range keys {
			if k == key {
				continue LOOP
			}
		}
		values.Del(key)
	}
	return
}

// fetch value of keys from map, if not exists set ""
func FetchMap(mp map[string]interface{}, keys ...string) (r map[string]interface{}) {
	r = map[string]interface{}{}
	for _, k := range keys {
		if v, ok := mp[k]; ok {
			r[k] = v
		} else {
			r[k] = ""
		}
	}
	return
}

// fetch value of keys from map, keys must be exists
func FetchMapExists(mp map[string]interface{}, keys ...string) (r map[string]interface{}, allExists bool) {
	for _, k := range keys {
		if v, ok := mp[k]; ok {
			r[k] = v
		}
	}
	allExists = len(r) == len(keys)
	return
}

func ToStruct(mp map[string]interface{}, p interface{}) {
	rv := reflect.Indirect(reflect.ValueOf(p))
	if !rv.IsValid() || rv.Kind() != reflect.Struct {
		LogDebug(0, nil, "ToStruct rv.Kind", rv.Kind())
		return
	}
	for k, val := range mp {
		k := strings.Title(k)
		v := rv.FieldByName(k)
		if !v.IsValid() || !v.CanSet() {
			LogDebug(0, nil, "ToStruct", k, v.Kind(), v.CanSet())
			continue
		}
		t := v.Type()
		switch t.Kind() {
		default:
			LogDebug(0, nil, "ToStruct", k, v.Kind(), v.CanSet())
		case reflect.String:
			v.SetString(to.String(val))
		case reflect.Int, reflect.Int64, reflect.Int32, reflect.Int16, reflect.Int8:
			v.SetInt(to.Int64(val))
		case reflect.Uint, reflect.Uint64, reflect.Uint32, reflect.Uint16, reflect.Uint8:
			v.SetUint(to.Uint64(val))
		case reflect.Float32, reflect.Float64:
			v.SetFloat(to.Float64(val))
		case reflect.Bool:
			v.SetBool(to.Bool(val))
		case reflect.Struct:
			switch t.PkgPath() + "." + t.Name() {
			case "time.Time":
				v.Set(reflect.ValueOf(to.Time(val)))
			case "time.Duration":
				v.Set(reflect.ValueOf(to.Duration(val)))
			}
		}
	}
}
