package main

import (
	"errors"
	"fmt"
	"strconv"
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
	// Operator +-*/
	TypeOperator
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

func (c *Calculator) SetDigits(digits int){
	if(digits < 0){
		digits = 0
	}
	c.digits = digits
}

// a + b - c * d - e / f
func (c Calculator) ComputeV1(formula string) (string, error) {
	elements, err := c.ParseFormula(formula)
	if err != nil{
		return "", err
	}
	// muliply or divide
	elements, err = c.computeAndSquash(elements, "*", "/")
	if err != nil{
		return "", err
	}
	// plus or minus
	elements, err = c.computeAndSquash(elements, "+", "-")
	if err != nil{
		return "", err
	}
	if len(elements) != 1{
		fmt.Printf("elements: %+v\n", elements)
		return "", errors.New("unexception elements in formula")
	}
	return elements[0], nil
}

func (c Calculator) computeAndSquash(elements []string, targets ...string)([]string, error){
	index := c.FindIndex(elements, targets...)
	var v1, v2, op, result string
	var prev, next []string
	var err error
	for index != -1{
		if index == 0 || index == len(elements) - 1{
			return nil, errors.New("error format of formula")
		}
		v1 = elements[index - 1]
		op = elements[index]
		v2 = elements[index + 1]
		if c.GetElementType(v1) != TypeNumber || c.GetElementType(v2) != TypeNumber{
			return nil, errors.New("error format of formula")
		}
		result, err = c.compute(v1, op, v2)
		if err != nil{
			return nil, errors.New("invalid element of formula")
		}
		prev = elements[:index - 1]
		next = elements[index + 2:]
		elements = append(prev, result)
		elements = append(elements, next...)
		index = c.FindIndex(elements, targets...)
	}
	return elements, nil
}

func (c Calculator) compute(v1 string, op string, v2 string)(string, error){
	var f1, f2 float64
	var err error
	if f1, err = strconv.ParseFloat(v1, 64); err != nil{
		return "", errors.New("cannot convert element to number")
	}
	if f2, err = strconv.ParseFloat(v2, 64); err != nil{
		return "", errors.New("cannot convert element to number")
	}
	var result float64
	switch(op){
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
	// fmt.Printf("format: %s, result: %f\n", format, result)
	return fmt.Sprintf(format, result), nil
}

func (c Calculator) FindIndex(elements []string, targets ...string)int{
	var t string
	for i, e := range elements{
		for _, t = range targets{
			if e == t{
				return i
			}
		}
	}
	return -1
}

func (c Calculator) ParseFormula(formula string) ([]string, error){
	var currType, t int = TypeNone, TypeNone
	var crossType int
	elements := []string{}
	rs := []rune{}
	for _, f := range formula {
		t = c.GetType(f)
		// fmt.Printf("'%s' -> type: %d\n", string(f), t)
		crossType = 10*currType + t
		// fmt.Printf("currType: %d, crossType: %d\n", currType, crossType)
		switch crossType {
		// None -> Number
		case 1:
			fallthrough
		// Operator -> Number
		case 21:
			rs = append(rs, f)
		// None -> Operator
		case 2:
			elements = append(elements, string(f))
			// fmt.Printf("Append: '%s'\n", string(f))
		// Number -> None
		case 10:
			elements = append(elements, string(rs))
			rs = rs[:0]
		// Number -> Number
		case 11:
			rs = append(rs, f)
		// Number -> Operator
		case 12:
			elements = append(elements, string(rs))
			rs = rs[:0]
			elements = append(elements, string(f))
		// Operator -> Operator
		case 22:
			return nil, errors.New("duplicate operator")
		// None -> None; Operator -> None
		default:
		}
		currType = t
	}
	if(len(rs) != 0){
		if currType == TypeOperator{
			return nil, errors.New("formula end with operator")
		}
		elements = append(elements, string(rs))
	}
	return elements, nil
}

func (c Calculator) GetElementType(s string) int {
	rs := []rune(s)
	if len(rs) > 1{
		return TypeNumber
	}else{
		return c.GetType(rs[0])
	}
}

func (c Calculator) GetType(b rune) int {
	if (40 <= b && b <= 43) || b == 45 || b == 47 {
		return TypeOperator
	} else if b == 46 || (48 <= b && b <= 57) {
		return TypeNumber
	} else {
		return TypeNone
	}
}

// a + (b - c) * (d - e) / f
func (c Calculator) Compute(formula string) (string, error) {
	// left := 0
	// right := 0
	// length := len(formula)
	// var index int
	// var f string

	// // Find parentheses
	// for i := 0; i < length; i++ {
	// 	f = string(formula[i])
	// 	switch f {
	// 	case "(":
	// 		if left == 0 {
	// 			index = i
	// 		}
	// 		left += 1
	// 	case ")":
	// 		right += 1
	// 		if left < right {
	// 			return "", errors.New("Error format of formula.")
	// 		} else if left == right {
	// 			value, err := c.Compute()
	// 		}
	// 	default:
	// 	}
	// }
	return formula, nil
}
