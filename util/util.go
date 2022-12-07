package util

import (
	"fmt"
	"io"

	"github.com/dstotijn/go-notion"
)

type Row struct {
	Cells map[string]string
}

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
	fmt.Fprintf(w, "%s\t%s\t%s\n", "page ID", "TITLE", "TAGS")
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
				// if header{//headerをkeyに設定する}
				// cells = append(cells, title)
				cells[k] = title
			case notion.DBPropTypeMultiSelect:
				multiSelect = fmt.Sprintf("%s", v.MultiSelect)
				cells[k] = multiSelect
				// if wide{}
			}
		}
		row.Cells = cells
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
	// header := true
	for _, page := range pages {
		props := page.Properties
		// fmt.Printf("%#v", props.(notion.DatabasePageProperties))
		var title string
		var multiSelect string
		row := Row{}
		cells := map[string]string{}
		cells["pageID"] = page.ID
		for k, v := range props.(notion.DatabasePageProperties) {
			// Rowを作っていく
			switch v.Type {
			case notion.DBPropTypeTitle:
				title = fmt.Sprintf("%s", v.Title[0].Text.Content) //pythonにおけるmap的な書き方できないんだろうか?
				// if header{//headerをkeyに設定する}
				// TODO: property typeとkをcellのmapのkeyとして一緒に格納する
				cells[k] = title
			case notion.DBPropTypeMultiSelect:
				multiSelect = fmt.Sprintf("%s", v.MultiSelect)
				cells[k] = multiSelect
				// if wide{}
			}
		}
		row.Cells = cells
		// fmt.Fprintf(w, "%s\t%s\t%s\n", page.ID, title, multiSelect)
		rows = append(rows, row)
	}
	table.Rows = rows
	// fmt.Println(table)

	return table
}

var printOrder = []string{}

// cellの文字数を決めてprintしたい。
func (t *TablePrinter) Print() {
	for _, r := range t.table.Rows {
		var output string
		// printする順番をどこかで定義しないといけない
		for _, c := range r.Cells {
			// if i == len(r.Cells)-1 {
			// 	output += fmt.Sprintf("%s\n", c)
			// 	continue
			// }
			output += fmt.Sprintf("%s\t", c)
		}
		// output += "\n"
		fmt.Fprint(t.writer, output)
	}
}
