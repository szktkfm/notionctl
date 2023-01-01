package util

import (
	"bytes"
	"testing"
	"time"

	"github.com/dstotijn/go-notion"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
)

func mustParseTime(layout, value string) time.Time {
	t, err := time.Parse(layout, value)
	if err != nil {
		panic(err)
	}
	return t
}

func mustParseDateTime(value string) notion.DateTime {
	dt, err := notion.ParseDateTime(value)
	if err != nil {
		panic(err)
	}
	return dt
}

func TestNewTable(t *testing.T) {

	tests := []struct {
		name string
		arg  []notion.Page
		want Table
	}{
		{
			name: "new table test",
			arg: []notion.Page{
				{
					ID: "7c6b1c95-de50-45ca-94e6-af1d9fd295ab",
					// CreatedTime:    mustParseTime(time.RFC3339Nano, "2021-05-18T17:50:22.371Z"),
					// LastEditedTime: mustParseTime(time.RFC3339Nano, "2021-05-18T17:50:22.371Z"),
					URL: "https://www.notion.so/Avocado-251d2b5f268c4de2afe9c71ff92ca95c",
					Parent: notion.Parent{
						Type:       notion.ParentTypeDatabase,
						DatabaseID: "39ddfc9d-33c9-404c-89cf-79f01c42dd0c",
					},
					Archived: false,
					Properties: notion.DatabasePageProperties{
						"Date": notion.DatabasePageProperty{
							ID:   "Q]uT",
							Type: notion.DBPropTypeDate,
							Name: "Date",
							Date: &notion.Date{
								Start: mustParseDateTime("2021-05-18T12:49:00.000-05:00"),
							},
						},
						"Name": notion.DatabasePageProperty{
							ID:   "title",
							Type: notion.DBPropTypeTitle,
							Name: "Name",
							Title: []notion.RichText{
								{
									Type: notion.RichTextTypeText,
									Text: &notion.Text{
										Content: "Foobar",
									},
									PlainText: "Foobar",
									Annotations: &notion.Annotations{
										Color: notion.ColorDefault,
									},
								},
							},
						},
						"Description": notion.DatabasePageProperty{
							ID:   "ja0H",
							Type: notion.DBPropTypeRichText,
							Name: "Name",
							RichText: []notion.RichText{
								{
									Type: notion.RichTextTypeText,
									Text: &notion.Text{
										Content: "Foobar",
									},
									PlainText: "Foobar",
									Annotations: &notion.Annotations{
										Color: notion.ColorDefault,
									},
								},
							},
						},
						"Food group": notion.DatabasePageProperty{
							ID:   "TJmr",
							Type: notion.DBPropTypeSelect,
							Select: &notion.SelectOptions{
								ID:    "96eb622f-4b88-4283-919d-ece2fbed3841",
								Name:  "🥦Vegetable",
								Color: notion.ColorGreen,
							},
						},
						"tags": notion.DatabasePageProperty{
							ID:   "L9/e",
							Type: notion.DBPropTypeMultiSelect,
							MultiSelect: []notion.SelectOptions{
								{
									ID:    "d209b920-212c-4040-9d4a-bdf349dd8b2a",
									Name:  "math",
									Color: notion.ColorBlue,
								},
								{
									ID:    "6c3867c5-d542-4f84-b6e9-a420c43094e7",
									Name:  "history",
									Color: notion.ColorYellow,
								},
							},
						},
						"Age": notion.DatabasePageProperty{
							ID:     "$9nb",
							Type:   notion.DBPropTypeNumber,
							Name:   "Age",
							Number: notion.Float64Ptr(42),
						},
						"People": notion.DatabasePageProperty{
							ID:   "1#nc",
							Type: notion.DBPropTypePeople,
							Name: "People",
							People: []notion.User{
								{
									BaseUser: notion.BaseUser{
										ID: "be32e790-8292-46df-a248-b784fdf483cf",
									},
									Name:      "Jane Doe",
									AvatarURL: "https://example.com/image.png",
									Type:      notion.UserTypePerson,
									Person: &notion.Person{
										Email: "jane@example.com",
									},
								},
							},
						},
						"Files": notion.DatabasePageProperty{
							ID:   "!$9x",
							Type: notion.DBPropTypeFiles,
							Name: "Files",
							Files: []notion.File{
								{
									Name: "foobar.pdf",
								},
							},
						},
						"Checkbox": notion.DatabasePageProperty{
							ID:       "49S@",
							Type:     notion.DBPropTypeCheckbox,
							Name:     "Checkbox",
							Checkbox: notion.BoolPtr(true),
						},
						"Calculation": notion.DatabasePageProperty{
							ID:   "s(4f",
							Type: notion.DBPropTypeFormula,
							Name: "Calculation",
							Formula: &notion.FormulaResult{
								Type:   notion.FormulaResultTypeNumber,
								Number: notion.Float64Ptr(float64(42)),
							},
						},
						"URL": notion.DatabasePageProperty{
							ID:   "93$$",
							Type: notion.DBPropTypeURL,
							Name: "URL",
							URL:  notion.StringPtr("https://example.com"),
						},
						"Email": notion.DatabasePageProperty{
							ID:    "xb3Q",
							Type:  notion.DBPropTypeEmail,
							Name:  "Email",
							Email: notion.StringPtr("jane@example.com"),
						},
						"PhoneNumber": notion.DatabasePageProperty{
							ID:          "c2#Q",
							Type:        notion.DBPropTypePhoneNumber,
							Name:        "PhoneNumber",
							PhoneNumber: notion.StringPtr("867-5309"),
						},
						// "CreatedTime": notion.DatabasePageProperty{
						// 	ID:          "s#0s",
						// 	Type:        notion.DBPropTypeCreatedTime,
						// 	Name:        "Created time",
						// 	CreatedTime: notion.TimePtr(mustParseTime(time.RFC3339Nano, "2021-05-24T15:44:09.123Z")),
						// },
						"CreatedBy": notion.DatabasePageProperty{
							ID:   "49S@",
							Type: notion.DBPropTypeCreatedBy,
							Name: "Created by",
							CreatedBy: &notion.User{
								BaseUser: notion.BaseUser{
									ID: "be32e790-8292-46df-a248-b784fdf483cf",
								},
								Name:      "Jane Doe",
								AvatarURL: "https://example.com/image.png",
								Type:      notion.UserTypePerson,
								Person: &notion.Person{
									Email: "jane@example.com",
								},
							},
						},
						"LastEditedTime": notion.DatabasePageProperty{
							ID:             "x#0s",
							Type:           notion.DBPropTypeLastEditedTime,
							Name:           "Last edited time",
							LastEditedTime: notion.TimePtr(mustParseTime(time.RFC3339Nano, "2021-05-24T15:44:09.123Z")),
						},
						"LastEditedBy": notion.DatabasePageProperty{
							ID:   "x9S@",
							Type: notion.DBPropTypeLastEditedBy,
							Name: "Last edited by",
							LastEditedBy: &notion.User{
								BaseUser: notion.BaseUser{
									ID: "be32e790-8292-46df-a248-b784fdf483cf",
								},
								Name:      "Jane Doe",
								AvatarURL: "https://example.com/image.png",
								Type:      notion.UserTypePerson,
								Person: &notion.Person{
									Email: "jane@example.com",
								},
							},
						},
						"Relation": notion.DatabasePageProperty{
							ID:   "Cxl[",
							Type: notion.DBPropTypeRelation,
							Name: "Relation",
							Relation: []notion.Relation{
								{
									ID: "2be9597f-693f-4b87-baf9-efc545d38ebe",
								},
							},
						},
						"Rollup": notion.DatabasePageProperty{
							ID:   "xyA}",
							Type: notion.DBPropTypeRollup,
							Name: "Rollup",
							Rollup: &notion.RollupResult{
								Type: notion.RollupResultTypeArray,
								Array: []notion.DatabasePageProperty{
									{
										Type:   notion.DBPropTypeNumber,
										Number: notion.Float64Ptr(42),
									},
								},
							},
						},
					},
				},
			},
			want: Table{
				Rows: []Row{
					{
						Cells: []propTypeValue{
							{
								propType: notion.DBPropTypeTitle,
								propName: "Name",
								value:    "Foobar",
							},
							// {
							// 	propType: "created_time",
							// 	propName: "AGE",
							// 	value:    "575d",
							// },
							{
								propType: notion.DBPropTypeRichText,
								propName: "Description",
								value:    "Foobar",
							},
							{
								propType: notion.DBPropTypeSelect,
								propName: "Food group",
								value:    "🥦Vegetable",
							},
							{
								propType: notion.DBPropTypeMultiSelect,
								propName: "tags",
								value:    "math,history",
							},
						},
					},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := newTable(tt.arg)

			opts := []cmp.Option{
				cmpopts.SortSlices(func(i, j propTypeValue) bool {
					return i.propName < j.propName
				}),
				cmp.AllowUnexported(propTypeValue{}),
			}

			if diff := cmp.Diff(tt.want, got, opts...); diff != "" {
				t.Errorf("Table value is mismatch : %s\n", diff)

			}
		})
	}
}

func TestTablePrinterPrint(t *testing.T) {

	tests := []struct {
		name  string
		table Table
		want  string
	}{
		{
			name: "test full column",
			table: Table{
				Rows: []Row{
					{
						Cells: []propTypeValue{
							{
								propType: notion.DBPropTypeTitle,
								propName: "Name",
								value:    "Foobar",
							},
							{
								propType: "created_time",
								propName: "AGE",
								value:    "575d",
							},
							{
								propType: notion.DBPropTypeRichText,
								propName: "Description",
								value:    "Foobar",
							},
							{
								propType: notion.DBPropTypeSelect,
								propName: "Food group",
								value:    "🥦Vegetable",
							},
							{
								propType: notion.DBPropTypeMultiSelect,
								propName: "tags",
								value:    "math,history",
							},
						},
					},
				},
			},
			want: "NAME\tAGE\tFOOD GROUP\tTAGS\tDESCRIPTION\t\nFoobar\t575d\t🥦Vegetable\tmath,history\tFoobar\t\n",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			buf := &bytes.Buffer{}
			printer := &TablePrinter{table: tt.table, writer: buf}
			printer.Print()
			got := buf.String()

			if tt.want != got {
				t.Errorf("print value is mismatch. want: %s, got: %s", tt.want, got)
			}
		})
	}
}
