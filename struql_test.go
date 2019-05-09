package struql

import (
	"errors"
	"fmt"
	"reflect"
	"strings"
	"testing"
)

type santa struct {
	ID     int
	Clause string
	Chain  []int
}

type custom struct {
	MaskI64 int64
	MaskF64 float64
}

type origins struct {
	ID     int32
	Code   string
	Descr  string
	Santa  []*santa
	Custom custom
}

type device struct {
	Number      int
	Model       string
	Version     float32
	Manufacture string
	Oi          []origins
	SomeByte    byte
}

var (
	dev = device{
		Number:      99,
		Model:       "JFQ",
		Version:     3.78,
		Manufacture: "Factory",
		SomeByte:    1,
		Oi: []origins{
			{
				ID:    190,
				Code:  "x256",
				Descr: "Debug",
				Santa: []*santa{
					{ID: 1, Clause: "Hoho", Chain: []int{8, 9, 21, 24, 3, 48, 69}},
					{ID: 2, Clause: "Tree", Chain: []int{1, 2, 3, 4, 5, 8, 9}},
					{ID: 3, Clause: "Tree", Chain: []int{31, 32, 33, 34, 35, 38, 39}},
				},
				Custom: custom{
					MaskF64: 1321.34,
					MaskI64: 800,
				},
			},
			{
				ID:    191,
				Code:  "x2599",
				Descr: "Release",
				Santa: []*santa{
					{ID: 11, Clause: "Tooo", Chain: []int{91, 12, 23, 43, 5, 38, 93}},
					{ID: 12, Clause: "Mooo", Chain: []int{11, 22}},
					{ID: 13, Clause: "Laaa", Chain: []int{1, 2, 39, 4, 5, 8, 9}},
					{ID: 14, Clause: "Qwerty", Chain: []int{1, 55, 86, 97, 39}},
				},
				Custom: custom{
					MaskF64: 2321.34,
					MaskI64: 12100,
				},
			},
			{
				ID:    200,
				Code:  "x25990",
				Descr: "Profile",
				Santa: []*santa{
					{ID: 31, Clause: "Moose"},
					{ID: 32, Clause: "Fox"},
					{ID: 33, Clause: "Black"},
					{ID: 34, Clause: "Baltic"},
					{ID: 35, Clause: "Pacific"},
					{ID: 36, Clause: "Rocky Mountains"},
					{ID: 37, Clause: "Spotlight"},
					{ID: 38, Clause: "Cortina"},
					{ID: 39, Clause: "TreeTree"},
					{ID: 40, Clause: "Tiger"},
					{ID: 41, Clause: "Tree"},
				},
			},
		},
	}

	sq StruQL
)

func modMe(s interface{}) interface{} {
	return strings.ToLower(s.(string))
}

func init() {
	sq.Init(dev)
}

type fieldValues struct {
	fieldName string
	values    interface{}
}

type CR struct {
	caseName string
	filter   []Filter
	err      error
	values   fieldValues
}

func TestPrinting(t *testing.T) {
	var isq StruQL
	dev := device{
		Number: 100,
		Model:  "Chip",
	}

	isq.Init(dev)
	fmt.Print(isq)

	filter := []Filter{
		{FieldName: "Model", Value: "Chip"},
	}
	dataSet := make(RowCollection, 0, 1)
	dataSet, _ = isq.Where(dataSet, filter...)

	fmt.Print(dataSet)
}

func TestInitializing(t *testing.T) {
	var isq StruQL

	err := isq.Init(29)
	if err == nil {
		t.Error("Should be initialized only by struct value")
	}
}

func TestGeneral(t *testing.T) {
	var err error
	// Unexisting field and field index
	filter := []Filter{
		{FieldName: "Oi.Descr", Value: "Debug"},
		{FieldName: "Oi.Santa.Clause", Value: "Tree"},
	}

	dataSet := make(RowCollection, 0, 2)
	dataSet, err = sq.Where(dataSet, filter...)
	if err != nil {
		t.Error(err)
	}

	fld := dataSet[0].FieldByIndex(999)
	if fld != nil {
		t.Error("Field should be nil when trying to access it by wrong index")
	}

	_, err = dataSet.CollectValues("unexpected field name")
	if err == nil {
		t.Error("Should be error on collecting values from unexisting field")
	}

	// Empty dataset
	filter = []Filter{
		{FieldName: "Oi.Descr", Value: "QWERTY"},
	}
	dataSet = make(RowCollection, 0, 2)
	dataSet, err = sq.Where(dataSet, filter...)
	if err != nil {
		t.Error(err)
	}
	_, err = dataSet.CollectValues("Oi.Santa.ID")
	if err == nil {
		t.Error("Should be error on collecting values from empty dataset")
	}

	// Try to filter empty dataset
	dataSet, err = dataSet.Where(dataSet, filter...)
	if err == nil {
		t.Error("Should be error when filtering empty dataset")
	}

	// Unexisting fieldname in filter
	filter = []Filter{
		{FieldName: "Oi.Descrucio", Value: "Debug"},
		{FieldName: "Oi.Santa.Clause", Value: "Tree"},
	}

	dataSet = make(RowCollection, 0, 2)
	dataSet, err = sq.Where(dataSet, filter...)
	if err == nil {
		t.Error("Should be error on quering with unexisting fieldname in filter")
	}
}

func TestQuering(t *testing.T) {

	cases := []CR{
		{
			caseName: "1. Simple filter",
			filter: []Filter{
				{FieldName: "Oi.Descr", Value: "Debug"},
				{FieldName: "Oi.Santa.ID", Value: 1, Operation: ComparisonGreater},
				{FieldName: "Oi.Santa.Clause", Value: "Tree"},
			},
			values: fieldValues{
				fieldName: "Oi.Santa.ID",
				values:    []interface{}{2, 3},
			},
		},
		{
			caseName: "2. String begin with + modifier",
			filter: []Filter{
				{FieldName: "Oi.Code", Value: "x25990"},
				{FieldName: "Oi.Santa.Clause", Value: "tree", Modifier: modMe, Operation: ComparisonBeginWith},
			},
			values: fieldValues{
				fieldName: "Oi.Santa.ID",
				values:    []interface{}{39, 41},
			},
		},
		{
			caseName: "3. Int32, lesser, exists in slice",
			filter: []Filter{
				{FieldName: "Oi.ID", Value: int32(200), Operation: ComparisonLesser},
				{FieldName: "Version", Value: float32(3.78), Operation: ComparisonEqual},
				{FieldName: "Oi.Santa.Chain", Value: 39, Operation: ComparisonIn},
			},
			values: fieldValues{
				fieldName: "Oi.Santa.ID",
				values:    []interface{}{3, 13, 14},
			},
		},
		{
			caseName: "4. String end with, int32 greater, string lesser, string end with, float32 greater",
			filter: []Filter{
				{FieldName: "Oi.ID", Value: int32(100), Operation: ComparisonGreater},
				{FieldName: "Version", Value: float32(1), Operation: ComparisonGreater},
				{FieldName: "Oi.Santa.Clause", Value: "ic", Operation: ComparisonEndWith},
				{FieldName: "Oi.Descr", Value: "ZZZZZZZZ", Operation: ComparisonLesser},
				{FieldName: "Oi.Descr", Value: "Ab", Operation: ComparisonGreater},
			},
			values: fieldValues{
				fieldName: "Oi.Santa.ID",
				values:    []interface{}{34, 35},
			},
		},
		{
			caseName: "5. Unsupported Comparison 1",
			filter: []Filter{
				{FieldName: "Number", Value: 1, Operation: ComparisonEndWith},
			},
			err: errors.New("unsupported comparison"),
		},
		{
			caseName: "6. Unsupported Comparison 2",
			filter: []Filter{
				{FieldName: "Number", Value: 1, Operation: ComparisonBeginWith},
			},
			err: errors.New("unsupported comparison"),
		},
		{
			caseName: "7. Unsupported Comparison 3",
			filter: []Filter{
				{FieldName: "SomeByte", Value: 1, Operation: ComparisonGreater},
			},
			err: errors.New("unsupported comparison"),
		},
		{
			caseName: "8. Unsupported Comparison 4",
			filter: []Filter{
				{FieldName: "SomeByte", Value: 1, Operation: ComparisonLesser},
			},
			err: errors.New("unsupported comparison"),
		},
		{
			caseName: "9. Unsupported Comparison 5",
			filter: []Filter{
				{FieldName: "Number", Value: 1, Operation: ComparisonIn},
			},
			err: errors.New("unsupported comparison"),
		},
		{
			caseName: "10. Unsupported Comparison 6",
			filter: []Filter{
				{FieldName: "Number", Value: 1, Operation: 99},
			},
			err: errors.New("unsupported comparison"),
		},
		{
			caseName: "11. Another filters",
			filter: []Filter{
				{FieldName: "Oi.Santa.ID", Value: 20, Operation: ComparisonLesser},
				{FieldName: "Version", Value: float32(9), Operation: ComparisonLesser},
				{FieldName: "Oi.Custom.MaskF64", Value: float64(2012.33), Operation: ComparisonLesser},
				{FieldName: "Oi.Custom.MaskF64", Value: float64(12.3), Operation: ComparisonGreater},
				{FieldName: "Oi.Custom.MaskI64", Value: int64(-454), Operation: ComparisonGreater},
				{FieldName: "Oi.Custom.MaskI64", Value: int64(900), Operation: ComparisonLesser},
			},
			values: fieldValues{
				fieldName: "Oi.Santa.ID",
				values:    []interface{}{1, 2, 3},
			},
		},
	}

	RunQueringTests(t, cases)
}

func RunQueringTests(t *testing.T, cases []CR) {
	var err error

	for _, testCase := range cases {
		t.Run(testCase.caseName, func(t *testing.T) {
			dataSet := make(RowCollection, 0, 2)
			dataSet, err = sq.Where(dataSet, testCase.filter...)
			if !compareError(err, testCase.err) {
				t.Errorf("Error check:\nWant: %v\n Got: %v\n", testCase.err, err)
				return
			}
			if err != nil {
				return
			}

			collectedVals, err := dataSet.CollectValues(testCase.values.fieldName)
			if err != nil {
				t.Error(err)
				return
			}

			if !reflect.DeepEqual(collectedVals, testCase.values.values) {
				t.Errorf("\nWant: %v\n Got: %v\n", testCase.values.values, collectedVals)
			}
		})
	}
}

func compareError(a, b error) bool {
	return (a == nil && b == nil) || (a != nil && b != nil && a.Error() == b.Error())
}

/*
	go test -v -cover
	go test -coverprofile=cover.out
	go tool cover -html=cover.out -o cover.html
*/
