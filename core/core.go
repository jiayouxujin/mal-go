package core

import (
	"fmt"
	"github.com/jiayouxujin/mal-go/types"
)

func AssertLength(args []types.MalType, expect int) error {
	if actual := len(args); actual != expect {
		return fmt.Errorf("incorrect number of arguments: expect %d but get %d", expect, actual)
	}
	return nil
}
func assertTwoNumbers(args []types.MalType) (int, int, error) {
	if err := AssertLength(args, 2); err != nil {
		return 0, 0, err
	}
	a, ok1 := args[0].(types.MalNumber)
	b, ok2 := args[1].(types.MalNumber)
	if !ok1 || !ok2 {
		return 0, 0, fmt.Errorf("invalid operand(s)")
	}
	return a.Value, b.Value, nil
}
func add(args ...types.MalType) (types.MalType, error) {
	a, b, err := assertTwoNumbers(args)
	if err != nil {
		return nil, err
	}
	return types.MalNumber{
		Value: a + b,
	}, nil
}

func sub(args ...types.MalType) (types.MalType, error) {
	a, b, err := assertTwoNumbers(args)
	if err != nil {
		return nil, err
	}
	return types.MalNumber{
		Value: a - b,
	}, nil
}

func mul(args ...types.MalType) (types.MalType, error) {
	a, b, err := assertTwoNumbers(args)
	if err != nil {
		return nil, err
	}
	return types.MalNumber{
		Value: a * b,
	}, nil
}

func div(args ...types.MalType) (types.MalType, error) {
	a, b, err := assertTwoNumbers(args)
	if err != nil {
		return nil, err
	}
	return types.MalNumber{
		Value: a / b,
	}, nil
}
