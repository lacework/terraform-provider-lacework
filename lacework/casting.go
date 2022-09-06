package lacework

import (
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

// extract an attribute from the provided ResourceData and convert it into a string slice
func castAttributeToStringSlice(d *schema.ResourceData, attr string) []string {
	return castAndTransformStringSlice(d.Get(attr).([]interface{}), func(s string) string { return s })
}

func castAndTitleCaseStringSlice(d *schema.ResourceData, attr string) []string {
	return castAndTransformStringSlice(d.Get(attr).([]interface{}), func(s string) string { return cases.Title(language.English).String(strings.ToLower(s)) })
}

func castAndUpperStringSlice(iArray []interface{}) []string {
	return castAndTransformStringSlice(iArray, func(s string) string { return strings.ToUpper(s) })
}

// turn an interface slice into a string slice
func castStringSlice(iArray []interface{}) []string {
	return castAndTransformStringSlice(iArray, func(s string) string { return s })
}

// turn an interface slice into a string slice and apply a transformation func
func castAndTransformStringSlice(iArray []interface{}, f func(string) string) []string {
	a := make([]string, 0)
	for _, v := range iArray {
		if v == nil {
			continue
		}
		a = append(a, f(v.(string)))
	}
	return a
}

// turn a string slice into an instance slice
func castStringSliceToInterface(strs []string) []interface{} {
	arr := make([]interface{}, len(strs))
	for i, str := range strs {
		arr[i] = str
	}
	return arr
}

// extract an attribute from the provided ResourceData and convert it into a map of strings
// with string keys
func castAttributeToStringKeyMapOfStrings(d *schema.ResourceData, attr string) map[string]string {
	mapStrings := d.Get(attr).(map[string]interface{})
	newParams := make(map[string]string, len(mapStrings))
	for key, val := range mapStrings {
		newParams[key] = val.(string)
	}

	return newParams
}

// extract an attribute from the provided ResourceData and convert it into an array of map of strings
// with string keys. (needed for API v2 ContainerRegistry Limits)
//
// Example of key/value TypeSet:
//
// label {
//   key = "foo"
//   value = "bar"
// }
//
// label {
//   key = "abc"
//   value = "xyz"
// }
//
// label {
//   key = "key"
//   value = "value"
// }
//
// The returned array of map of strings with string keys:
//
// []map[string]string{
//   {"foo": "bar"},
//   {"abc": "xyz"},
//   {"key": "value"},
// }
func castAttributeToArrayOfKeyValueMap(d *schema.ResourceData, attr string) []map[string]string {
	list := d.Get(attr).(*schema.Set).List()
	aMap := make([]map[string]string, len(list))
	for i, v := range list {
		val := v.(map[string]interface{})
		aMap[i] = map[string]string{val["key"].(string): val["value"].(string)}
	}

	return aMap
}

// extract an attribute from the provided ResourceData and convert it into a map of strings
// with string keys
func castAttributeToArrayKeyMapOfStrings(d *schema.ResourceData, attr string) []map[string]string {
	if castMap, ok := d.Get(attr).(map[string]interface{}); ok {
		mapList := make([]map[string]string, 1)
		stringMap := make(map[string]string)
		for key, val := range castMap {
			stringMap[key] = val.(string)
		}
		mapList[0] = stringMap
		return mapList
	}

	return []map[string]string{}
}

func castAttributeToArrayOfCustomKeyValueMap(d *schema.ResourceData, attr string, key string, value string) []map[string]string {
	list := d.Get(attr).(*schema.Set).List()
	aMap := make([]map[string]string, len(list))
	for i, v := range list {
		val := v.(map[string]interface{})
		aMap[i] = map[string]string{val[key].(string): val[value].(string)}
	}

	return aMap
}

// convert an array of map of strings with string keys to a key/value TypeSet
// needed for API v2 ContainerRegistry Limits
//
// @afiune This function reverts the array of map (from APIv2) to a TypeSet array
//         which is pretty much what castAttributeToArrayOfKeyValueMap() returns
//
// Example of an array of map of strings with string keys:
//
// []map[string]string{
//   {"foo": "bar"},
//   {"abc": "xyz"},
//   {"key": "value"},
// }
//
// The returned array of key/value TypeSet:
//
// mockLabels = []map[string]string{
//   {"key": "foo", "value": "bar"},
//   {"key": "abc", "value": "xyz"},
//   {"key": "key", "value": "value"},
// }
func castArrayOfStringKeyMapOfStringsToLimitByLabelSet(list []map[string]string) []map[string]string {
	aMap := make([]map[string]string, len(list))

	for i, mStrings := range list {
		aMap[i] = map[string]string{}
		for k, v := range mStrings {
			aMap[i]["key"] = k
			aMap[i]["value"] = v
		}
	}

	return aMap
}

// TODO @afiune remove this function when we release v1.0
func joinMapStrings(m map[string]string, delimit string) string {
	out := make([]string, 0)
	for _, val := range m {
		out = append(out, val)
	}
	return strings.Join(out, delimit)
}

func ContainsStr(array []string, expected string) bool {
	for _, value := range array {
		if expected == value {
			return true
		}
	}
	return false
}
