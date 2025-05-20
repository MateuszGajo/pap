package main

import (
	"encoding/binary"
	"fmt"
	"log"
	"os"
	"strings"
)

var registers_8 = []string{"AL", "CL", "DL", "BL", "AH", "CH", "DH", "BH"}
var registers_16 = []string{"AX", "CX", "DX", "BX", "SP", "BP", "SI", "DI"}
var effectiveAddress = [][]string{{"BX", "SI"}, {"BX", "DI"}, {"BP", "SI"}, {"BP", "DI"}, {"SI"}, {"DI"}, {"BP"}, {"BX"}}

func buildValue(low, high byte) int {
	val := int16(uint16(high)<<8 | uint16(low))
	return -int(val)
}

func main() {
	data, err := os.ReadFile("listing_0037_single_register_mov")

	if err != nil {
		log.Fatal(err)
	}

	// fmt.Println(data)
	// fmt.Printf("%08b\n", data)

	for len(data) > 0 {
		firstByte := data[0]

		fmt.Printf("%0b \n", data[0])
		if len(data) > 0 {
			data = data[1:]
		}

		// mov operation register/memory to/from register
		if firstByte>>2 == 0b00100010 {

			secondByte := data[0]
			if len(data) > 0 {
				data = data[1:]
			}

			d := (firstByte >> 1) & 1

			w := firstByte & 1
			mod := (secondByte >> 6) & 0b11
			reg := (secondByte >> 3) & 0b111
			rm := secondByte & 0b111

			registers := registers_8
			if w == 0b1 {
				registers = registers_16
			}

			var dst, src string
			if mod == 0b11 {
				dst = registers[rm]
				src = registers[reg]
				if d == 0b1 {
					dst, src = src, dst
				}

			} else if mod == 0b00 {
				dst = strings.Join(effectiveAddress[rm], "+")
				src = registers[reg]

			} else if mod == 0b01 {
				thirdByte := data[0]
				if len(data) > 0 {
					data = data[1:]
				}
				val := int8(thirdByte)
				dst = strings.Join(effectiveAddress[rm], "+")

				if val > 0 {
					dst += fmt.Sprintf("+%d", val)
				} else if val < 0 {
					dst += fmt.Sprintf("%d", val)
				}

				src = registers[reg]

			} else if mod == 0b10 {
				thirdByte := data[0]
				fourthByte := data[1]
				val := int16(binary.LittleEndian.Uint16(append([]byte{}, thirdByte, fourthByte)))
				if len(data) > 0 {
					data = data[2:]
				}

				dst = strings.Join(effectiveAddress[rm], "+")
				if val > 0 {
					dst += fmt.Sprintf("+%d", val)
				} else if val < 0 {
					dst += fmt.Sprintf("%d", val)
				}

				src = registers[reg]

			} else {
				fmt.Printf("%0b \n", mod)
				fmt.Printf("%0b \n", rm)
				panic("not supported mod")
			}

			if d == 0b1 {
				dst, src = src, dst
			}

			fmt.Printf("mov %v, %v \n", dst, src)
		} else if firstByte>>1 == 0b01100011 {

			secondByte := data[0]
			if len(data) > 0 {
				data = data[1:]
			}

			w := firstByte & 1
			mod := (secondByte >> 6) & 0b11
			rm := secondByte & 0b111

			registers := registers_8
			if w == 0b1 {
				registers = registers_16
			}

			if w == 0b1 {
				registers = registers_16
			}

			var res string

			var dst string
			if mod == 0b11 {
				dst = registers[rm]
				// src = registers[reg]

			} else if mod == 0b00 {
				dst = strings.Join(effectiveAddress[rm], "+")

				// src = registers[reg]

			} else if mod == 0b01 {
				thirdByte := data[0]
				if len(data) > 0 {
					data = data[1:]
				}
				val := int8(thirdByte)
				dst = strings.Join(effectiveAddress[rm], "+")

				if val > 0 {
					dst += fmt.Sprintf("+%d", val)
				} else if val < 0 {
					dst += fmt.Sprintf("%d", val)
				}

				// src = registers[reg]

			} else if mod == 0b10 {
				thirdByte := data[0]
				fourthByte := data[1]
				val := int16(binary.LittleEndian.Uint16(append([]byte{}, thirdByte, fourthByte)))
				if len(data) > 0 {
					data = data[2:]
				}

				dst = strings.Join(effectiveAddress[rm], "+")
				if val > 0 {
					dst += fmt.Sprintf("+%d", val)
				} else if val < 0 {
					dst += fmt.Sprintf("%d", val)
				}

				// src = registers[reg]

			} else {
				fmt.Printf("%0b \n", mod)
				fmt.Printf("%0b \n", rm)
				panic("not supported mod")
			}

			if w == 0b1 {

				val := int16(binary.LittleEndian.Uint16(append([]byte{}, data[0], data[1])))
				res = fmt.Sprintf("word %d", val)
				if len(data) > 0 {
					data = data[2:]
				}
			} else {
				res = fmt.Sprintf("byte %d", int(int8(data[0])))
				if len(data) > 0 {
					data = data[1:]
				}
			}

			fmt.Printf("mov %v, %v \n", dst, res)
		} else if firstByte>>4 == 0b00001011 {

			w := firstByte >> 3 & 1
			reg := firstByte & 0b111

			registers := registers_8
			secondByte := data[0]

			if len(data) > 0 {
				data = data[1:]
			}

			res := int(int8(secondByte))

			if w == 0b1 {
				registers = registers_16
				thirdByte := data[0]
				if len(data) > 0 {
					data = data[1:]
				}
				val := int16(binary.LittleEndian.Uint16(append([]byte{}, secondByte, thirdByte)))
				res = int(val)
			}
			dst := registers[reg]
			fmt.Printf("mov %v, %v\n", dst, res)

		} else {
			fmt.Printf("%0b \n", firstByte)
			panic("Not supported byte code")
		}

	}
}
