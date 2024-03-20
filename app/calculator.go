package main

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

/*
( => 40
) => 41
* => 42
+ => 43
- => 45
. => 46
/ => 47
0 - 9 => 48 - 57
*/
const (
	TypeNone int = iota
	// Number 0-9 .
	TypeNumber
	// Operator +*/
	TypeOperator
	// Operator -
	TypeMinus
)

type Calculator struct {
	// 小數點保留位數
	digits int
}

func NewCalculator() *Calculator {
	return &Calculator{
		digits: 5,
	}
}

func (c *Calculator) SetDigits(digits int) {
	if digits < 0 {
		digits = 0
	}
	c.digits = digits
}

// a + (b - c) * (d - e) / f
func (c Calculator) ComputeV2(formula string) (string, error) {
	left := 0
	right := 0
	length := len(formula)
	var index int
	var f string

	// Find parentheses
	for i := 0; i < length; i++ {
		f = string(formula[i])
		switch f {
		case "(":
			if left == 0 {
				index = i
			}
			left += 1
		case ")":
			right += 1
			if left < right {
				return "", errors.New("error format of formula")
			} else if left == right {
				value, err := c.ComputeV2(formula[index+1 : i])
				if err != nil {
					return "", err
				}
				return c.ComputeV2(formula[:index] + value + formula[i+1:])
			}
		default:
		}
	}
	return c.ComputeV1(formula)
}

// a + b - c * d - e / f
func (c Calculator) ComputeV1(formula string) (string, error) {
	elements, err := c.ParseFormula(formula)
	if err != nil {
		return "", err
	}
	// muliply or divide
	elements, err = c.computeAndSquash(elements, "*", "/")
	if err != nil {
		return "", err
	}
	// plus or minus
	elements, err = c.computeAndSquash(elements, "+", "-")
	if err != nil {
		return "", err
	}
	result := "0"
	var element string
	for _, element = range elements {
		result, err = c.compute(result, "+", element)
		if err != nil {
			return "", errors.New("invalid element of formula")
		}
	}
	return result, nil
}

func (c Calculator) ParseFormula(formula string) ([]string, error) {
	elements, err := c.Parse([]string{formula}, "*", "/")
	if err != nil {
		return nil, err
	}
	buffers := []string{}
	var temps []string
	var element string
	for _, element = range elements {
		temps, err = c.Parse([]string{strings.TrimSpace(element)}, "+")
		if err != nil {
			return nil, err
		}
		buffers = append(buffers, temps...)
	}
	elements = elements[:0]
	for _, element = range buffers {
		temps, err = c.ParseMinus(strings.TrimSpace(element))
		if err != nil {
			return nil, err
		}
		elements = append(elements, temps...)
	}
	return elements, nil
}

func (c Calculator) Parse(elements []string, delimiters ...string) ([]string, error) {
	if len(delimiters) == 0 {
		return elements, nil
	}
	results := []string{}
	delimiter := rune(delimiters[0][0])
	delimiters = delimiters[1:]
	var element string
	var i, index int
	var e rune
	for _, element = range elements {
		index = 0
		for i, e = range element {
			if e == delimiter {
				results = append(results, element[index:i])
				results = append(results, string(e))
				index = i + 1
			}
		}
		if index != len(element) {
			results = append(results, element[index:])
		}
	}
	if len(delimiters) > 0 {
		return c.Parse(results, delimiters...)
	}
	return results, nil
}

// rune: - => 45
func (c Calculator) ParseMinus(formula string) ([]string, error) {
	count := 0
	var r rune
	for _, r = range formula {
		if r == 45 {
			count += 1
		}
	}
	if count == 0 {
		return []string{formula}, nil
	}
	elements := []string{}
	rs := []rune{}
	var currType int = TypeNone
	var crossType int
	var t int
	var err error
	for _, r = range formula {
		t = c.GetType(r)
		// None, Number, Minus
		crossType = 10*currType + t
		switch crossType {
		// None -> None
		case 0:
			// NOTE: Do nothing
		// None -> Number
		case 1:
			rs = append(rs, r)
		// None -> Minus
		case 3:
			rs = append(rs, r)
		// Number -> None
		case 10:
			elements, rs, err = c.squashMinus(elements, rs)
			if err != nil {
				return nil, err
			}
		// Number -> Number
		case 11:
			rs = append(rs, r)
		// Number -> Minus
		case 13:
			elements, rs, err = c.squashMinus(elements, rs)
			if err != nil {
				return nil, err
			}
			rs = append(rs, r)
		// Minus -> None
		case 30:
			elements, rs, err = c.squashMinus(elements, rs)
			if err != nil {
				return nil, err
			}
		// Minus -> Number
		case 31:
			rs = append(rs, r)
		// Minus -> Minus
		case 33:
			if len(rs) <= 1 {
				rs = append(rs, r)
			} else {
				return nil, errors.New("error format of formula")
			}
		default:
			return nil, errors.New("error format of formula")
		}
		currType = t
	}
	if len(rs) > 0 {
		elements, _, err = c.squashMinus(elements, rs)
		if err != nil {
			return nil, err
		}
	}
	return elements, nil
}

func (c Calculator) squashMinus(elements []string, rs []rune) ([]string, []rune, error) {
	// 0: number only
	// 1: negative number (-x)
	// 2: minus a negative number (- -x)
	count := 0
	var i int
	var r rune
	for i, r = range rs {
		if r == 45 {
			if i != 0 && i != 1 {
				count = -1
			}
			count += 1
			if count > 2 {
				count = -1
			}
		}
	}

	switch count {
	case 0:
		fallthrough
	case 1:
		elements = append(elements, string(rs))
	case 2:
		elements = append(elements, "-")
		elements = append(elements, string(rs[1:]))
	default:
		return nil, nil, errors.New("error format of formula")
	}
	rs = rs[:0]
	return elements, rs, nil
}

func (c Calculator) computeAndSquash(elements []string, targets ...string) ([]string, error) {
	index := c.FindIndex(elements, targets...)
	var v1, v2, op, result string
	var prev, next []string
	var err error
	for index != -1 {
		if index == 0 || index == len(elements)-1 {
			return nil, errors.New("error format of formula")
		}
		v1 = elements[index-1]
		op = elements[index]
		v2 = elements[index+1]
		if c.GetElementType(v1) != TypeNumber || c.GetElementType(v2) != TypeNumber {
			return nil, errors.New("error format of formula")
		}
		result, err = c.compute(v1, op, v2)
		if err != nil {
			return nil, errors.New("invalid element of formula")
		}
		prev = elements[:index-1]
		next = elements[index+2:]
		elements = append(prev, result)
		elements = append(elements, next...)
		index = c.FindIndex(elements, targets...)
	}
	return elements, nil
}

func (c Calculator) FindIndex(elements []string, targets ...string) int {
	var t string
	for i, e := range elements {
		for _, t = range targets {
			if e == t {
				return i
			}
		}
	}
	return -1
}

func (c Calculator) compute(v1 string, op string, v2 string) (string, error) {
	var f1, f2 float64
	var err error
	if f1, err = strconv.ParseFloat(v1, 64); err != nil {
		return "", errors.New("cannot convert element to number")
	}
	if f2, err = strconv.ParseFloat(v2, 64); err != nil {
		return "", errors.New("cannot convert element to number")
	}
	var result float64
	switch op {
	case "+":
		result = f1 + f2
	case "-":
		result = f1 - f2
	case "*":
		result = f1 * f2
	case "/":
		result = f1 / f2
	default:
		return "", errors.New("invalid operator")
	}
	format := fmt.Sprintf("%%.%df", c.digits)
	return fmt.Sprintf(format, result), nil
}

func (c Calculator) GetElementType(s string) int {
	rs := []rune(s)
	length := len(rs)
	switch length {
	case 0:
		return TypeNone
	case 1:
		return c.GetType(rs[0])
	default:
		return TypeNumber
	}
}

func (c Calculator) GetType(b rune) int {
	if b == 45 {
		return TypeMinus
	} else if (40 <= b && b <= 43) || b == 47 {
		return TypeOperator
	} else if b == 46 || (48 <= b && b <= 57) {
		return TypeNumber
	} else {
		return TypeNone
	}
}
