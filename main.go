package main

import (
	"bufio"
	"fmt"
	"os"
)

func READ() (string, error) {
	reader := bufio.NewReader(os.Stdin)
	text, err := reader.ReadString('\n')
	if err != nil {
		return "", err
	}
	return text, nil
}

func EVAL(text string) string {
	return text
}

func PRINT(text string) {
	fmt.Print(text)
}

func rep() {
	text, _ := READ()
	res := EVAL(text)
	PRINT(res)
}

func main() {
	for {
		fmt.Print("user> ")
		rep()
	}
}
