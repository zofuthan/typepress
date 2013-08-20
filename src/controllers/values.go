package controllers

import (
	"net/url"
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
