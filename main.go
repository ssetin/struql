package main

import (
	"fmt"

	struql "github.com/ssetin/struql/pkg"
)

// Santa ...
type Santa struct {
	ID     int
	Clause string
}

// Origins ...
type Origins struct {
	Code  string
	Descr string
	Santa []Santa
}

// Device ...
type Device struct {
	Number      int
	Model       string
	Manufacture string
	Oi          []Origins
}

func main() {
	var sq struql.StruQL
	dev := Device{
		Number:      99,
		Model:       "JFQ",
		Manufacture: "Factory",
		Oi: []Origins{
			{
				Code:  "x256",
				Descr: "Debug",
				Santa: []Santa{
					{ID: 1, Clause: "Hoho"},
					{ID: 2, Clause: "Tree"},
					{ID: 3, Clause: "Tree"},
				},
			},
			{
				Code:  "x2599",
				Descr: "Release",
				Santa: []Santa{
					{ID: 1, Clause: "Tooo"},
					{ID: 2, Clause: "Mooo"},
					{ID: 3, Clause: "Laaa"},
				},
			},
		},
	}

	sq.Init(dev)

	sq.Print()

	filter := []struql.Filter{
		{FieldName: "Oi.Descr", Value: "Debug"},
		{FieldName: "Oi.Santa.Clause", Value: "Tree"},
	}
	fmt.Println("Result: ")
	result := make(struql.RowCollection, 0, 5)

	sq.Where(result, filter...).Print()

	fmt.Println("FIN")
}
