package main

import (
	"fmt"
	"strings"

	"github.com/ssetin/struql"
)

type santa struct {
	ID     int
	Clause string
}

type origins struct {
	Code  string
	Descr string
	Santa []santa
}

type device struct {
	Number      int
	Model       string
	Manufacture string
	Oi          []origins
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

	dev := device{
		Number:      99,
		Model:       "JFQ",
		Manufacture: "Factory",
		Oi: []origins{
			{
				Code:  "x256",
				Descr: "Debug",
				Santa: []santa{
					{ID: 1, Clause: "Hoho"},
					{ID: 2, Clause: "Tree"},
					{ID: 3, Clause: "Tree"},
				},
			},
			{
				Code:  "x2599",
				Descr: "Release",
				Santa: []santa{
					{ID: 1, Clause: "Tooo"},
					{ID: 2, Clause: "Mooo"},
					{ID: 3, Clause: "Laaa"},
					{ID: 4, Clause: "Qwerty"},
				},
			},
			{
				Code:  "x25990",
				Descr: "Profile",
				Santa: []santa{
					{ID: 1, Clause: "Moose"},
					{ID: 2, Clause: "Fox"},
					{ID: 3, Clause: "Black"},
					{ID: 4, Clause: "Baltic"},
					{ID: 5, Clause: "Pacific"},
					{ID: 6, Clause: "Rocky Mountains"},
					{ID: 7, Clause: "Spotlight"},
					{ID: 8, Clause: "Cortina"},
					{ID: 9, Clause: "TreeTree"},
					{ID: 10, Clause: "Tiger"},
					{ID: 11, Clause: "Tree"},
				},
			},
		},
	}

	sq.Init(dev)

	fmt.Print(sq)

	filter := []struql.Filter{
		{FieldName: "Oi.Descr", Value: modMe("debUg"), Modifier: modMe},
		{FieldName: "Oi.Santa.Clause", Value: "ree", Operation: struql.ComparisonEndWith},
		{FieldName: "Number", Value: 500, Operation: struql.ComparisonLesser, Modifier: modMeInt},
	}
	fmt.Println("Result: ")
	dataSet := make(struql.RowCollection, 0, 5)

	dataSet, err = sq.Where(dataSet, filter...)
	if err != nil {
		fmt.Printf("Error: %s\n", err)
		return
	}

	vals, _ := dataSet.CollectValues("Oi.Santa.ID")
	fmt.Printf("Oi.Santa.ID: %v\n", vals)

	fmt.Println("FIN")
}
