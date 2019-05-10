# struql

[![GoDoc](https://godoc.org/github.com/ssetin/struql?status.svg)](https://godoc.org/github.com/ssetin/struql)
[![Build Status](https://travis-ci.org/ssetin/struql.svg?branch=master)](https://travis-ci.org/ssetin/struql)
[![Coverage Status](https://coveralls.io/repos/github/ssetin/struql/badge.svg?branch=master)](https://coveralls.io/github/ssetin/struql?branch=master)
[![Go Report Card](https://goreportcard.com/badge/github.com/ssetin/struql)](https://goreportcard.com/report/github.com/ssetin/struql)

Allows to filter data in the structure, presenting it as a table.
Note, that all filters values should be converted to appropriate field types.

## Install
```
go get github.com/ssetin/struql
```

## Using
Lets consider you have such struct:

```go
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
```
struql represents it as table:

| Number | Model | Details.Id | Details.Value |
|--------|-------|------------|---------------|
| 1      | OEM   | 1          | Fisrt         |
| 1      | OEM   | 2          | Soul          |
| 1      | OEM   | 3          | Seven         |

So now you can filter that data like this:
```go
	sq.Init(dev)
	filter := []struql.Filter{
		{FieldName: "Details.Id", Value: int32(1), Operation: struql.ComparisonGreater},
		{FieldName: "Details.Value", Value: "S", Operation: struql.ComparisonBeginWith},
	}
	dataSet := make(struql.RowCollection, 1)
	dataSet, _ = sq.Where(dataSet, filter...)

	values, _ := dataSet.CollectValues("Details.Id")
	fmt.Printf("Details.Id: %v\n", vals)
```
This way of searching data no so fast as just iterating the struct, but could be useful in some cases.