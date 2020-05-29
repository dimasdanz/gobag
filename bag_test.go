package gobag_test

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"testing"

	"github.com/dimasdanz/gobag"
	"github.com/stretchr/testify/suite"
)

type BagTestSuite struct {
	suite.Suite
}

const payload = `
{
	"string": "ok",
	"int": 1,
	"float": 0.5,
	"bool": true,
	"boolstring": "true",
	"array": [
		"foo",
		"bar"
	],
	"object": {
		"foo": "bar"
	},
	"array_of_object": [
		{
			"foo": "bar"
		},
		{
			"baz": "bar"
		}
	],
	"array_of_int": [
		1,
		2,
		3
	],
	"map_string": {
		"foo": "bar",
		"baz": 1
	}
}`

func (s *BagTestSuite) TestEmpty() {
	bag := gobag.Bagify(nil)

	s.True(bag.IsEmpty())
	s.Nil(bag.Get("anykeyshouldreturnnilonemptybag"))
}

func (s *BagTestSuite) TestCorrect() {
	body, _ := ioutil.ReadAll(ioutil.NopCloser(bytes.NewReader([]byte(payload))))
	var data map[string]interface{}
	_ = json.Unmarshal(body, &data)

	bag := gobag.Bagify(data)

	s.NotNil(bag.Get(""))
	s.Equal("ok", bag.GetString("string"))
	s.Equal(1, bag.GetInt("int"))
	s.Equal(0.5, bag.GetFloat("float"))
	s.Equal(true, bag.GetBool("bool"))
	s.Equal(true, bag.GetBool("boolstring"))
	s.Equal("foo", bag.GetString("array.0"))
	s.Equal("bar", bag.GetString("array.1"))
	s.Equal("bar", bag.GetString("object.foo"))
	s.Equal("bar", bag.GetString("array_of_object.0.foo"))
	s.Equal("bar", bag.GetString("array_of_object.1.baz"))
	s.Equal([]string{"foo", "bar"}, bag.GetArrayString("array"))
	s.Equal([]int{1, 2, 3}, bag.GetArrayInt("array_of_int"))
	s.Equal([]float64{1, 2, 3}, bag.GetArrayFloat("array_of_int"))
	s.Equal(map[string]interface{}{"baz": float64(1), "foo": "bar"}, bag.GetMapString("map_string"))
	s.Equal(
		[]interface{}{map[string]interface{}{"foo": "bar"}, map[string]interface{}{"baz": "bar"}},
		bag.GetArray("array_of_object"),
	)
}

func (s *BagTestSuite) TestIncorrect() {
	body, _ := ioutil.ReadAll(ioutil.NopCloser(bytes.NewReader([]byte(payload))))
	var data map[string]interface{}
	_ = json.Unmarshal(body, &data)

	bag := gobag.Bagify(data)

	s.Equal("", bag.GetString("int"))
	s.Equal(0, bag.GetInt("string"))
	s.Equal(float64(0), bag.GetFloat("string"))
	s.Equal(false, bag.GetBool("string"))
	s.Equal("", bag.GetString("array.2"))
	s.Equal("", bag.GetString("array.first"))
	s.Equal("", bag.GetString("string.0"))
	s.Equal([]string{}, bag.GetArrayString("string"))
	s.Equal([]string{"", "", ""}, bag.GetArrayString("array_of_int"))
	s.Equal([]int{}, bag.GetArrayInt("string"))
	s.Equal([]int{0, 0}, bag.GetArrayInt("array"))
	s.Equal([]float64{}, bag.GetArrayFloat("string"))
	s.Equal([]float64{0, 0}, bag.GetArrayFloat("array"))
	s.Equal(map[string]interface{}{}, bag.GetMapString("array"))
	s.Equal([]interface{}{}, bag.GetArray("map_string"))
}

func (s *BagTestSuite) TestNotJSON() {
	body, _ := ioutil.ReadAll(ioutil.NopCloser(bytes.NewReader([]byte(payload))))
	var data map[string]interface{}
	_ = json.Unmarshal(body, &data)

	bag := gobag.Bagify(map[string]interface{}{
		"string":       "ok",
		"int":          1,
		"float":        0.5,
		"bool":         true,
		"array_of_int": []interface{}{1, 2, 3},
	})

	s.Equal("ok", bag.GetString("string"))
	s.Equal(1, bag.GetInt("int"))
	s.Equal(float64(0.5), bag.GetFloat("float"))
	s.Equal(true, bag.GetBool("bool"))
	s.Equal([]int{1, 2, 3}, bag.GetArrayInt("array_of_int"))
}

func TestResponseBag(t *testing.T) {
	suite.Run(t, new(BagTestSuite))
}
