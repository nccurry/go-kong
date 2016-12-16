package main

import (
	"encoding/json"
	"fmt"
	"github.com/fatih/structs"
	"strings"
	"bytes"
	"io"
	"log"
	"io/ioutil"
)

type Parent struct {
	FieldOne string `json:"field_one"`
	FieldTwo map[string]interface{} `json:"field_two"`
}

type Child struct {
	Parent
	FieldTwo OtherDef `json:"field_two"`
}

func (c *Child) ToParent() *Parent {

	s := structs.New(c.FieldTwo)
	fieldNames := s.Names()

	fieldTwo := make(map[string]interface{})

	for _, v := range fieldNames {
		jsonTags := strings.Split(s.Field(v).Tag("json"), ",")
		fieldTwo[jsonTags[0]] = s.Field(v).Value()
	}

	parent := &c.Parent
	parent.FieldTwo = fieldTwo
	return parent
}

type OtherDef struct {
	FieldThree string `json:"field_three,omitempty"`
	FieldFour int `json:"field_four,omitempty"`
}

func main() {



	myJson := "{\"field_one\": \"yay\", \"field_two\": {\"field_three\":\"test\", \"field_four\": 10}}"

	var child Child
	json.Unmarshal([]byte(myJson), &child)

	test := child.ToParent()
	fmt.Printf("%+v", test)
	//myj, _ := json.Marshal(test)

	var buf io.ReadWriter
	buf = new(bytes.Buffer)
	err := json.NewEncoder(buf).Encode(&test)
	if err != nil {
		log.Fatal(err)
	}

	b, _ := ioutil.ReadAll(buf)
	fmt.Println(string(b))
}
