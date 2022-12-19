package util

import (
	"fmt"
	"io"
	"sort"
	"strings"
	"time"

	"github.com/dstotijn/go-notion"
)

var printOrderOfDBPropType = []notion.DatabasePropertyType{
	notion.DBPropTypeTitle,
	notion.DBPropTypeSelect,
	notion.DBPropTypeMultiSelect,
	notion.DBPropTypeCreatedTime,
	notion.DBPropTypeRichText,
}

func getIdxOfPrintOrder(query notion.DatabasePropertyType) int {
	for i, s := range printOrderOfDBPropType {
		if query == s {
			return i
		}
	}
	return -1
}

type propTypeValue struct {
	propType notion.DatabasePropertyType
	propName string
	value    string
}

type Row struct {
	Cells []propTypeValue
}

type Table struct {
	Rows []Row
}

type TablePrinter struct {
	table  Table
	writer io.Writer
}

func NewDatabaseQueryResponcePrinter(res notion.DatabaseQueryResponse, w io.Writer) *TablePrinter {
	return &TablePrinter{table: newTable(res.Results), writer: w}

}

func NewCreatePageResponcePrinter(res notion.Page, w io.Writer) *TablePrinter {
	return &TablePrinter{table: newTable([]notion.Page{res}), writer: w}

}

func richTextToString(r []notion.RichText) string {
	if len(r) == 0 {
		return " "
	} else {
		return r[0].Text.Content
	}
}

func selectMultiOptionsToStrig(so []notion.SelectOptions) string {
	result := ""
	for _, v := range so {
		result += v.Name + ","
	}
	return result
}

func selectOptionToString(so *notion.SelectOptions) string {
	if so == nil {
		return " "
	}
	return so.Name

}

func newTable(pages []notion.Page) Table {
	table := Table{}
	rows := []Row{}
	// header := []string{}

	for _, page := range pages {
		props := page.Properties
		row := Row{}
		cells := []propTypeValue{}
		for k, v := range props.(notion.DatabasePageProperties) {
			// Rowを作っていく
			switch v.Type {
			case notion.DBPropTypeTitle:
				cells = append(cells, propTypeValue{
					propType: notion.DBPropTypeTitle,
					propName: k,
					value:    fmt.Sprintf("%s", richTextToString(v.Title)),
				})

			case notion.DBPropTypeMultiSelect:
				cells = append(cells, propTypeValue{
					propType: notion.DBPropTypeMultiSelect,
					propName: k,
					value:    selectMultiOptionsToStrig(v.MultiSelect),
				})

			case notion.DBPropTypeSelect:
				cells = append(cells, propTypeValue{
					propType: notion.DBPropTypeSelect,
					propName: k,
					value:    selectOptionToString(v.Select),
				})
			case notion.DBPropTypeCreatedTime:
				cells = append(cells, propTypeValue{
					propType: notion.DBPropTypeCreatedTime,
					propName: "AGE",
					value:    fmt.Sprintf("%v", HumanDuration(time.Now().Sub(*v.CreatedTime))), // calcurate age
				})
			case notion.DBPropTypeRichText:
				cells = append(cells, propTypeValue{
					propType: notion.DBPropTypeRichText,
					propName: k,
					value:    richTextToString(v.RichText),
				})
			}
		}
		row.Cells = cells
		rows = append(rows, row)
	}
	table.Rows = rows

	return table
}

func (t *TablePrinter) Print() {

	// TODO: 全角を考慮したcellの幅を計算し、paddingする。

	for i, r := range t.table.Rows {
		var output string
		sort.Slice(r.Cells, func(i, j int) bool {
			return getIdxOfPrintOrder(r.Cells[i].propType) < getIdxOfPrintOrder(r.Cells[j].propType)
		})

		// write header
		if i == 0 {
			for _, c := range r.Cells {
				// TODO: string builderを使う！
				output += fmt.Sprintf("%s\t", c.propName)
			}
			output += "\n"
			fmt.Fprint(t.writer, strings.ToUpper(output))
			output = ""
		}

		for _, c := range r.Cells {
			// TODO: string builderを使う！
			if c.propType == notion.DBPropTypeRichText {
				output += fmt.Sprintf("%.24s\t", c.value)
				continue
			}
			output += fmt.Sprintf("%s\t", c.value)
		}
		output += "\n"
		fmt.Fprint(t.writer, output)
	}
}
