package parser

import (
	"bytes"
	"fmt"
	"strconv"

	tkn "github.com/chriserin/pgplanparser/tokenizer"
)

type PlannedStatement struct {
	Plantree PlanNode
	Rtables  []Rtable
}

type Rtable struct {
	Rindex int
	Relid  int
}

type PlanNode struct {
	Nodetype  string
	Relid     int
	Lefttree  *PlanNode
	Righttree *PlanNode
	Tablename string
}

func (stmt PlannedStatement) String() string {
	return fmt.Sprintf("plantree: --- \n %s \n rtables: --- \n %s", stmt.Plantree, stmt.Rtables)
}

func (node Rtable) String() string {
	return "RTABLE!!"
}

func (node PlanNode) String() string {
	var b bytes.Buffer
	if node.Nodetype != "" {
		b.WriteString("{ " + node.Nodetype + " ")
		b.WriteString(fmt.Sprintf("name: %v ", node.Tablename))
		if node.Lefttree != nil {
			b.WriteString(fmt.Sprintf("%v", node.Lefttree))
		}
		if node.Righttree != nil {
			b.WriteString(fmt.Sprintf("%v", node.Righttree))
		}
		b.WriteString(" }")
		return b.String()
	}
	return ""
}

func ParsePlan(planTokens []tkn.Token) (PlannedStatement, error) {

	statement, _ := parseStatement(planTokens)

	return statement, nil
}

func parseStatement(tokens []tkn.Token) (PlannedStatement, error) {
	var stmt PlannedStatement
	for cursor := 0; cursor < len(tokens); cursor++ {
		if tkn.IsKey(tokens[cursor], "planTree") {
			plantree, _ := parseNode(&cursor, tokens)
			stmt.Plantree = plantree
		}

		if tkn.IsKey(tokens[cursor], "rtable") {
			cursor++
			rtables, _ := parseRtables(&cursor, tokens)
			stmt.Rtables = rtables
		}
	}

	return stmt, nil
}

func parseRtables(cursor *int, tokens []tkn.Token) ([]Rtable, error) {
	var reftables []Rtable

	if tokens[*cursor].Token != tkn.ListStart {
		return reftables, fmt.Errorf("Rtables must be a list")
	}

	for i := 0; *cursor < len(tokens); i++ {
		*cursor++
		if tokens[*cursor].Token == tkn.ItemStart {
			reftable, _ := parseRtable(cursor, tokens, i)
			reftables = append(reftables, reftable)
		}

		if tokens[*cursor].Token == tkn.ListEnd {
			break
		}
	}

	return reftables, nil
}

func parseRtable(cursor *int, tokens []tkn.Token, rtableIndex int) (Rtable, error) {
	var table Rtable
	table.Rindex = rtableIndex + 1
	var currentLevel int
	for *cursor < len(tokens) {
		*cursor++
		currentToken := tokens[*cursor]
		if currentToken.Token == tkn.ItemId && currentLevel == 0 {
			currentLevel = currentToken.Depth
			continue
		}

		if tkn.IsKey(currentToken, "relid") {
			*cursor++
			relid, _ := strconv.Atoi(tokens[*cursor].Value)
			table.Relid = relid
			continue
		}

		if currentToken.Token == tkn.ItemEnd && currentToken.Depth == currentLevel {
			break
		}
	}

	return table, nil
}

func parseNode(cursor *int, tokens []tkn.Token) (PlanNode, error) {
	var node PlanNode
	var currentLevel int
	for *cursor < len(tokens) {
		*cursor++
		currentToken := tokens[*cursor]
		nextToken := tokens[*cursor+1]

		if currentToken.Token == tkn.ItemId && currentLevel == 0 {
			currentLevel = currentToken.Depth
			node.Nodetype = currentToken.Value
			continue
		}

		if tkn.IsKey(currentToken, "lefttree") && nextToken.Token == tkn.ItemStart {
			lefttree, _ := parseNode(cursor, tokens)
			node.Lefttree = &lefttree
		}

		if tkn.IsKey(currentToken, "righttree") && nextToken.Token == tkn.ItemStart {
			righttree, _ := parseNode(cursor, tokens)
			node.Righttree = &righttree
		}

		if tkn.IsKey(currentToken, "relid") {
			*cursor++
			value := tokens[*cursor].Value
			relid, _ := strconv.Atoi(value)
			node.Relid = relid
		}

		if currentToken.Token == tkn.ItemEnd && currentToken.Depth == currentLevel {
			break
		}
	}

	return node, nil
}
