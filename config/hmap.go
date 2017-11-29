package config

import "github.com/qwantix/qxeye/util"

type hmap map[string]interface{}

func (h *hmap) Has(key string) bool {
	return (*h)[key] != nil
}

func (h *hmap) String(key string) string {
	if h.Has(key) {
		switch (*h)[key].(type) {
		case string:
			return (*h)[key].(string)
		default:
			util.Warn("Invalid value type for ", key)
		}
	}
	return ""
}

func (h *hmap) Bool(key string) bool {
	if h.Has(key) {
		switch (*h)[key].(type) {
		case bool:
			return (*h)[key].(bool)
		default:
			util.Warn("Invalid value type for ", key)
		}
	}
	return false
}

func (h *hmap) Int(key string) int {
	if h.Has(key) {
		switch (*h)[key].(type) {
		case int:
			return (*h)[key].(int)
		case float64:
			return int((*h)[key].(float64))
		default:
			util.Warn("Invalid value type for ", key)
		}
	}
	return 0
}

func (h *hmap) Float(key string) float64 {
	if h.Has(key) {
		switch (*h)[key].(type) {
		case float64:
			return (*h)[key].(float64)
		default:
			util.Warn("Invalid value type for ", key)
		}
	}
	return 0.0
}
