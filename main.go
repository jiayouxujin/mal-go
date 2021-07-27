package main

import (
	"bufio"
	"fmt"
	"github.com/jiayouxujin/mal-go/printer"
	"github.com/jiayouxujin/mal-go/reader"
	"github.com/jiayouxujin/mal-go/types"
	"os"
)

func READ() (types.MalType, error) {
	in := bufio.NewReader(os.Stdin)
	fmt.Print("user> ")
	text, err := in.ReadString('\n')
	if err != nil {
		return "", err
	}
	ast, err := reader.ReadStr(text)
	if err != nil {
		return "", err
	}
	return ast, nil
}

func EVAL(ast types.MalType) (types.MalType, error) {
	return ast, nil
}

func PRINT(exp types.MalType) (string, error) {
	return printer.PrStr(exp, true), nil
}

func rep() (types.MalType, error) {
	var exp types.MalType
	var res string
	var e error
	if exp, e = READ(); e != nil {
		return nil, e
	}
	if exp, e = EVAL(exp); e != nil {
		return nil, e
	}
	if res, e = PRINT(exp); e != nil {
		return nil, e
	}
	return res, nil
}

func main() {
	for {
		out, err := rep()
		if err != nil {
			fmt.Printf("Error: %v\n", err)
		}
		fmt.Printf("%v\n", out)
	}
}
