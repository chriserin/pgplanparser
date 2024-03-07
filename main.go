package main

import (
	"bytes"
	"fmt"
	"os"

	tkn "github.com/chriserin/plan_parser/tokenizer"
)

type plannedstatement struct {
	plantree plannode
}

type plannode struct {
	nodetype  string
	lefttree  *plannode
	righttree *plannode
}

func (stmt plannedstatement) String() string {
	return fmt.Sprintf("plantree: --- \n %s", stmt.plantree)
}

func (node plannode) String() string {
	var b bytes.Buffer
	if node.nodetype != "" {
		b.WriteString("{ " + node.nodetype + " ")
		if node.lefttree != nil {
			b.WriteString(fmt.Sprintf("%v", node.lefttree))
		}
		if node.righttree != nil {
			b.WriteString(fmt.Sprintf("%v", node.righttree))
		}
		b.WriteString(" }")
		return b.String()
	}
	return ""
}

func main() {
	input := os.Args[1]

	plan := []rune(input)
	value, err := parsePlan(plan)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println(value)
}

func parsePlan(plan []rune) (plannedstatement, error) {
	planTokens := tkn.Tokenize(plan)

	fmt.Println(planTokens)

	statement, _ := parseStatement(planTokens)

	return statement, nil
}

func parseStatement(tokens []tkn.Token) (plannedstatement, error) {
	var stmt plannedstatement
	for cursor := 0; cursor < len(tokens); cursor++ {
		if tkn.IsKey(tokens[cursor], "plantree") {
			plantree, _ := parseNode(&cursor, tokens)
			stmt.plantree = plantree
		}
	}

	return stmt, nil
}

func parseNode(cursor *int, tokens []tkn.Token) (plannode, error) {
	var node plannode
	var currentLevel int
	for *cursor < len(tokens) {
		*cursor++
		currentToken := tokens[*cursor]

		if currentToken.Token == tkn.ItemId && currentLevel == 0 {
			currentLevel = currentToken.Depth
			fmt.Println("ENTER", currentLevel)
			node.nodetype = currentToken.Value
			continue
		}

		if tkn.IsKey(currentToken, "lefttree") {
			lefttree, _ := parseNode(cursor, tokens)
			node.lefttree = &lefttree
		}

		if tkn.IsKey(currentToken, "righttree") {
			righttree, _ := parseNode(cursor, tokens)
			node.righttree = &righttree
		}

		if currentToken.Token == tkn.ItemEnd && currentToken.Depth == currentLevel {
			fmt.Println("Leave", currentLevel)
			break
		}
	}

	return node, nil
}
