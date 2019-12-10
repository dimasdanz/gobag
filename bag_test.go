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

	s.Assert().True(bag.IsEmpty())
	s.Assert().Nil(bag.Get("anykeyshouldreturnnilonemptybag"))
}

func (s *BagTestSuite) TestCorrect() {
	body, _ := ioutil.ReadAll(ioutil.NopCloser(bytes.NewReader([]byte(payload))))
	var data map[string]interface{}
	_ = json.Unmarshal(body, &data)

	bag := gobag.Bagify(data)

	s.Assert().NotNil(bag.Get(""))
	s.Assert().Equal("ok", bag.GetString("string"))
	s.Assert().Equal(1, bag.GetInt("int"))
	s.Assert().Equal(0.5, bag.GetFloat("float"))
	s.Assert().Equal(true, bag.GetBool("bool"))
	s.Assert().Equal(true, bag.GetBool("boolstring"))
	s.Assert().Equal("foo", bag.GetString("array.0"))
	s.Assert().Equal("bar", bag.GetString("array.1"))
	s.Assert().Equal("bar", bag.GetString("object.foo"))
	s.Assert().Equal("bar", bag.GetString("array_of_object.0.foo"))
	s.Assert().Equal("bar", bag.GetString("array_of_object.1.baz"))
	s.Assert().Equal([]string{"foo", "bar"}, bag.GetArrayString("array"))
	s.Assert().Equal([]int{1, 2, 3}, bag.GetArrayInt("array_of_int"))
	s.Assert().Equal([]float64{1, 2, 3}, bag.GetArrayFloat("array_of_int"))
	s.Assert().Equal(map[string]interface{}{"baz": float64(1), "foo": "bar"}, bag.GetMapString("map_string"))
}

func (s *BagTestSuite) TestIncorrect() {
	body, _ := ioutil.ReadAll(ioutil.NopCloser(bytes.NewReader([]byte(payload))))
	var data map[string]interface{}
	_ = json.Unmarshal(body, &data)

	bag := gobag.Bagify(data)

	s.Assert().Equal("", bag.GetString("int"))
	s.Assert().Equal(0, bag.GetInt("string"))
	s.Assert().Equal(float64(0), bag.GetFloat("string"))
	s.Assert().Equal(false, bag.GetBool("string"))
	s.Assert().Equal("", bag.GetString("array.2"))
	s.Assert().Equal("", bag.GetString("array.first"))
	s.Assert().Equal("", bag.GetString("string.0"))
	s.Assert().Equal([]string{}, bag.GetArrayString("string"))
	s.Assert().Equal([]string{"", "", ""}, bag.GetArrayString("array_of_int"))
	s.Assert().Equal([]int{}, bag.GetArrayInt("string"))
	s.Assert().Equal([]int{0, 0}, bag.GetArrayInt("array"))
	s.Assert().Equal([]float64{}, bag.GetArrayFloat("string"))
	s.Assert().Equal([]float64{0, 0}, bag.GetArrayFloat("array"))
	s.Assert().Equal(map[string]interface{}{}, bag.GetMapString("array"))
}

func TestResponseBag(t *testing.T) {
	suite.Run(t, new(BagTestSuite))
}
