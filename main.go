package main

import (
	"bytes"
	"encoding/json"
	"fmt"
)

func main() {
	// String containing arbitrary bytes
	str := "\xad"
	fmt.Printf("raw %q: %x\n", str, str)
	got := bytes.Buffer{}
	encoder := json.NewEncoder(&got)
	encoder.SetEscapeHTML(false)
	err := encoder.Encode(&str)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	fmt.Printf("encode '%q': %x\n", got.Bytes(), got.Bytes())
	fmt.Printf("encode '%s'\n", got.Bytes())

	mgot, err := json.Marshal(str)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}
	fmt.Printf("marshal '%q': %x\n", mgot, mgot)
	fmt.Printf("marshal '%s'\n", mgot)

	var theMStr string
	err = json.Unmarshal(mgot, &theMStr)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}
	fmt.Printf("unmarshal theMStr %q: %x\n", theMStr, theMStr)

	var theStr string
	err = json.Unmarshal(got.Bytes(), &theStr)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}
	fmt.Printf("unmarshal %q: %x\n", theStr, theStr)
}
