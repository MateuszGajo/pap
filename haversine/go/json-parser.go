package main

import (
	"errors"
	"fmt"
	"log"
	"os"
	"reflect"
	"time"
	"unsafe"
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
		// a lot of allocations, maybe do it once in a while at the end?
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
	//maybe sliding windows instead of parsing byte by byte
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
		//maybe one method skipp white spaces, get next byte?
		parser.nextByte()

		parser.skipWhitespace()
		//inlining maybe?
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

func isDigit(ch byte) bool {
	return ch >= '0' && ch <= '9'
}

func fastParseFloat(b []byte) (float64, error) {
	var i int
	var neg bool
	if b[0] == '-' {
		neg = true
		i++
	}
	var intPart int64
	for i < len(b) && isDigit(b[i]) {
		intPart = intPart*10 + int64(b[i]-'0')
		i++
	}
	var fracPart float64
	if i < len(b) && b[i] == '.' {
		i++
		divisor := 10.0
		for i < len(b) && isDigit(b[i]) {
			fracPart += float64(b[i]-'0') / divisor
			divisor *= 10
			i++
		}
	}
	f := float64(intPart) + fracPart
	if neg {
		f = -f
	}
	return f, nil
}

func (parser *Parser) parseNumber() (interface{}, error) {
	start := parser.pos
	// fmt.Println("start", start)
	if parser.getByte() == '-' {
		parser.nextByte()
	}
	digitsFound := false

	for (len(parser.data) > parser.pos) && isDigit(parser.getByte()) {
		parser.nextByte()
		digitsFound = true
	}

	if !digitsFound {
		return nil, errors.New("invalid number")
	}

	if parser.getByte() == '.' {
		parser.nextByte()
		if !isDigit(parser.getByte()) {
			return nil, errors.New("invalid float format")
		}
		for (len(parser.data) > parser.pos) && isDigit(parser.getByte()) {
			parser.nextByte()
		}
	}

	// fmt.Println("end", parser.pos)

	// diff := parser.pos - start

	num, err := fastParseFloat(parser.data[start:parser.pos])
	if err != nil {
		return nil, fmt.Errorf("invalid number: %v", err)
	}
	return num, nil
}

// [{"x0":-68.10232982601129,"
func (parser *Parser) parseValue() (interface{}, error) {
	// is if else faster?
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

// inlining
func (p *Parser) skipWhitespace() {
	for p.pos < len(p.data) && (p.data[p.pos] == ' ' || p.data[p.pos] == '\n' || p.data[p.pos] == '\r' || p.data[p.pos] == '\t') {
		p.pos++
	}
}

func (parser *Parser) parseString() (string, error) {
	if parser.getByte() != '"' {
		return "", fmt.Errorf("Expected key to start with \"")
	}
	parser.nextByte()
	start := parser.pos
	for {
		if parser.getByte() == '"' {
			key := unsafeBytesToString(parser.data[start:parser.pos])
			parser.nextByte()

			return key, nil
		}
		parser.nextByte()
	}

	panic("should never enter here")
}

func unsafeBytesToString(b []byte) string {
	return *(*string)(unsafe.Pointer(&b))
}
