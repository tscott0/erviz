package main

import (
	"encoding/csv"
	"fmt"
	"os"
)

type ForeignKey struct {
	Table string
	Index string
	Col1  string
	Col2  string
	Col3  string

	ForeignTable string
	ForeignIndex string
	ForeignCol1  string
	ForeignCol2  string
	ForeignCol3  string
}

type Column struct {
	Column string
	Type   string
}

// go run main.go | dot -Tpng -o test.png && feh test.png
func main() {
	// Foreign Keys
	csvFile, err := os.Open("./foreign_keys")

	if err != nil {
		fmt.Println(err)
	}

	defer csvFile.Close()

	reader := csv.NewReader(csvFile)

	reader.Comma = '|'

	reader.FieldsPerRecord = -1

	csvData, err := reader.ReadAll()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	var fk ForeignKey
	var allForeignKeys []ForeignKey

	for _, each := range csvData {
		fk.Table = each[0]
		fk.Index = each[1]
		fk.Col1 = each[2]
		fk.Col2 = each[3]
		fk.Col3 = each[4]

		fk.ForeignTable = each[5]
		fk.ForeignIndex = each[6]
		fk.ForeignCol1 = each[7]
		fk.ForeignCol2 = each[8]
		fk.ForeignCol3 = each[9]

		allForeignKeys = append(allForeignKeys, fk)
	}

	// Columns
	csvFile, err = os.Open("./columns")

	if err != nil {
		fmt.Println(err)
	}

	defer csvFile.Close()

	reader = csv.NewReader(csvFile)

	reader.Comma = '|'

	reader.FieldsPerRecord = -1

	csvData, err = reader.ReadAll()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	var c Column
	tables := make(map[string][]Column)

	for _, each := range csvData {
		table := each[0]

		c.Column = each[1]
		c.Type = each[2]

		cols, ok := tables[table]
		if !ok {
			cols = make([]Column, 0)
		}

		cols = append(cols, c)

		tables[table] = cols
	}

	printDot(allForeignKeys, tables)
}

func printDot(foreign_keys []ForeignKey, tables map[string][]Column) {

	fmt.Println(`digraph structs {
   rankdir=LR;
   node [shape=record,style=filled,fillcolor=".7 .3 1.0"];`)

	// Tables
	for table, columns := range tables {
		line := fmt.Sprintf("   %v [shape=record,label=\"%v\\n\\n", table, table)
		for _, c := range columns {
			line += fmt.Sprintf("|<%v>%v %v\\l", c.Column, c.Column, c.Type)
		}
		line += "\"];"
		fmt.Println(line)
	}

	//person:f1 -> score:f0;
	//person:f2 -> game:here;
	for _, fk := range foreign_keys {
		fmt.Printf("   %v:%v -> %v:%v [label=\"%v\"];\n",
			fk.Table, fk.Col1, fk.ForeignTable, fk.ForeignCol1, fk.Index)
	}

	fmt.Println("}")
}
