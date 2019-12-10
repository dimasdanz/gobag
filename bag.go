// Package gobag provides a simple utility to read json response without any custom struct.
// Works for nested json using dot notation
//
// A quick usage
//
// supposed body is a json like this
//  {"name": "gobag", "amount": 1}
// then
//  var data map[string]interface{}
//  _ = json.Unmarshal(body, &data)
//  bag := gobag.Bagify(data)
//  bag.GetString("name") // gobag
//  bag.GetInt("amount") // 1
package gobag

import (
	"strconv"
	"strings"
)

// Bag only contains contains filtered or unexported fields
type Bag struct {
	items interface{}
}

// Bagify constructs Bag by providing anything
func Bagify(data interface{}) Bag {
	return Bag{
		items: data,
	}
}

// IsEmpty returns whether the data constructed is empty or not
func (b Bag) IsEmpty() bool {
	return b.items == nil
}

// Get a value of any key including nested key by using dot notation
func (b Bag) Get(key string) interface{} {
	if b.IsEmpty() {
		return nil
	}

	if key == "" {
		return b.items
	}

	if !strings.Contains(key, ".") {
		return value(b.items, key)
	}

	i := b.items
	for _, k := range strings.Split(key, ".") {
		i = value(i, k)

		// shortcircuit
		if i == nil {
			return nil
		}
	}

	return i
}

// GetString returns value in string
func (b Bag) GetString(key string) string {
	val, ok := b.Get(key).(string)
	if ok {
		return val
	}

	return ""
}

// GetArrayString returns value in array of string
// defaults to "" if any of the member is not integer
func (b Bag) GetArrayString(key string) []string {
	val, ok := b.Get(key).([]interface{})
	if !ok {
		return []string{}
	}

	v := make([]string, len(val))
	for i := range val {
		f, ok := val[i].(string)
		if ok {
			v[i] = f
		} else {
			v[i] = ""
		}
	}

	return v
}

// GetInt returns value in integer. Works for json numbers
func (b Bag) GetInt(key string) int {
	// json integer is float64
	val, ok := b.Get(key).(float64)
	if ok {
		return int(val)
	}

	return 0
}

// GetArrayInt returns value in array of integer
// defaults to 0 if any of the member is not integer
func (b Bag) GetArrayInt(key string) []int {
	val, ok := b.Get(key).([]interface{})
	if !ok {
		return []int{}
	}

	v := make([]int, len(val))
	for i := range val {
		f, ok := val[i].(float64)
		if ok {
			v[i] = int(f)
		} else {
			v[i] = 0
		}
	}

	return v
}

// GetFloat returns value in integer. Works for json numbers
func (b Bag) GetFloat(key string) float64 {
	val, ok := b.Get(key).(float64)
	if ok {
		return val
	}

	return 0
}

// GetArrayFloat returns value in array of float64
// defaults to 0 if any of the member is not float64
func (b Bag) GetArrayFloat(key string) []float64 {
	val, ok := b.Get(key).([]interface{})
	if !ok {
		return []float64{}
	}

	v := make([]float64, len(val))
	for i := range val {
		f, ok := val[i].(float64)
		if ok {
			v[i] = f
		} else {
			v[i] = 0
		}
	}

	return v
}

// GetBool returns value in boolean
func (b Bag) GetBool(key string) bool {
	val, ok := b.Get(key).(bool)
	if ok {
		return val
	}

	if b.GetString(key) == "true" {
		return true
	}

	return false
}

// GetArray returns map[string]interface{}
func (b Bag) GetMapString(key string) map[string]interface{} {
	val, ok := b.Get(key).(map[string]interface{})
	if ok {
		return val
	}

	return map[string]interface{}{}
}

// value returns its underlying data by key
// returns nil when something wrong happens
// for shortcircuiting
func value(data interface{}, key string) interface{} {
	if items, ok := data.(map[string]interface{}); ok {
		return items[key]
	}

	if items, ok := data.([]interface{}); ok {
		key, err := strconv.Atoi(key)
		if err != nil {
			return nil
		}

		if len(items) > key {
			return items[key]
		}

		return nil
	}

	return nil
}
