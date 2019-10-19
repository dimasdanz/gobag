# gobag

Built with ❤️

[![Build Status](https://github.com/dimasdanz/gobag/workflows/Publish/badge.svg)](https://github.com/dimasdanz/gobag/actions)
[![Code Coverage](https://codecov.io/gh/dimasdanz/gobag/branch/master/graph/badge.svg)](https://codecov.io/gh/dimasdanz/gobag)
[![Go Report Card](https://goreportcard.com/badge/github.com/dimasdanz/gobag)](https://goreportcard.com/report/github.com/dimasdanz/gobag)
[![GoDoc](https://godoc.org/github.com/dimasdanz/gobag?status.svg)](https://godoc.org/github.com/dimasdanz/gobag)

A simple utility to read json response without any custom struct. Works for nested json using dot notation

## Description

Tired of making custom structs for every API call you made? Especially when the response is nested?  Well, this package simplify everything.
For whatever structure your json response is you can simplify struct construction
```go
type Person struct {
    Name string
    Birthdate string
    Height int
}

// somehow it nested somewhere deep in the json response
person := Person{
    Name: bag.GetString("profiles.0.name"),
    Birthdate: bag.GetString("profiles.0.birthdate"),
    Height: bag.GetInt("profiles.0.height")
}
```

## Usage

Supposed you make an HTTP call and received a response like this
```json
{
    "name": "Against The Current",
    "members": [
        "Chrissy Costanza",
        "Daniel Gow",
        "William Ferri"
    ],
    "formed_in": 2011,
    "profiles": [
        {
            "name": "Chrissy Costanza ❤️",
            "birthdate": "23 August 1995",
            "height": 155
        }
    ],
    "informations": {
        "labels": "Fueled by Ramen",
        "website": "atcofficial.com"
    },
    "years_of_active": [
        2011,
        2012,
        2013,
        2014,
        2015,
        2016,
        2017,
        2018,
        2019,
    ],
}
```
And in your code might look like this
```go
resp, _ := http.Get("https://example.com/against-the-current")
body, _ := ioutil.ReadAll(resp.Body)

var data map[string]interface{}
_ = json.Unmarshal(body, &data)

bag := gobag.Bagify(data)
```
That will construct the Bag struct, now we can access everything
```go
// simple get
name := bag.GetString("name") // Against The Current
formed_in := bag.GetInt("formed_in") // 2011

// accessing inside an array directly
firstElement := bag.GetString("members.0") // Chrissy Costanza

// array of string
var members []string
members = bag.GetArrayString("members") // ["Chrissy Costanza", "Daniel Gow", "William Ferri"]

// array of integer
var yearsOfActive []int
yearsOfActive = bag.GetArrayInt("years_of_active") // [2011, 2012, 2013, 2014, ...]

// wanna get a string inside a nested json?
website := bag.GetString("informations.website")

// or a deep nested object
height := bag.GetInt("profiles.0.height")
```
**Error Handling**  
In the case that the json isn't in the type expected, it will returns its zero-value
```go
// the key is missing
foo := bag.GetString("missing_key") // ""
bar := bag.GetInt("missing_key") // 0
baz := bag.GetArrayInt("missing_key") // []int{}

// or it's an integer not a string
qux := bag.GetString("formed_in") // ""
```

## Why?

This package is particularly useful when you don't want to create a custom struct for every API call response.  
Instead of creating a complex embedded struct of a response when you only need a few of the data
here's an example, you wanna create an array of people struct
```go
type Person struct {
    Name string
    Birthdate string
    Height int
}

func main() {
    bag := gobag.Bagify(jsonResponse)
    data, _ := bag.Get("data").([]interface{})

	var persons []Person
	for _, v := range data {
        // actually, you can bagify almost everything, including a json request on your handler
        b := gobag.Bagify(v)

		persons = append(persons, Person{
			Name:      b.GetString("name"),
			Birthdate: b.GetString("birthdate"),
			Height:    b.GetInt("height"),
		})
    }
    
    fmt.Println(persons)
}
```
