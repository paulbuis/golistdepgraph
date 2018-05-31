package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strings"
)

type JsonObject map[string]interface{}

func NewJsonObject() JsonObject {
	tmp := make(map[string]interface{})
	return tmp
}

func (jo JsonObject) String() string {
	var buf bytes.Buffer
	encoder := json.NewEncoder(&buf)
	encoder.Encode(jo)
	return buf.String()
}

func (jo JsonObject) GetString(key string) string {
	i, ok := jo[key]
	if !ok {
		return ""
	}
	switch v := i.(type) {
	case string:
		return v
	}

	return ""
}

func (jo JsonObject) GetStringSlice(key string) []string {
	i, ok := jo[key]
	if !ok {
		return []string{}
	}
	switch v := i.(type) {
	case []interface{}:
		a := make([]string, len(v))
		for index, elem := range v {
			switch e := elem.(type) {
			case string:
				a[index] = e
			default:
				a[index] = ""
			}
		}
		return a
	}

	return []string{}
}

func (jo JsonObject) GetBool(key string) bool {
	i, ok := jo[key]
	if !ok {
		return false
	}
	switch v := i.(type) {
	case bool:
		return v
	}

	return false
}

type JsonSeq []JsonObject

func NewJsonSeq(buf []byte) JsonSeq {
	seq := []JsonObject{}
	reader := bytes.NewReader(buf)
	decoder := json.NewDecoder(reader)
	var err error
	for err == nil {
		jo := NewJsonObject()
		err = decoder.Decode(&jo)
		seq = append(seq, jo)
	}
	if len(seq) == 0 {
		return seq
	} else {
		return seq[:len(seq)-1]
	}
}

func (js JsonSeq) String() string {
	var b strings.Builder
	fmt.Fprintln(&b, "[")
	for i, jo := range js {
		if i != 0 {
			fmt.Fprintln(&b, ",")
		}
		fmt.Fprint(&b, jo.String())
	}
	fmt.Fprintln(&b, "]")
	return b.String()
}
