package printer

import (
	"bytes"
	"fmt"
	"strconv"
	"strings"

	psr "github.com/chriserin/pgplanparser/parser"
)

type line struct {
	depth   int
	content string
}

func Print(stmt psr.PlannedStatement) {
	lines := []line{}
	depth := 0
	getLines(&lines, &stmt.Plantree, depth)
	output := printPlan(lines)
	fmt.Println(output)
}

func getLines(lines *[]line, node *psr.PlanNode, depth int) {
	depth++
	var b bytes.Buffer
	b.WriteString(strategyStr(node.Nodetype, node.Strategy))
	b.WriteString(node.Nodetype)
	b.WriteString(cmdStr(node.Nodetype, node.Cmd))
	b.WriteString(" ")
	b.WriteString(node.Tablename)
	*lines = append(*lines, line{depth, b.String()})
	if node.Lefttree != nil {
		getLines(lines, node.Lefttree, depth)
	}
	if node.Righttree != nil {
		getLines(lines, node.Righttree, depth)
	}
}

func cmdStr(nodetype string, cmd int) string {
	if nodetype != "SETOP" {
		return ""
	}
	switch cmd {
	case 0:
		return "Intersect"
	case 2:
		return "Except"
	}

	return ""
}

func strategyStr(nodetype string, strategy int) string {
	if nodetype != "SETOP" {
		return ""
	}

	switch strategy {
	case 1:
		return "Hash"
	}

	return ""
}

func printPlan(lines []line) string {
	hSize := getMaxLength(lines) + 20

	top := "┌" + strings.Repeat("─", hSize) + "┐\n"
	title := "│    Query Plan " + strings.Repeat(" ", hSize-15) + "│\n"
	sep := "├" + strings.Repeat("─", hSize) + "┤\n"
	bottom := "└" + strings.Repeat("─", hSize) + "┘\n"

	output := top + title + sep

	for i := 0; i < len(lines); i++ {
		output += ("│" + strconv.Itoa(lines[i].depth) + " " + lines[i].content + strings.Repeat(" ", hSize-len(lines[i].content)-2) + "│" + "\n")
	}

	return output + bottom
}

func getMaxLength(lines []line) int {

	max := 0

	for i := 0; i < len(lines); i++ {
		max = Max(max, len(lines[i].content))
	}

	return max
}

func Max(x int, y int) int {
	if x > y {
		return x
	}
	return y
}
