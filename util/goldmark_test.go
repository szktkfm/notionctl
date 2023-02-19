package util

import (
	"testing"

	"github.com/dstotijn/go-notion"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
)

// func testrenderblockquote(t *testing.T) {
// 	source := []byte(
// 		`> foo
// `,
// 	)
// 	want := []notion.Block{
// 		&notion.QuoteBlock{
// 			RichText: []notion.RichText{
// 				{
// 					Type: notion.RichTextTypeText,
// 					Text: &notion.Text{
// 						Content: "foo",
// 					},
// 				},
// 			},
// 		},
// 	}

// 	got := MDToNotionBlock(source)
// 	opt := cmpopts.IgnoreUnexported(notion.QuoteBlock{})
// 	if diff := cmp.Diff(want, got, opt); diff != "" {
// 		t.Errorf("Table value is mismatch : %s\n", diff)
// 	}
// }

func TestRenderHeading(t *testing.T) {
	source := []byte(
		`# title1 
## title2
		`,
	)
	want := []notion.Block{
		&notion.Heading1Block{
			RichText: []notion.RichText{
				{
					Type: notion.RichTextTypeText,
					Text: &notion.Text{
						Content: "title1",
					},
				},
			},
		},
		&notion.Heading2Block{
			RichText: []notion.RichText{
				{
					Type: notion.RichTextTypeText,
					Text: &notion.Text{
						Content: "title2",
					},
				},
			},
		},
	}

	got := MDToNotionBlock(source)
	opt := cmpopts.IgnoreUnexported(notion.Heading2Block{}, notion.Heading1Block{})
	if diff := cmp.Diff(want, got, opt); diff != "" {
		t.Errorf("Table value is mismatch : %s\n", diff)
	}
}
