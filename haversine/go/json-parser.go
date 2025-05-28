package main

import (
	"errors"
	"fmt"
	"log"
	"os"
	"strconv"
	"unicode"
)

type Parser struct {
	data []byte
	pos  int
}

func (parser *Parser) nextByte() {
	if parser.pos < len(parser.data)-1 {
		parser.pos++
	}
}

func (parser *Parser) getByte() byte {
	return parser.data[parser.pos]
}

func parseJson() {
	file, err := os.ReadFile("./output.json")

	if err != nil {
		log.Fatal(err)
	}
	parser := Parser{
		data: file,
		pos:  0,
	}
	var output interface{}
	if parser.getByte() == '[' {
		fmt.Println("arr")
		output, err = parser.parseArray()
	} else if parser.getByte() == '{' {
		output, err = parser.parseObject()
	} else {
		panic("not supported char")
	}

	fmt.Println("output", output)
}

func (parser *Parser) parseArray() ([]interface{}, error) {
	if parser.getByte() != '[' {
		panic("expected to start with {")
	}
	parser.nextByte()
	parser.skipWhitespace()

	arr := []interface{}{}

	for {
		parser.skipWhitespace()
		if parser.getByte() == ']' {
			parser.nextByte()
			break
		}

		val, err := parser.parseValue()
		if err != nil {
			panic(err)
		}
		fmt.Println("val")
		fmt.Println(val)
		arr = append(arr, val)

		parser.skipWhitespace()

		if parser.getByte() == ',' {
			parser.nextByte()
			parser.skipWhitespace()
		} else if parser.getByte() != ']' {
			fmt.Println(string(parser.getByte()))
			panic("expected data  to be ]")
		}
	}

	return arr, nil
}

func (parser *Parser) parseObject() (map[string]interface{}, error) {
	obj := make(map[string]interface{})
	if parser.getByte() != '{' {
		panic("expected to start with {")
	}
	parser.nextByte()
	for {
		parser.skipWhitespace()
		if parser.getByte() == '}' {
			fmt.Println("end?")
			parser.nextByte()
			break
		}

		key, err := parser.parseString()
		if err != nil {
			panic(err)
		}
		parser.skipWhitespace()
		if parser.getByte() != ':' {
			panic("invalid struct, after key expected to be :")
		}
		parser.nextByte()

		parser.skipWhitespace()
		val, err := parser.parseValue()
		fmt.Println("parsed value number", val)
		if err != nil {
			panic(err)
		}
		obj[key] = val

		parser.skipWhitespace()
		if parser.getByte() == ',' {
			parser.nextByte()
			parser.skipWhitespace()
		} else if parser.getByte() != '}' {
			fmt.Println(string(parser.getByte()))
			panic("expected data  to be }")
		}
	}
	return obj, nil
}

func (parser *Parser) parseNumber() (interface{}, error) {
	start := parser.pos
	if parser.getByte() == '-' {
		parser.nextByte()
	}
	digitsFound := false

	for (len(parser.data) > parser.pos) && unicode.IsDigit(rune(parser.getByte())) {
		parser.nextByte()
		digitsFound = true
	}

	if !digitsFound {
		return nil, errors.New("invalid number")
	}

	if parser.getByte() == '.' {
		parser.nextByte()
		if !unicode.IsDigit(rune(parser.getByte())) {
			return nil, errors.New("invalid float format")
		}
		for (len(parser.data) > parser.pos) && unicode.IsDigit(rune(parser.getByte())) {
			parser.nextByte()
		}
	}

	numStr := parser.data[start:parser.pos]
	num, err := strconv.ParseFloat(string(numStr), 64)
	if err != nil {
		return nil, fmt.Errorf("invalid number: %v", err)
	}
	return num, nil
}

func (parser *Parser) parseValue() (interface{}, error) {
	switch ch := parser.getByte(); {
	case ch == '{':
		return parser.parseObject()
	case ch == '[':
		return parser.parseArray()
	case ch == '"':
		return parser.parseString()
	case ch == '-' || (ch >= '0' && ch <= '9'):
		return parser.parseNumber()
	default:
		fmt.Println(string(parser.getByte()))
		panic("not supported")
	}
}

func (parser *Parser) skipWhitespace() {
	for len(parser.data) > parser.pos && unicode.IsSpace(rune(parser.data[parser.pos])) {
		parser.pos++
	}
}

func (parser *Parser) parseString() (string, error) {
	if parser.getByte() != '"' {
		fmt.Println(string(parser.getByte()))
		return "", fmt.Errorf("Expected key to start with \"")
	}
	parser.nextByte()
	key := ""
	for {
		if parser.getByte() == '"' {
			parser.nextByte()
			break
		}
		key += string(parser.getByte())
		parser.nextByte()
	}

	return key, nil
}
