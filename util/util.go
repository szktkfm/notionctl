package util

import (
	"fmt"
	"io"
	"sort"

	"github.com/dstotijn/go-notion"
)

func getPos(query notion.DatabasePropertyType) int {
	var printOrder = []notion.DatabasePropertyType{
		notion.DBPropTypeTitle,
		notion.DBPropTypeMultiSelect,
		notion.DBPropTypeDate,
	}
	for i, s := range printOrder {
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

// type Row struct {
// 	Cells map[string]string
// }

//	type TableHeader struct {
//		Header []string
//	}
type Table struct {
	Rows []Row
}

type TablePrinter struct {
	table  Table
	writer io.Writer
}

// tablePrinterを生成
func PrintDatabaseQueryResponce(res notion.DatabaseQueryResponse, w io.Writer) error {
	fmt.Fprintf(w, "%s\t%s\t%s\n", "page ID", "TITLE", "TAGS")
	PrintPage(res.Results, w)

	return nil
}

func NewDatabaseQueryResponcePrinter(res notion.DatabaseQueryResponse, w io.Writer) *TablePrinter {
	//fmt.Fprintf(w, "%s\t%s\t%s\n", "page ID", "TITLE", "TAGS")
	// fmt.Fprintf(w, "%s\t%s\n", "TITLE", "TAGS")
	return &TablePrinter{table: newTable(res.Results), writer: w}

}

// build table method
// 具体的には、1つのfor loopでRowを生成する
func PrintPage(pages []notion.Page, w io.Writer) error {
	table := &Table{}
	rows := []Row{}
	// header := true
	for _, page := range pages {
		props := page.Properties
		// fmt.Printf("%#v", props.(notion.DatabasePageProperties))
		var title string
		var multiSelect string
		row := Row{}
		cells := map[string]string{}
		for k, v := range props.(notion.DatabasePageProperties) {
			// Rowを作っていく
			switch v.Type {
			case notion.DBPropTypeTitle:
				title = fmt.Sprintf("%s", v.Title[0].Text.Content) //pythonにおけるmap的な書き方できないんだろうか?
				// title = fmt.Sprintf("%s", v.Value()) //pythonにおけるmap的な書き方できないんだろうか?
				// if header{//headerをkeyに設定する}
				// cells = append(cells, title)
				cells[k] = title
			case notion.DBPropTypeMultiSelect:
				multiSelect = fmt.Sprintf("%s", v.MultiSelect)
				cells[k] = multiSelect
				// if wide{}
			}
		}
		// row.Cells = cells
		fmt.Fprintf(w, "%s\t%s\t%s\n", page.ID, title, multiSelect)
		rows = append(rows, row)
	}
	table.Rows = rows
	// fmt.Println(table)

	return nil
}

func newTable(pages []notion.Page) Table {
	table := Table{}
	rows := []Row{}
	// header := []string{}

	for _, page := range pages {
		props := page.Properties
		// fmt.Printf("%#v", props.(notion.DatabasePageProperties))
		var title string
		var multiSelect string
		row := Row{}
		// cells := map[propTypeValue]string{}
		cells := []propTypeValue{}
		// cells["pageID"] = page.ID
		for k, v := range props.(notion.DatabasePageProperties) {
			// Rowを作っていく
			switch v.Type {
			case notion.DBPropTypeTitle:
				// TODO: Untitledだとindex out of rangeになる
				title = fmt.Sprintf("%s", v.Title[0].Text.Content) //pythonにおけるmap的な書き方できないんだろうか?
				// title = fmt.Sprintf("%s", v.Value()) //pythonにおけるmap的な書き方できないんだろうか?
				cells = append(cells, propTypeValue{
					propType: notion.DBPropTypeTitle,
					propName: k,
					value:    title,
				})

			case notion.DBPropTypeMultiSelect:
				multiSelect = fmt.Sprintf("%s", v.MultiSelect)
				cells = append(cells, propTypeValue{
					propType: notion.DBPropTypeMultiSelect,
					propName: k,
					value:    multiSelect,
				})
			}
		}
		row.Cells = cells
		rows = append(rows, row)
	}
	table.Rows = rows
	// fmt.Println(table)

	return table
}

// cellの文字数を決めてprintしたい。
func (t *TablePrinter) Print() {

	for i, r := range t.table.Rows {
		var output string
		// printする順番をどこかで定義しないといけない
		sort.Slice(r.Cells, func(i, j int) bool { return getPos(r.Cells[i].propType) < getPos(r.Cells[j].propType) })

		// write header
		// TODO: propNameがからのときの対応
		if i == 0 {
			for _, c := range r.Cells {
				// TODO: string builderを使う！
				output += fmt.Sprintf("%s\t", c.propName)
			}
			output += "\n"
			fmt.Fprint(t.writer, output)
			output = ""
		}

		for _, c := range r.Cells {
			// TODO: string builderを使う！
			output += fmt.Sprintf("%s\t", c.value)
		}
		output += "\n"
		fmt.Fprint(t.writer, output)
	}
}
