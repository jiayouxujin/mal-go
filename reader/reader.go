package reader

import (
	"fmt"
	"github.com/jiayouxujin/mal-go/types"
	"regexp"
	"strconv"
)

var (
	tokenRegexp = `[\s,]*(~@|[\[\]{}()'` + "`" +
		`~^@]|"(?:\\.|[^\\"])*"?|;.*|[^\s\[\]{}('"` + "`" + `,;)]*)`
)

// Reader Next() return the token at the current position and increments the position
//Peek() return the token at the current position
type Reader interface {
	Next() (string, error)
	Peek() (string, error)
}

type TokenReader struct {
	tokens   []string
	position int
}

func (r *TokenReader) anyError() error {
	if r == nil {
		return fmt.Errorf("tokenReader is nil")
	}
	if r.position >= len(r.tokens) {
		return fmt.Errorf("the value of position isn't correct")
	}
	return nil
}
func (r *TokenReader) Next() (string, error) {
	if err := r.anyError(); err != nil {
		return "", err
	}
	res := r.tokens[r.position]
	r.position++
	return res, nil
}

func (r *TokenReader) Peek() (string, error) {
	if err := r.anyError(); err != nil {
		return "", err
	}
	return r.tokens[r.position], nil
}

func ReadStr(input string) (types.MalType, error) {
	//cal tokenize
	tokens, err := tokenize(input)
	if err != nil {
		return nil, err
	}
	if len(tokens) == 0 {
		return nil, fmt.Errorf("empty input")
	}
	//create a new Reader instance
	tr := TokenReader{tokens: tokens, position: 0}
	return readForm(&tr)
}

func tokenize(token string) ([]string, error) {
	re, err := regexp.Compile(tokenRegexp)
	if err != nil {
		return nil, err
	}
	tokens := make([]string, 0)
	for _, group := range re.FindAllStringSubmatch(token, -1) {
		tmp := group[1]
		//ignore whitespaces or commas
		if tmp == "" || tmp[0] == ';' {
			continue
		}
		tokens = append(tokens, tmp)
	}
	return tokens, nil
}

func readForm(rd Reader) (types.MalType, error) {
	token, err := rd.Peek()
	if err != nil {
		return nil, err
	}
	switch token {
	case "(":
		return readList(rd)
	case ")":
		return nil, fmt.Errorf("unexpected ')'")
	case "[":
		return readVector(rd)
	case "]":
		return nil, fmt.Errorf("unexpected ']'")
	case "{":
		return readHashmap(rd)
	case "}":
		return nil, fmt.Errorf("unexpected '}'")
	default:
		return readAtom(rd)
	}
}

func readAtom(r Reader) (types.MalType, error) {
	token, err := r.Next()
	if err != nil {
		return nil, err
	}
	if matched, _ := regexp.MatchString(`^[-+]?\d+$`, token); matched { //number
		number, err := strconv.Atoi(token)
		if err != nil {
			return nil, err
		}
		return types.MalNumber{Value: number}, nil
	} else if matched, _ := regexp.MatchString(`^"(?:\\.|[^\\"])*"?$`, token); matched { // string
		if matched, _ := regexp.MatchString(`^"(?:\\.|[^\\"])*"$`, token); !matched {
			return nil, fmt.Errorf("unclosed string: %s", token)
		}
		unquoted, err := strconv.Unquote(token)
		if err != nil {
			return nil, err
		}
		return types.MalString{Value: unquoted}, nil
	} else if token == "nil" {
		return types.MalNil, nil
	} else if token == "true" {
		return types.MalTrue, nil
	} else if token == "false" {
		return types.MalFalse, nil
	} else if token[0] == ':' {
		return types.MalKeyword{Value: token[1:]}, nil
	} else {
		return types.MalSymbol{Value: token}, nil
	}
}

func readList(r Reader) (types.MalType, error) {
	return readStartEnd(r, "(", ")")
}

func readVector(r Reader) (types.MalType, error) {
	list, err := readStartEnd(r, "[", "]")
	if err != nil {
		return nil, err
	}
	return list, nil
}

func readHashmap(r Reader) (types.MalType, error) {
	list, err := readStartEnd(r, "{", "}")
	if err != nil {
		return nil, err
	}
	if len(list)%2 != 0 {
		return nil, fmt.Errorf("incorrect number of elements for a hashmap")
	}
	hashmap := make(types.MalHashmap)
	for i := 0; i < len(list); i += 2 {
		switch t := list[i].(type) {
		case types.MalKeyword, types.MalString:
			hashmap[t] = list[i+1]
		default:
			return nil, fmt.Errorf("hashmap keys only accept string or keyword")
		}
	}
	return hashmap, nil
}

func readStartEnd(rd Reader, start, end string) (types.MalList, error) {
	first, _ := rd.Next()
	if first != start {
		return nil, fmt.Errorf("incorrect starting token: expect '%s' but get '%s'", start, first)
	}
	astList := types.MalList{}
	for token, err := rd.Peek(); token != end; token, err = rd.Peek() {
		if err != nil {
			return nil, err
		}
		ast, err := readForm(rd)
		if err != nil {
			return nil, err
		}
		astList = append(astList, ast)
	}
	_, _ = rd.Next()
	return astList, nil
}
