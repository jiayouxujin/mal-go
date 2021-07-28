package core

import "github.com/jiayouxujin/mal-go/types"

var NameSpace = map[string]types.MalFunction{
	"+": add,
	"-": sub,
	"*": mul,
	"/": div,
}
