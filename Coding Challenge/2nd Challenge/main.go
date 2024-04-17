package main

import (
	"encoding/json"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"
)

func main() {
	input := `{
			"number_1": {"N": "1.50"},
			"string_1": {"S": "784498 "},
			"string_2": {"S": "2014-07-16T20:55:46Z"},
			"map_1": {
				"M": {
					"bool_1": {"BOOL": "truthy"},
					"null_1": {"NULL ": "true"},
					"list_1": {"L": [
						{"S": ""},
						{"N": "011"},
						{"N": "5215s"},
						{"BOOL": "f"},
						{"NULL": "0"}
					]}
				}
			},
			"list_2": {"L": "noop"},
			"list_3": {"L": ["noop"]},
			"": {"S": "noop"}
		}`

	var rawData map[string]json.RawMessage
	if err := json.Unmarshal([]byte(input), &rawData); err != nil {
		log.Fatalf("Error parsing JSON: %v", err)
	}

	transformed := transform(rawData)
	outputJSON, err := json.MarshalIndent([]interface{}{transformed}, "", "  ")
	if err != nil {
		log.Fatalf("Error marshalling JSON: %v", err)
	}
	fmt.Println(string(outputJSON))
}

// this one transforms map of josn to itemType and value
func transform(rawData map[string]json.RawMessage) map[string]interface{} {
	result := make(map[string]interface{})
	for key, value := range rawData {
		if key == "" || strings.TrimSpace(key) == "" {
			continue
		}

		var item map[string]json.RawMessage
		json.Unmarshal(value, &item)
		for itemType, itemValue := range item {
			key = strings.TrimSpace(key)
			switch itemType {
			case "N":
				result[key] = transformNumber(itemValue, true)
			case "S":
				result[key] = transformString(itemValue, true)
			case "BOOL":
				result[key] = transformBool(itemValue, true)
			case "NULL ":
				if transformNull(itemValue, true) {
					result[key] = nil
				}
			case "L":
				list, _ := transformList(itemValue)
				if len(list) > 0 {
					result[key] = list
				}
			case "M":
				mapResult := transformMap(itemValue)
				if len(mapResult) > 0 {
					result[key] = mapResult
				}
			}
		}
	}
	return result
}

func convertToString(data json.RawMessage) string {
	var str string
	if err := json.Unmarshal(data, &str); err != nil {
		return ""
	}
	return str
}

// transform to number
func transformNumber(data json.RawMessage, convert bool) interface{} {
	str := convertToString(data)
	str = strings.TrimSpace(str)
	str = strings.TrimLeft(str, "0")
	// Now we convert the string to a float64
	if num, err := strconv.ParseFloat(str, 64); err == nil {
		return num
	}
	return 0
}

// transform to string
func transformString(data json.RawMessage, convert bool) interface{} {
	str := convertToString(data)
	str = strings.TrimSpace(str)
	if t, err := time.Parse(time.RFC3339, str); err == nil {
		return t.Unix()
	}
	return str
}

// transform to bool value
func transformBool(data json.RawMessage, convert bool) interface{} {
	str := convertToString(data)
	str = strings.TrimSpace(strings.ToLower(string(data)))
	switch str {
	case "1", "t", "true":
		return true
	case "0", "f", "false":
		return false
	}
	return nil
}

// // transform to bool null
func transformNull(data json.RawMessage, convert bool) bool {
	str := convertToString(data)
	str = strings.TrimSpace(strings.ToLower(string(data)))
	return str == "true"
}

// transform to list
func transformList(data json.RawMessage) ([]interface{}, error) {
	var items []json.RawMessage
	if err := json.Unmarshal(data, &items); err != nil {
		return nil, err
	}
	var result []interface{}
	for _, item := range items {
		var elem map[string]json.RawMessage
		json.Unmarshal(item, &elem)
		for itemType, itemValue := range elem {
			switch itemType {
			case "N":
				if num := transformNumber(itemValue, false); num != nil {
					result = append(result, num)
				}
			case "BOOL":
				if boolVal := transformBool(itemValue, false); boolVal != nil {
					result = append(result, boolVal)
				}
			case "NULL":
				if transformNull(itemValue, false) {
					result = append(result, nil)
				}
			}
		}
	}
	return result, nil
}

// transform to map
func transformMap(data json.RawMessage) map[string]interface{} {
	var rawMap map[string]json.RawMessage
	json.Unmarshal(data, &rawMap)
	return transform(rawMap)
}
