package reader

import (
	"errors"
	"fmt"
	. "github.com/jiayouxujin/mal-go/types"
	"regexp"
	"strconv"
)

var (
	tokenRegexp = `[\s,]*(~@|[\[\]{}()'` + "`" +
		`~^@]|"(?:\\.|[^\\"])*"?|;.*|[^\s\[\]{}('"` + "`" + `,;)]*)`
)

//Reader Next() returns the token at the current position and increments the position
//Reader Peek() returns the token at the current posistion
type Reader interface {
	Next() (string, error)
	Peek() (string, error)
}

type TokenReader struct {
	tokens   []string
	position int
}

func (t *TokenReader) anyError() error {
	if t == nil {
		return errors.New("tokenReader is nil")
	}
	if t.position >= len(t.tokens) {
		return errors.New("position out of tokens")
	}
	return nil
}

func (t *TokenReader) Next() (string, error) {
	if err := t.anyError(); err != nil {
		return "", err
	}
	tmp := t.tokens[t.position]
	t.position++
	return tmp, nil
}

func (t *TokenReader) Peek() (string, error) {
	if err := t.anyError(); err != nil {
		return "", err
	}
	return t.tokens[t.position], nil
}

func ReadStr(input string) (MalType, error) {
	//call tokenize
	tokens, err := tokenize(input)
	if err != nil {
		return nil, err
	}
	if len(tokens) == 0 {
		return nil, fmt.Errorf("<empty input>")
	}
	//create a new Reader object
	r := TokenReader{
		tokens:   tokens,
		position: 0,
	}
	//call readForm
	return readForm(&r)
}

func tokenize(input string) ([]string, error) {
	re, err := regexp.Compile(tokenRegexp)
	if err != nil {
		return nil, err
	}
	res := make([]string, 0)
	for _, token := range re.FindAllStringSubmatch(input, -1) {
		tmp := token[1]
		if tmp == "" || tmp[0] == ';' { //ignore whitespaces or commas
			continue
		}
		res = append(res, tmp)
	}
	return res, nil
}

func readForm(t *TokenReader) (MalType, error) {
	first, err := t.Peek()
	if err != nil {
		return nil, err
	}
	switch first[0] {
	case '[':
		return readVector(t)
	case ']':
		return nil, errors.New("unexpected ']'")
	case '(':
		return readList(t)
	case ')':
		return nil, errors.New("unexpected ')")
	case '{':
		return readHashmap(t)
	case '}':
		return nil, errors.New("unexpected '}")
	default:
		return readAtom(t)
	}
}

func readAtom(t *TokenReader) (MalType, error) {
	token, err := t.Next()
	if err != nil {
		return nil, err
	}
	if matched, _ := regexp.MatchString(`^[-+]?\d+$`, token); matched { //number
		number, err := strconv.Atoi(token)
		if err != nil {
			return nil, err
		}
		return MalNumber{Value: number}, nil
	} else if matched, _ := regexp.MatchString(`^"(?:\\.|[^\\"])*"?$`, token); matched { // string
		if matched, _ := regexp.MatchString(`^"(?:\\.|[^\\"])*"$`, token); !matched {
			return nil, fmt.Errorf("unclosed string: %s", token)
		}
		unquoted, err := strconv.Unquote(token)
		if err != nil {
			return nil, err
		}
		return MalString{Value: unquoted}, nil
	} else if token == "nil" {
		return MalNil, nil
	} else if token == "true" {
		return MalTrue, nil
	} else if token == "false" {
		return MalFalse, nil
	} else if token[0] == ':' {
		return MalKeyword{Value: token[1:]}, nil
	} else {
		return MalSymbol{Value: token}, nil
	}
}

func readHashmap(t *TokenReader) (MalType, error) {
	tmp, err := readStartEnd(t, "{", "}")
	if err != nil {
		return nil, err
	}
	if len(tmp)%2 != 0 {
		return nil, fmt.Errorf("the length of hashmap is even,but you get %d\n", len(tmp))
	}
	res := make(MalHashmap)
	for i := 0; i < len(tmp); i += 2 {
		switch t := tmp[i].(type) {
		case MalKeyword, MalString:
			res[t] = tmp[i+1]
		default:
			return nil, fmt.Errorf("hashmap keys only accept string of keyword")
		}
	}
	return res, nil
}

func readList(t *TokenReader) (MalType, error) {
	return readStartEnd(t, "(", ")")
}

func readStartEnd(t *TokenReader, start, end string) (MalList, error) {
	first, _ := t.Next()
	if first != start {
		return nil, fmt.Errorf("unexpected hapend,you want %s,bug get %s", start, first)
	}
	res := MalList{}
	for cur, err := t.Peek(); cur != end; cur, err = t.Peek() {
		if err != nil {
			return nil, err
		}
		tmp, err := readForm(t)
		if err != nil {
			return nil, err
		}
		res = append(res, tmp)
	}
	_, _ = t.Next()
	return res, nil
}

func readVector(t *TokenReader) (MalType, error) {
	tmp, err := readStartEnd(t, "[", "]")
	if err != nil {
		return nil, err
	}
	return MalVector(tmp), nil
}
