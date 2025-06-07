package main

import (
	"errors"
	"fmt"
	"log"
	"os"
	"reflect"
	"strconv"
	"time"
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

func parseJson() (interface{}, int, float64) {
	start := time.Now()

	file, err := os.ReadFile("./output.json")

	fileLoadTimeSec = time.Since(start).Seconds()

	if err != nil {
		log.Fatal(err)
	}
	parser := Parser{
		data: file,
		pos:  0,
	}
	var output interface{}
	if parser.getByte() == '[' {
		output, err = parser.parseArray()
	} else if parser.getByte() == '{' {
		output, err = parser.parseObject()
	} else {
		panic("not supported char")
	}

	var data DataStruct

	if err := assign(output, &data); err != nil {
		panic(err)
	}

	return output, len(file), fileLoadTimeSec
}

func assign(data any, out any) error {
	val := reflect.ValueOf(out)
	if val.Kind() != reflect.Ptr || val.IsNil() {
		return errors.New("output must be a non-nil pointer")
	}

	return assignValue(data, val.Elem())
}

func assignValue(data any, target reflect.Value) error {
	if !target.CanSet() {
		return errors.New("target is not settable")
	}

	switch target.Kind() {
	case reflect.Struct:
		m, ok := data.(map[string]interface{})
		if !ok {
			return fmt.Errorf("expected map[string]interface{} for struct, got %T", data)
		}

		for i := 0; i < target.NumField(); i++ {
			field := target.Type().Field(i)

			// Skip unexported fields
			if !target.Field(i).CanSet() {
				continue
			}

			jsonKey := field.Tag.Get("json")
			if jsonKey == "" {
				jsonKey = field.Name
			}

			val, found := m[jsonKey]
			if !found {
				return fmt.Errorf("missing field '%s' in input", jsonKey)
			}

			if err := assignValue(val, target.Field(i)); err != nil {
				return fmt.Errorf("error assigning field '%s': %w", jsonKey, err)
			}
		}

		return nil

	case reflect.Slice:
		arr, ok := data.([]interface{})
		if !ok {
			return errors.New("expected array for slice assignment")
		}

		slice := reflect.MakeSlice(target.Type(), len(arr), len(arr))
		for i, v := range arr {
			if err := assignValue(v, slice.Index(i)); err != nil {
				return err
			}
		}
		target.Set(slice)
		return nil

	case reflect.Float64:
		f, ok := data.(float64)
		if !ok {
			return errors.New("expected float64")
		}
		target.SetFloat(f)
		return nil

	case reflect.String:
		s, ok := data.(string)
		if !ok {
			return errors.New("expected string")
		}
		target.SetString(s)
		return nil

	case reflect.Bool:
		b, ok := data.(bool)
		if !ok {
			return errors.New("expected bool")
		}
		target.SetBool(b)
		return nil

	default:
		return errors.New("unsupported type: " + target.Kind().String())
	}
}

func (parser *Parser) parseArray() ([]interface{}, error) {
	if parser.getByte() != '[' {
		panic("expected to start with {")
	}
	parser.nextByte()
	parser.skipWhitespace()

	arr := make([]interface{}, 5_000_000)
	index := 0

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
		arr[index] = val
		index++

		parser.skipWhitespace()

		if parser.getByte() == ',' {
			parser.nextByte()
			parser.skipWhitespace()
		} else if parser.getByte() != ']' {
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
		if err != nil {
			panic(err)
		}
		obj[key] = val

		parser.skipWhitespace()
		if parser.getByte() == ',' {
			parser.nextByte()
			parser.skipWhitespace()
		} else if parser.getByte() != '}' {
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
