package tokenizer

import (
	"bytes"
	"fmt"
	"strings"
	"unicode"
)

type TokenType int

const (
	ItemStart TokenType = iota
	ItemEnd
	ItemId
	ItemKey
	ItemValue
	ListStart
	ListEnd
	ListValue
	NullValue
)

func (t TokenType) String() string {
	return [...]string{"ItemStart", "ItemEnd", "ItemId", "ItemKey", "ItemValue", "ListStart", "ListEnd", "ListValue", "NullValue"}[t]
}

type Token struct {
	location int
	Depth    int
	Token    TokenType
	Value    string
}

func (t Token) String() string {
	return fmt.Sprintf("{%v %v '%v'}", t.Depth, t.Token, t.Value)
}

func Tokenize(plan []rune) []Token {
	acc := []Token{}

	depth := 0

	for i := 0; i < len(plan); i++ {
		if plan[i] == '{' {
			depth++
			acc = append(acc, Token{i, depth, ItemStart, "{"})
			continue
		}

		if plan[i] == '}' {
			acc = append(acc, Token{i, depth, ItemEnd, "}"})
			depth--
			continue
		}

		if plan[i] == '(' {
			depth++
			acc = append(acc, Token{i, depth, ListStart, "("})
			continue
		}

		if plan[i] == ')' {
			acc = append(acc, Token{i, depth, ListEnd, ")"})
			depth--
			continue
		}

		if isNull(&i, plan) {
			acc = append(acc, Token{i, depth, NullValue, "<>"})
			continue
		}

		yesWord, word := isWord(&i, plan)
		if yesWord && isItemStart(acc) {
			acc = append(acc, Token{i, depth, ItemId, word})
			continue
		}

		if yesWord && isItemKey(acc) {
			acc = append(acc, Token{i, depth, ItemValue, word})
			continue
		}

		if yesWord && isListContext(acc) {
			acc = append(acc, Token{i, depth, ListValue, word})
			continue
		}

		if yesWord && !isListContext(acc) && !isItemKey(acc) {
			acc = append(acc, Token{i, depth, ItemKey, word})
			continue
		}
	}

	return acc
}

func isItemStart(tokens []Token) bool {
	lastToken := tokens[len(tokens)-1]
	return lastToken.Token == ItemStart
}

func isItemKey(tokens []Token) bool {
	lastToken := tokens[len(tokens)-1]
	return lastToken.Token == ItemKey
}

func isListContext(tokens []Token) bool {
	for j := len(tokens) - 1; j > 0; j-- {
		currentToken := tokens[j]
		if currentToken.Token == ListStart {
			return true
		}

		if currentToken.Token == ItemStart {
			return false
		}

		if currentToken.Token == ListEnd {
			return false
		}
	}
	return false
}

func isWord(i *int, plan []rune) (bool, string) {
	var b bytes.Buffer
	var j int
	for j = *i; j < len(plan); j++ {
		c := plan[j]
		if unicode.IsLetter(c) || unicode.IsNumber(c) || c == ':' || c == '.' || c == '_' {
			b.WriteRune(c)
		} else {
			break
		}
	}

	if b.Len() > 0 {
		*i = j - 1
		return true, b.String()
	} else {
		return false, ""
	}
}

func isNull(i *int, plan []rune) bool {
	if plan[*i] == '<' && plan[*i+1] == '>' {
		*i = *i + 1
		return true
	}
	return false
}

func IsKey(token Token, keyName string) bool {
	return token.Token == ItemKey && strings.Contains(token.Value, keyName)
}
