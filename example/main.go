package main

import (
	"fmt"
	"strings"

	"github.com/ssetin/struql"
)

// Santa ...
type Santa struct {
	ID     int
	Clause string
}

// Origins ...
type Origins struct {
	Code       string
	Descr      string
	Santa      []*Santa
	Collection []int32
}

// Device ...
type Device struct {
	Number      int
	Model       string
	Manufacture string
	Oi          []Origins
}

func modMe(s interface{}) interface{} {
	return strings.ToLower(s.(string))
}

func modMeInt(s interface{}) interface{} {
	return s.(int) + 1
}

func main() {
	var (
		sq  struql.StruQL
		err error
	)

	dev := Device{
		Number:      99,
		Model:       "JFQ",
		Manufacture: "Factory",
		Oi: []Origins{
			{
				Code:       "x256",
				Descr:      "Debug",
				Collection: []int32{1, 6, 99, 100, 11},
				Santa: []*Santa{
					{ID: 1, Clause: "Hoho"},
					{ID: 2, Clause: "Tree"},
					{ID: 3, Clause: "Tree"},
				},
			},
			{
				Code:       "x2599",
				Descr:      "Release",
				Collection: []int32{101, 102, 103},
				Santa: []*Santa{
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
		{FieldName: "Oi.Descr", Value: modMe("debUg"), Modifier: modMe},
		{FieldName: "Oi.Santa.Clause", Value: "ree", Operation: struql.ComparsionEndWith},
		{FieldName: "Oi.Collection", Value: int32(100), Operation: struql.ComparsionIn},
		{FieldName: "Number", Value: 500, Operation: struql.ComparsionLesser, Modifier: modMeInt},
	}
	fmt.Println("Result: ")
	dataSet := make(struql.RowCollection, 0, 5)

	dataSet, err = sq.Where(dataSet, filter...)
	if err != nil {
		fmt.Printf("Error: %s\n", err)
		return
	}

	dataSet.Print()

	fmt.Println("FIN")
}
