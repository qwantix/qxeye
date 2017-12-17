package config

import "github.com/qwantix/qxeye/util"

type Hmap map[string]interface{}

func (h *Hmap) Has(key string) bool {
	return (*h)[key] != nil
}

func (h *Hmap) String(key string) string {
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

func (h *Hmap) Bool(key string) bool {
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

func (h *Hmap) Int(key string) int {
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

func (h *Hmap) Float(key string) float64 {
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
