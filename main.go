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
	//IDs   []int
	Santa []Santa
}

// Device ...
type Device struct {
	Number int
	Model  string
	Man    string
	Oi     Origins
}

func main() {
	var sq struql.StruQL
	dev := Device{
		Number: 99,
		Model:  "JFQ",
		Man:    "Sony",
		Oi: Origins{
			Code:  "x256",
			Descr: "Debug",
			Santa: []Santa{
				Santa{ID: 1, Clause: "Hoho"},
				Santa{ID: 2, Clause: "Rum"},
				Santa{ID: 3, Clause: "Tree"},
			},
		}}

	sq.Init(dev)

	sq.Print()

	//fmt.Print(sq)
	filter := struql.Filter{FieldName: "Oi.Santa.Clause", Value: "Rum"}
	fmt.Println("Result: ")
	sq.Where(filter).Print()

	fmt.Println("FIN")
}
