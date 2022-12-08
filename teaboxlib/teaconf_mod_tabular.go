package teaboxlib

import (
	"fmt"
	"strings"

	wzlib_logger "github.com/infra-whizz/wzlib/logger"
)

/*
Tabular data definition and parser.
*/

type TeaConfTabularRow struct {
	valueHidden bool
	labels      []string
	value       interface{}
}

func NewTeaConfTabularRow(data []interface{}, valueIdx int, valueHidden bool) *TeaConfTabularRow {
	r := &TeaConfTabularRow{
		value:       data[valueIdx],
		valueHidden: valueHidden,
	}

	for _, e := range data {
		r.labels = append(r.labels, fmt.Sprintf("%v", e))
	}

	return r
}

// GetLabels returns you the row labels for cell building
func (r *TeaConfTabularRow) GetLabels() []string {
	return r.labels
}

// GetValue returns you the value of the row. If this is a header,
// this returns you an index of value-carrier column.
func (r *TeaConfTabularRow) GetValue() interface{} {
	return r.value
}

// IsValueHidden returns a bool, which indicates if the value column is hidden.
func (r *TeaConfTabularRow) IsValueHidden() bool {
	return r.valueHidden
}

type TeaConfTabularData struct {
	header []string
	rows   [][]*TeaConfTabularRow
	data   []interface{}

	wzlib_logger.WzLogger
}

func NewTeaConfTabularData(data []interface{}) *TeaConfTabularData {
	return &TeaConfTabularData{
		header: []string{},
		rows:   [][]*TeaConfTabularRow{},
		data:   data,
	}
}

func (tcd *TeaConfTabularData) Make() []*TeaConfCmdOption {
	options := []*TeaConfCmdOption{}

	// Bail-out if there is no data
	if tcd.data == nil {
		return options
	}

	// Get table header
	header, ok := tcd.data[0].([]interface{})
	if !ok {
		panic("Wrong type of tabular header") // XXX: remove this
		//return options
	}

	// Find which field is a response value and setup its attrs
	labels, attributes, index := tcd.getHeaderAttributes(header)
	headerRow := &TeaConfTabularRow{
		labels: labels,
		value:  -1,
	}

	for _, attr := range attributes {
		switch attr {
		case "hidden":
			headerRow.valueHidden = true
		case "value":
			headerRow.value = index
		}
	}

	options = append(options, &TeaConfCmdOption{
		label:      "",
		optionType: "tabular:header",
		value:      headerRow,
	})

	// Get tabular body
	for _, rdata := range tcd.data[1:] {
		row, ok := rdata.([]interface{})
		if !ok {
			panic("Wrong tabular data: should be an array")
		}

		options = append(options, &TeaConfCmdOption{
			label:      "",
			optionType: "tabular:row",
			value:      NewTeaConfTabularRow(row, headerRow.value.(int), headerRow.valueHidden),
		})
	}

	return options
}

// Get attributes for the header: return clean labels, and attributes
// Field attributes are:
//
//   - value
//   - hidden
//
// They are written in any order, like so:
//
//	value:hidden:"Whatever label it is"
func (tcd *TeaConfTabularData) getHeaderAttributes(header []interface{}) ([]string, []string, int) {
	labels := []string{}
	attrs := []string{}
	idx := -1

	for i, r := range header {
		label := r.(string)
		if strings.Contains(label, ":") && idx < 0 {
			attrs = strings.SplitN(label, ":", 3)
			label = attrs[(len(attrs) - 1)]
			attrs = attrs[:(len(attrs) - 1)]
			idx = i
		}
		labels = append(labels, label)
	}

	return labels, attrs, idx
}
