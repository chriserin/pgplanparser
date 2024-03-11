package main

import (
	"context"
	"fmt"
	"os"
	"slices"

	psr "github.com/chriserin/plan_parser/parser"
	ptr "github.com/chriserin/plan_parser/printer"
	tkn "github.com/chriserin/plan_parser/tokenizer"
	pgx "github.com/jackc/pgx/v5"
)

func main() {
	input := os.Args[1]

	parsedPlan := processPlan(input)

	tables := getTables()

	populateTableNames(&parsedPlan, tables)

	ptr.Print(parsedPlan)
}

func processPlan(planInput string) psr.PlannedStatement {
	plan := []rune(planInput)
	planTokens := tkn.Tokenize(plan)
	parsedPlan, err := psr.ParsePlan(planTokens)

	if err != nil {
		fmt.Println(err)
		return psr.PlannedStatement{}
	}

	return parsedPlan
}

func populateTableNames(parsedPlan *psr.PlannedStatement, tables []postgresTable) {
	setTableName(&parsedPlan.Plantree, parsedPlan.Rtables, tables)
}

func setTableName(node *psr.PlanNode, rtables []psr.Rtable, tables []postgresTable) {
	rtIndex := slices.IndexFunc(rtables, func(rt psr.Rtable) bool { return rt.Rindex == node.Relid })
	if rtIndex >= 0 {
		rtable := rtables[rtIndex]
		ptIndex := slices.IndexFunc(tables, func(pt postgresTable) bool { return pt.relid == rtable.Relid })
		if ptIndex >= 0 {
			pgTable := tables[ptIndex]
			(*node).Tablename = pgTable.relname
		}
	}

	if node.Lefttree != nil {
		setTableName(node.Lefttree, rtables, tables)
	}

	if node.Righttree != nil {
		setTableName(node.Righttree, rtables, tables)
	}
}

type postgresTable struct {
	relid   int
	relname string
}

func getTables() []postgresTable {
	var tables []postgresTable
	urlExample := "postgres://postgres@localhost:5433/postgres_air"
	conn, err := pgx.Connect(context.Background(), urlExample)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	}
	defer conn.Close(context.Background())

	var relid int
	var relname string
	rows, err := conn.Query(context.Background(), "select oid::int, relname from pg_catalog.pg_class where relnamespace = $1", 16389)

	_, err = pgx.ForEachRow(rows, []any{&relid, &relname}, func() error {
		tables = append(tables, postgresTable{relid, relname})
		return nil
	})

	if err != nil {
		fmt.Fprintf(os.Stderr, "QueryRow failed: %v\n", err)
		os.Exit(1)
	}

	return tables
}
