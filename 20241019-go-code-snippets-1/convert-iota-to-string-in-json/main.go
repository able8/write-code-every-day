package main

import (
	"bytes"
	"encoding/json"
	"fmt"
)

// https://www.golangcode.com/convert-iota-to-string-in-json/

type Gender int

const (
	GenderNotSet = iota
	GenderMale
	GenderFemale
	GenderOther
)

var toString = map[Gender]string{
	GenderNotSet: "Not Set",
	GenderMale:   "Male",
	GenderFemale: "Female",
	GenderOther:  "Other",
}

var toID = map[string]Gender{
	"Not Set": GenderNotSet,
	"Male":    GenderMale,
	"Female":  GenderFemale,
	"Other":   GenderOther,
}

func (g Gender) MarshalJSON() ([]byte, error) {
	buffer := bytes.NewBufferString(`"`)
	buffer.WriteString(toString[g])
	buffer.WriteString(`"`)
	return buffer.Bytes(), nil
}

func (g *Gender) UnmarshalJSON(b []byte) error {
	var j string
	if err := json.Unmarshal(b, &j); err != nil {
		return err
	}
	*g = toID[j]
	return nil
}

type Human struct {
	Gender Gender `json:"gender"`
}

func main() {
	me := Human{
		Gender: GenderMale,
	}
	prettyJSON, _ := json.MarshalIndent(me, "", "    ")
	fmt.Println(string(prettyJSON))
}
