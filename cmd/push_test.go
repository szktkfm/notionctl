package cmd

import (
	"bytes"
	"testing"

	"github.com/dstotijn/go-notion"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
)

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
			pushOptions := newTestPushOptions()
			got := pushOptions.newCreatePageParams()

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
