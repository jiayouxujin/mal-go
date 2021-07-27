package printer

import (
	"github.com/jiayouxujin/mal-go/types"
	"strconv"
)

func printList(lst types.MalList, start, end string, readable bool) string {
	result := start
	for index, item := range lst {
		if index != 0 {
			result += " "
		}
		result += PrStr(item, readable)
	}
	result += end
	return result
}

func printHashmap(hm types.MalHashmap, readable bool) string {
	result := "{"
	flag := true
	for k, v := range hm {
		if !flag {
			result += " "
		}
		flag = false
		result += PrStr(k, readable)
		result += " "
		result += PrStr(v, readable)
	}
	result += "}"
	return result
}
func PrStr(data types.MalType, readable bool) string {
	switch t := data.(type) {
	case types.MalNumber:
		return strconv.Itoa(t.Value)
	case types.MalSymbol:
		return t.Value
	case types.MalString:
		if readable {
			return strconv.Quote(t.Value)
		}
		return t.Value
	case types.MalLiteral:
		return string(t)
	case types.MalKeyword:
		return ":" + t.Value
	case types.MalList: //(foo bar)
		return printList(t, "(", ")", readable)
	case types.MalVector: //[foo bar]
		return printList(types.MalList{t}, "[", "]", readable)
	case types.MalHashmap:
		return printHashmap(t, readable)
	default:
		return "/UNKNOWN VALUE/"
	}
}
