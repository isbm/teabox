package teaboxlib

import (
	"fmt"

	wzlib_logger "github.com/infra-whizz/wzlib/logger"
)

/*
Tabular data definition and parser.
*/

type TeaConfTabularRow struct {
	labels []string
}

func NewTeaConfTabularRow(data []interface{}) *TeaConfTabularRow {
	r := &TeaConfTabularRow{}

	for _, e := range data {
		r.labels = append(r.labels, fmt.Sprintf("%v", e))
	}

	return r
}

// GetLabels returns you the row labels for cell building
func (r *TeaConfTabularRow) GetLabels() []string {
	return r.labels
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

func (tcd *TeaConfTabularData) MakeOptionsData() []*TeaConfCmdOption {
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
	headerRow := &TeaConfTabularRow{labels: tcd.getHeaderLabels(header)}
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
			value:      NewTeaConfTabularRow(row),
		})
	}

	return options
}

// Get header labels
func (tcd *TeaConfTabularData) getHeaderLabels(header []interface{}) []string {
	labels := []string{}

	for _, r := range header {
		labels = append(labels, fmt.Sprintf("%v", r))
	}

	return labels
}
