package env

import (
	"fmt"
	"github.com/jiayouxujin/mal-go/core"
	"github.com/jiayouxujin/mal-go/types"
)

type Env struct {
	outer types.MalEnv
	data  map[string]types.MalType
}

//Set takes a symbol key and a mal value and adds to the data structure
func (e *Env) Set(Key types.MalSymbol, value types.MalType) error {
	if e == nil {
		return fmt.Errorf("set value in nil environment")
	}
	e.data[Key.Value] = value
	return nil
}

//Find takes a symbol key and if the current env contains that key then return the env
//if no key is found and outer is not nil then call find on the outer env
func (e *Env) Find(key types.MalSymbol) types.MalEnv {
	if e == nil {
		return nil
	}
	if _, ok := e.data[key.Value]; ok {
		return e
	} else if e.outer != nil {
		return e.outer.Find(key)
	} else {
		return nil
	}
}

//Get takes a symbol key and uses the find method to locate the env with the key,then return the matching the value
func (e *Env) Get(key types.MalSymbol) (types.MalType, error) {
	env := e.Find(key)
	if env == nil {
		return nil, fmt.Errorf("not found")
	} else {
		return env.(*Env).data[key.Value], nil
	}
}

func CreateEnv(outer types.MalEnv, binds types.MalList, exps types.MalList) (*Env, error) {
	env := &Env{
		outer: outer,
		data:  make(map[string]types.MalType),
	}

	flag := false
	for i, k := range binds {
		symbol, ok := k.(types.MalSymbol)
		if !ok {
			return nil, fmt.Errorf("invalid symbol(s) in variable bindings")
		}
		if symbol.Value == "&" {
			if i != len(binds)-2 {
				return nil, fmt.Errorf("invalid position for '&' in bindings")
			}
			flag = true
		}
	}

	if flag && len(binds)-2 > len(exps) {
		return nil, fmt.Errorf("not enough expressions for a variadic function")
	} else if !flag && len(binds) != len(exps) {
		return nil, fmt.Errorf(
			"different numbers of bindings and expressions for a non-variadic function")
	}
	for i, k := range binds {
		if flag && i == len(binds)-2 {
			continue
		}
		symbol, _ := k.(types.MalSymbol)
		var v types.MalType
		if flag && i == len(binds)-1 {
			v = exps[i-1:]
		} else {
			v = exps[i]
		}
		err := env.Set(symbol, v)
		if err != nil {
			return nil, err
		}
	}
	return env, nil
}

func GetInitEnv() (e *Env) {
	e, _ = CreateEnv(nil, nil, nil)
	for k, v := range core.NameSpace {
		err := e.Set(types.MalSymbol{Value: k}, v)
		if err != nil {
			return nil
		}
	}
	return
}
