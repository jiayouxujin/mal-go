package main

import (
	"fmt"
	"github.com/jiayouxujin/mal-go/env"
	"github.com/jiayouxujin/mal-go/printer"
	"github.com/jiayouxujin/mal-go/reader"
	"github.com/jiayouxujin/mal-go/readline"
	. "github.com/jiayouxujin/mal-go/types"
)

func READ(input string) (MalType, error) {
	ast, err := reader.ReadStr(input)
	if err != nil {
		return "", err
	}
	return ast, nil
}

func evalAst(ast MalType, env *env.Env) (MalType, error) {
	switch t := ast.(type) {
	case MalSymbol:
		if fun, err := env.Get(t); err == nil {
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
func EVAL(ast MalType, e *env.Env) (MalType, error) {
	for{
		switch t := ast.(type) {
		case MalList:
			if len(t) == 0 {
				return t, nil //ast is empty list return ast unchanged
			}
			first := ""
			if symbol, ok := t[0].(MalSymbol); ok {
				first = symbol.Value
			}
			switch first {
			case "def!":
				if len(t) != 3 {
					return nil, fmt.Errorf("incorrect number of parameters for 'def!'")
				}
				k, ok := t[1].(MalSymbol)
				if !ok {
					return nil, fmt.Errorf("the first parameter is expected to be a symbol")
				}
				v, err := EVAL(t[2], e)
				if err != nil {
					return nil, err
				}
				err = e.Set(k, v)
				return v, err
			case "let*":
				if len(t) != 3 {
					return nil, fmt.Errorf("incorrect number of arguments for 'let*'")
				}
				bindings, ok := t[1].(MalList)
				if !ok || len(bindings)%2 != 0 {
					return nil, fmt.Errorf("the first parameter is expected to be a list of even length")
				}
				tmpEnv, _ := env.CreateEnv(e, nil, nil)
				for i := 0; i < len(bindings); i += 2 {
					k, ok := bindings[i].(MalSymbol)
					if !ok {
						return nil, fmt.Errorf("invalid symbol(s) in variable bindings")
					}
					v, err := EVAL(bindings[i+1], tmpEnv)
					if err != nil {
						return nil, err
					}
					err = tmpEnv.Set(k, v)
					if err != nil {
						return nil, err
					}
				}
				ast, e = t[2], tmpEnv
			default:
				evaluatedList, err := evalAst(t, e)
				if err != nil {
					return nil, err
				}
				return evaluatedList.(MalList)[0].(MalFunction)(evaluatedList.(MalList)[1:]...)
			}
		default: //ast is not a list,call evalAst
			return evalAst(ast, e)
		}
	}
}

func PRINT(exp MalType) (string, error) {
	return printer.PrStr(exp, true), nil
}

func rep(input string, replEnv *env.Env) (MalType, error) {
	var exp MalType
	var res string
	var e error
	if exp, e = READ(input); e != nil {
		return nil, e
	}
	if exp, e = EVAL(exp, replEnv); e != nil {
		return nil, e
	}
	if res, e = PRINT(exp); e != nil {
		return nil, e
	}
	return res, nil
}

func main() {
	defer readline.Close()

	replEnv := env.GetInitEnv()
	for {
		input, err := readline.PromptAndRead("user> ")
		if err != nil {
			break
		}
		res, err := rep(input, replEnv)
		if err != nil {
			fmt.Printf("%v\n", err)
		} else {
			fmt.Printf("%v\n", res)
		}
	}
}
