package util

import (
	"bytes"
	"encoding/json"
	"testing"

	"github.com/dstotijn/go-notion"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/yuin/goldmark"
)

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

	got := mdToNotionBlock(source)
	opt := cmpopts.IgnoreUnexported(notion.Heading2Block{}, notion.Heading1Block{})
	if diff := cmp.Diff(want, got, opt); diff != "" {
		t.Errorf("Table value is mismatch : %s\n", diff)
	}
}

func mdToNotionBlock(source []byte) []notion.Block {
	// 共通的な処理なので外だししい。
	var buf bytes.Buffer

	buf.WriteString(`{"results": [`)

	md := goldmark.New(
		goldmark.WithExtensions(NotionExtension),
	)
	md.Convert(source, &buf)

	buf.Truncate(buf.Len() - 1)
	buf.WriteString(`]}`)

	var got notion.BlockChildrenResponse
	json.Unmarshal(buf.Bytes(), &got)
	return got.Results
}
