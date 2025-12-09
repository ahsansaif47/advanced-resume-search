package utils

import (
	"encoding/json"
	"fmt"
	"log"
	"regexp"
	"strings"
)

func InterfaceSliceToStringSlice(input any) []string {
	if input == nil {
		return nil
	}

	ifaceSlice, ok := input.([]interface{})
	if !ok {
		return nil
	}

	result := make([]string, len(ifaceSlice))
	for i, v := range ifaceSlice {
		if ok {
			switch v := v.(type) {
			case string:
				result[i] = v
			case map[string]any:
				result[i] = MarshalToString(v)
			default:
				fmt.Printf("unsupported type: %T", v)
			}
		}
	}
	return result
}

func MarshalToString(data map[string]any) string {
	data_str, err := json.Marshal(data)
	if err != nil {
		log.Fatalln("Could not convert to string")
	}
	return string(data_str)
}

func SanitizeFileName(name string) string {
	name = strings.ToLower(name)
	re := regexp.MustCompile(`[^a-z0-9-]+`)
	name = re.ReplaceAllString(name, "-")
	name = strings.Trim(name, "-")
	return name
}
