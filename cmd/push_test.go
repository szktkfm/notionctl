package cmd

import (
	"bytes"
	"fmt"
	"testing"

	"github.com/dstotijn/go-notion"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
)

func TestPushCmdRun(t *testing.T) {

	tests := []struct {
		name      string
		option    *GetOptions
		queryMock string
		want      string
	}{
		{
			name: "test full column",
			option: &GetOptions{
				DB:   "hogedb",
				Page: "",
				Wide: false,
			},
			// TODO: define []notion.Database and unmashal it
			queryMock: `{"object":"list","results":[{"object":"page","id":"711e6ef1-28c6-482e-91f2-75dd26dfd041","created_time":"2022-12-05T11:58:00.000Z","last_edited_time":"2022-12-24T16:44:00.000Z","created_by":{"object":"user","id":"5a639bd5-f786-4565-bc65-5d9281ef3944"},"last_edited_by":{"object":"user","id":"25b9e72d-a868-4007-bf74-841efc304d3e"},"cover":{"type":"external","external":{"url":"https://upload.wikimedia.org/wikipedia/commons/6/62/Tuscankale.jpg"}},"icon":{"type":"emoji","emoji":"🥬"},"parent":{"type":"database_id","database_id":"98079428-d5d0-436f-a316-b2d36da049c2"},"archived":false,"properties":{"Food group":{"id":"B%60Ts","type":"select","select":{"id":"26ab5fc8-7e6b-4d11-b6c5-6864e614c3ed","name":"Vegetable","color":"purple"}},"Description":{"id":"oBRk","type":"rich_text","rich_text":[{"type":"text","text":{"content":"A dark green leafy vegetable","link":null},"annotations":{"bold":false,"italic":false,"strikethrough":false,"underline":false,"code":false,"color":"default"},"plain_text":"A dark green leafy vegetable","href":null}]},"Created time":{"id":"rlOQ","type":"created_time","created_time":"2022-12-05T11:58:00.000Z"},"Tags":{"id":"urJ%5B","type":"multi_select","multi_select":[]},"Price":{"id":"xyUL","type":"number","number":2.5},"Name":{"id":"title","type":"title","title":[{"type":"text","text":{"content":"Tuscan Kale","link":null},"annotations":{"bold":false,"italic":false,"strikethrough":false,"underline":false,"code":false,"color":"default"},"plain_text":"Tuscan Kale","href":null}]}},"url":"https://www.notion.so/Tuscan-Kale-711e6ef128c6482e91f275dd26dfd041"}],"next_cursor":null,"has_more":false,"type":"page","page":{}}`,
			want:      "NAME\tAGE\tFOOD GROUP\tTAGS\tDESCRIPTION\t\nTuscan Kale\t19d\tVegetable\t-\tA dark green leafy veget\t\n",
		},
	}
	for _, tt := range tests {
		fmt.Println(tt)
	}
}

func TestNewCreatePageParam(t *testing.T) {

	tests := []struct {
		name string
		want notion.CreatePageParams
	}{
		{
			name: "test full column",
			// TODO: define []notion.Database and unmashal it
			want: notion.CreatePageParams{
				ParentType: "database_id",
				DatabasePageProperties: &notion.DatabasePageProperties{
					"Description": notion.DatabasePageProperty{
						RichText: []notion.RichText{
							{
								Text: &notion.Text{
									Content: "test description",
								},
							},
						},
					},
					"Name": notion.DatabasePageProperty{
						Title: []notion.RichText{
							{
								Text: &notion.Text{
									Content: "test title",
								},
							},
						},
					},
				},
				Children: []notion.Block{
					notion.ParagraphBlock{
						RichText: []notion.RichText{
							{
								Text: &notion.Text{
									Content: "paragraph line",
								},
							},
						},
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			option := newTestPushOptions()
			got := option.newCreatePageParams()

			opt := cmpopts.IgnoreUnexported(notion.ParagraphBlock{})
			if diff := cmp.Diff(tt.want, got, opt); diff != "" {
				t.Errorf("Table value is mismatch : %s\n", diff)

			}

		})
	}
}

func newTestPushOptions() *PushOptions {
	buf := new(bytes.Buffer)
	buf.Write([]byte("paragraph line"))

	return &PushOptions{
		Title:       "test title",
		Description: "test description",
		In:          buf,
		targetDB: notion.Database{
			ID: "668d797c-76fa-4934-9b05-ad288df2d136",
			CreatedBy: notion.BaseUser{
				ID: "71e95936-2737-4e11-b03d-f174f6f13087",
			},
			LastEditedBy: notion.BaseUser{
				ID: "5ba97cc9-e5e0-4363-b33a-1d80a635577f",
			},
			URL: "https://www.notion.so/668d797c76fa49349b05ad288df2d136",
			Title: []notion.RichText{
				{
					Type: notion.RichTextTypeText,
					Text: &notion.Text{
						Content: "Grocery List",
					},
					Annotations: &notion.Annotations{
						Color: notion.ColorDefault,
					},
					PlainText: "Grocery List",
				},
			},
			Properties: notion.DatabaseProperties{
				"Name": notion.DatabaseProperty{
					ID:    "title",
					Type:  notion.DBPropTypeTitle,
					Title: &notion.EmptyMetadata{},
				},
				"Description": notion.DatabaseProperty{
					ID:   "J@cS",
					Type: notion.DBPropTypeRichText,
				},
				"In stock": notion.DatabaseProperty{
					ID:       "{xYx",
					Type:     notion.DBPropTypeCheckbox,
					Checkbox: &notion.EmptyMetadata{},
				},
				"Food group": notion.DatabaseProperty{
					ID:   "TJmr",
					Type: notion.DBPropTypeSelect,
					Select: &notion.SelectMetadata{
						Options: []notion.SelectOptions{
							{
								ID:    "96eb622f-4b88-4283-919d-ece2fbed3841",
								Name:  "🥦Vegetable",
								Color: notion.ColorGreen,
							},
							{
								ID:    "bb443819-81dc-46fb-882d-ebee6e22c432",
								Name:  "🍎Fruit",
								Color: notion.ColorRed,
							},
							{
								ID:    "7da9d1b9-8685-472e-9da3-3af57bdb221e",
								Name:  "💪Protein",
								Color: notion.ColorYellow,
							},
						},
					},
				},
				"Price": notion.DatabaseProperty{
					ID:   "cU^N",
					Type: notion.DBPropTypeNumber,
					Number: &notion.NumberMetadata{
						Format: notion.NumberFormatDollar,
					},
				},
				"Cost of next trip": {
					ID:   "p:sC",
					Type: notion.DBPropTypeFormula,
					Formula: &notion.FormulaMetadata{
						Expression: `if(prop("In stock"), 0, prop("Price"))`,
					},
				},
				"Last ordered": notion.DatabaseProperty{
					ID:   "]\\R[",
					Type: notion.DBPropTypeDate,
					Date: &notion.EmptyMetadata{},
				},
				"Meals": notion.DatabaseProperty{
					ID:   "lV]M",
					Type: notion.DBPropTypeRelation,
					Relation: &notion.RelationMetadata{
						DatabaseID: "668d797c-76fa-4934-9b05-ad288df2d136",
						Type:       notion.RelationTypeDualProperty,
						DualProperty: &notion.DualPropertyRelation{
							SyncedPropID:   "IJi<",
							SyncedPropName: "Related to Test database (Relation Test)",
						},
					},
				},
				"Number of meals": notion.DatabaseProperty{
					ID:   "Z\\Eh",
					Type: notion.DBPropTypeRollup,
					Rollup: &notion.RollupMetadata{
						RollupPropName:   "Name",
						RelationPropName: "Meals",
						RollupPropID:     "title",
						RelationPropID:   "mxp^",
						Function:         notion.RollupFunctionCountAll,
					},
				},
				"Store availability": notion.DatabaseProperty{
					ID:   "=_>D",
					Type: notion.DBPropTypeMultiSelect,
					MultiSelect: &notion.SelectMetadata{
						Options: []notion.SelectOptions{
							{
								ID:    "d209b920-212c-4040-9d4a-bdf349dd8b2a",
								Name:  "Duc Loi Market",
								Color: notion.ColorBlue,
							},
							{
								ID:    "6c3867c5-d542-4f84-b6e9-a420c43094e7",
								Name:  "Gus's Community Market",
								Color: notion.ColorYellow,
							},
						},
					},
				},
				"+1": notion.DatabaseProperty{
					ID:     "aGut",
					Type:   notion.DBPropTypePeople,
					People: &notion.EmptyMetadata{},
				},
				"Photo": {
					ID:    "aTIT",
					Type:  "files",
					Files: &notion.EmptyMetadata{},
				},
			},
			Parent: notion.Parent{
				Type:   notion.ParentTypePage,
				PageID: "b8595b75-abd1-4cad-8dfe-f935a8ef57cb",
			},
		},
	}
}
