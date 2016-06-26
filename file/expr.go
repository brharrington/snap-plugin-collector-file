package file

import (
	"strings"
	"fmt"
	"errors"
	"math"
)

func eval(vars map[string]interface{}, expr string) (interface{}, error) {
	var err error
	stack := []interface{}{}
	parts := strings.Split(expr, ",")
	for _, part := range parts {
		switch strings.Trim(part, " \t\r\n") {
		case ":add":
			stack, err = binaryOp(stack, add)
			if err != nil {
				return nil, err
			}
		case ":sub":
			stack, err = binaryOp(stack, sub)
			if err != nil {
				return nil, err
			}
		case ":mul":
			stack, err = binaryOp(stack, mul)
			if err != nil {
				return nil, err
			}
		case ":div":
			stack, err = binaryOp(stack, div)
			if err != nil {
				return nil, err
			}
		default:
			var v interface{} = part
			vlen := len(part)
			if vlen > 2 && part[0] == '{' && part[vlen - 1] == '}' {
				varname := part[1:vlen - 1]
				if varvalue, ok := vars[varname]; ok {
					v = varvalue
				} else {
					return nil, errors.New(fmt.Sprintf("unknown variable: '%v'", varname))
				}
			} else {
				v = parseValue(part)
			}
			stack = append(stack, v)
		}
	}

	if len(stack) != 1 {
		msg := fmt.Sprintf("stack should have single item, found: %v", stack)
		return nil, errors.New(msg)
	}

	return stack[0], nil
}

func add(v1 float64, v2 float64) float64 {
	return v1 + v2
}

func sub(v1 float64, v2 float64) float64 {
	return v1 - v2
}

func mul(v1 float64, v2 float64) float64 {
	return v1 * v2
}

func div(v1 float64, v2 float64) float64 {
	return v1 / v2
}

func binaryOp(stack []interface{}, op func(float64,float64)float64) ([]interface{}, error) {
	end := len(stack)
	if end < 2 {
		return nil, errors.New(fmt.Sprintf("need at least two arguments on the stack: %v", stack))
	}

	v2, err := toNumber(stack[end - 1])
	if err != nil {
		return nil, err
	}

	v1, err := toNumber(stack[end - 2])
	if err != nil {
		return nil, err
	}

	return append(stack[:end - 2], op(v1, v2)), nil
}

func toNumber(v interface{}) (float64, error) {
	switch i := v.(type) {
	case int:
		return float64(i), nil
	case int8:
		return float64(i), nil
	case int16:
		return float64(i), nil
	case int32:
		return float64(i), nil
	case int64:
		return float64(i), nil
	case uint:
		return float64(i), nil
	case uint8:
		return float64(i), nil
	case uint16:
		return float64(i), nil
	case uint32:
		return float64(i), nil
	case uint64:
		return float64(i), nil
	case float32:
		return float64(i), nil
	case float64:
		return float64(i), nil
	default:
		return math.NaN(), errors.New(fmt.Sprintf("not a number: '%v' %T", v, v))
	}
}