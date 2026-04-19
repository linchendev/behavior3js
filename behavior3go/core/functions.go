package core

import (
	"crypto/rand"
	"fmt"
	"time"
)

var Now = func() int64 {
	return time.Now().UnixMilli()
}

func createUUID() string {
	b := make([]byte, 16)
	if _, err := rand.Read(b); err != nil {
		return fmt.Sprintf("%d", time.Now().UnixNano())
	}

	b[6] = (b[6] & 0x0f) | 0x40
	b[8] = (b[8] & 0x3f) | 0x80

	return fmt.Sprintf(
		"%08x-%04x-%04x-%04x-%012x",
		b[0:4],
		b[4:6],
		b[6:8],
		b[8:10],
		b[10:16],
	)
}

func copyMap(source map[string]any) map[string]any {
	if source == nil {
		return nil
	}
	if len(source) == 0 {
		return map[string]any{}
	}

	target := make(map[string]any, len(source))
	for key, value := range source {
		if nested, ok := value.(map[string]any); ok {
			target[key] = copyMap(nested)
			continue
		}
		target[key] = value
	}
	return target
}

func clonePropertiesForLoad(source map[string]any) map[string]any {
	if len(source) == 0 {
		return nil
	}
	return copyMap(source)
}

func toInt(value any) (int, bool) {
	switch typed := value.(type) {
	case int:
		return typed, true
	case int8:
		return int(typed), true
	case int16:
		return int(typed), true
	case int32:
		return int(typed), true
	case int64:
		return int(typed), true
	case uint:
		return int(typed), true
	case uint8:
		return int(typed), true
	case uint16:
		return int(typed), true
	case uint32:
		return int(typed), true
	case uint64:
		return int(typed), true
	case float32:
		return int(typed), true
	case float64:
		return int(typed), true
	default:
		return 0, false
	}
}

func ToInt64(value any) (int64, bool) {
	switch typed := value.(type) {
	case int:
		return int64(typed), true
	case int8:
		return int64(typed), true
	case int16:
		return int64(typed), true
	case int32:
		return int64(typed), true
	case int64:
		return typed, true
	case uint:
		return int64(typed), true
	case uint8:
		return int64(typed), true
	case uint16:
		return int64(typed), true
	case uint32:
		return int64(typed), true
	case uint64:
		return int64(typed), true
	case float32:
		return int64(typed), true
	case float64:
		return int64(typed), true
	default:
		return 0, false
	}
}

func GetIntProperty(properties map[string]any, key string, defaultValue int) (int, bool) {
	if properties == nil {
		return defaultValue, false
	}

	value, ok := properties[key]
	if !ok {
		return defaultValue, false
	}

	integer, parsed := toInt(value)
	return integer, parsed
}

func boolValue(value any) bool {
	typed, ok := value.(bool)
	return ok && typed
}
