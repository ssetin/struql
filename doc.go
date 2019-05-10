/*
Package struql helps to filter data in the structure, presenting it as a table.

Typical usage

	type detail struct {
		Id    int32
		Value string
	}

	type device struct {
		Number  int
		Model   string
		Details []detail
	}

	var dev = device{
		Number: 1,
		Model:  "OEM",
		Details: []detail{
			{Id: 1, Value: "First"},
			{Id: 2, Value: "Soul"},
			{Id: 3, Value: "Seven"},
		},
	}

	var sq struql.StruQL
	sq.Init(dev)
	filter := []struql.Filter{
		{FieldName: "Details.Id", Value: int32(1), Operation: struql.ComparisonGreater},
		{FieldName: "Details.Value", Value: "S", Operation: struql.ComparisonBeginWith},
	}
	dataSet := make(struql.RowCollection, 1)
	dataSet, _ = sq.Where(dataSet, filter...)

	values, _ := dataSet.CollectValues("Details.Id")
	fmt.Printf("Details.Id: %v\n", vals)

*/
package struql
