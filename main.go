package main

import (
	"bufio"
	"fmt"
	"github.com/jiayouxujin/mal-go/core"
	"github.com/jiayouxujin/mal-go/printer"
	"github.com/jiayouxujin/mal-go/reader"
	. "github.com/jiayouxujin/mal-go/types"
	"os"
)

var replEnv = map[string]MalType{
	"+": func() {

	},
}

func READ() (MalType, error) {
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

func evalAst(ast MalType, env map[string]MalFunction) (MalType, error) {
	switch t := ast.(type) {
	case MalSymbol:
		if fun, ok := env[t.Value]; ok {
			return fun, nil
		}
		return nil, fmt.Errorf("failed to look up '%s' in environments", t.Value)
	case MalList:
		evaluatedList := make(MalList, 0)
		for _, ori := range t {
			if evaluated, err := EVAL(ori, env); err == nil {
				evaluatedList = append(evaluatedList, evaluated)
			} else {
				return nil, err
			}
		}
		return evaluatedList, nil
	case MalVector:
		evaluatedLIst, err := evalAst(MalList(t), env)
		if err != nil {
			return nil, err
		}
		return MalVector(evaluatedLIst.(MalList)), nil
	case MalHashmap:
		result := make(MalHashmap)
		for k, v := range t {
			v, err := EVAL(v, env)
			if err != nil {
				return nil, err
			}
			result[k] = v
		}
		return result, nil
	default:
		return ast, nil
	}
}
func EVAL(ast MalType, env map[string]MalFunction) (MalType, error) {
	switch t := ast.(type) {
	case MalList:
		if len(t) == 0 {
			return t, nil //ast is empty list return ast unchanged
		} else {
			evaluatedList, err := evalAst(ast, env)
			if err != nil {
				return nil, err
			}
			f := evaluatedList.(MalList)[0].(MalFunction)
			return f(evaluatedList.(MalList)[1:]...)
		}
	default: //ast is not a list,call evalAst
		return evalAst(ast, env)
	}
}

func PRINT(exp MalType) (string, error) {
	return printer.PrStr(exp, true), nil
}

func rep() (MalType, error) {
	var exp MalType
	var res string
	var e error
	if exp, e = READ(); e != nil {
		return nil, e
	}
	if exp, e = EVAL(exp, core.NameSpace); e != nil {
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
			fmt.Printf("%v\n", err)
		} else {
			fmt.Printf("%v\n", out)
		}
	}
}
