package markdown

import (
	"fmt"
	"strings"
	"unicode/utf8"

	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/renderer"
	"github.com/yuin/goldmark/renderer/html"
	"github.com/yuin/goldmark/util"
)

type NotionRenderer struct {
	html.Config
}

func NewNotionRenderer(opts ...html.Option) renderer.NodeRenderer {
	r := &NotionRenderer{
		Config: html.NewConfig(),
	}
	for _, opt := range opts {
		opt.SetHTMLOption(&r.Config)
	}
	return r
}

// RegisterFuncs implements renderer.NodeRenderer.RegisterFuncs.
func (r *NotionRenderer) RegisterFuncs(reg renderer.NodeRendererFuncRegisterer) {
	// reg.Register(ast.KindTaskCheckBox, r.renderTaskCheckBox)
	reg.Register(ast.KindHeading, r.renderHeading)
	reg.Register(ast.KindParagraph, r.renderParagraph)

	reg.Register(ast.KindDocument, r.renderDocument)
	reg.Register(ast.KindHeading, r.renderHeading)
	reg.Register(ast.KindBlockquote, r.renderBlockquote)
	reg.Register(ast.KindCodeBlock, r.renderCodeBlock)
	reg.Register(ast.KindFencedCodeBlock, r.renderFencedCodeBlock)
	reg.Register(ast.KindList, r.renderList)
	reg.Register(ast.KindListItem, r.renderListItem)
	reg.Register(ast.KindParagraph, r.renderParagraph)
	reg.Register(ast.KindTextBlock, r.renderTextBlock)

	// inlines

	reg.Register(ast.KindLink, r.renderLink)
	reg.Register(ast.KindText, r.renderText)
	reg.Register(ast.KindString, r.renderString)

}

func (r *NotionRenderer) renderDocument(w util.BufWriter, source []byte, node ast.Node, entering bool) (ast.WalkStatus, error) {
	// nothing to do
	return ast.WalkContinue, nil
}

func (r *NotionRenderer) renderParagraph(w util.BufWriter, source []byte, n ast.Node, entering bool) (ast.WalkStatus, error) {
	if entering {
		s := `{"object": "block",
			"type": "paragraph",
			"paragraph": {
				"rich_text": [
					{
						"type": "text",
						"text": {
						"content": "`

		_, _ = w.WriteString(
			strings.ReplaceAll(strings.ReplaceAll(s, "\n", ""), "\t", ""),
		)

	} else {
		_, _ = w.WriteString(
			`"}}]}}`,
		)
	}
	return ast.WalkContinue, nil
}

func (r *NotionRenderer) renderHeading(w util.BufWriter, source []byte, node ast.Node, entering bool) (ast.WalkStatus, error) {
	n := node.(*ast.Heading)
	if entering {
		var s string
		if n.Level == 1 {
			s = `{
				"object": "block",
				"type": "heading_1",
				"heading_1": {
					"rich_text": [{ "type": "text", "text": { "content": "`

		} else if n.Level == 2 {
			s = `{
				"object": "block",
				"type": "heading_2",
				"heading_2": {
					"rich_text": [{ "type": "text", "text": { "content": "`
		} else if n.Level == 4 {
			s = `{
				"object": "block",
				"type": "heading_3",
				"heading_3": {
					"rich_text": [{ "type": "text", "text": { "content": "`

		} else {
			s = `{"object": "block",
			"type": "paragraph",
			"paragraph": {
				"rich_text": [
					{
						"type": "text",
						"text": {
						"content": "`
		}

		_, _ = w.WriteString(
			strings.ReplaceAll(strings.ReplaceAll(s, "\n", ""), "\t", ""),
		)
	} else {
		_, _ = w.WriteString(
			`"}}]}}`,
		)
	}
	return ast.WalkContinue, nil
}

// TODO: nested block
func (r *NotionRenderer) renderBlockquote(w util.BufWriter, source []byte, n ast.Node, entering bool) (ast.WalkStatus, error) {
	if entering {
		// s := `{
		// 	"object": "block",
		// 	"type": "quote",
		// 	"quote": {
		// 		"rich_text": [
		// 			{"type": "text",
		// 			"text": {"content": ""}
		// 			}
		// 		],
		// 		"children": [`

		// _, _ = w.WriteString(
		// 	strings.ReplaceAll(strings.ReplaceAll(s, "\n", ""), "\t", ""),
		// )

	} else {

		// _, _ = w.WriteString(
		// 	`]}}`,
		// )
	}
	return ast.WalkContinue, nil
}

func (r *NotionRenderer) renderCodeBlock(w util.BufWriter, source []byte, n ast.Node, entering bool) (ast.WalkStatus, error) {
	if entering {
		r.writeLines(w, source, n)
	} else {
		_, _ = w.WriteString("</code></pre>\n")
	}
	return ast.WalkContinue, nil
}

func (r *NotionRenderer) renderFencedCodeBlock(w util.BufWriter, source []byte, node ast.Node, entering bool) (ast.WalkStatus, error) {
	n := node.(*ast.FencedCodeBlock)
	if entering {
		language := n.Language(source)
		if len(language) == 0 {
			language = []byte("go")
		}
		s := fmt.Sprintf(`{
			"object": "block","type": "code",
			"code": {
				"language": "%s", "rich_text": [{"type": "text","text": {"content": "`, language)

		_, _ = w.WriteString(
			strings.ReplaceAll(strings.ReplaceAll(s, "\n", ""), "\t", ""),
		)

		r.writeLines(w, source, n)
	} else {
		_, _ = w.WriteString(
			`"}}]}}`,
		)
	}
	return ast.WalkContinue, nil
}

// TODO: number list
func (r *NotionRenderer) renderList(w util.BufWriter, source []byte, node ast.Node, entering bool) (ast.WalkStatus, error) {
	return ast.WalkContinue, nil
}

// TODO: If the list is nested, this function is called recursively and the json is not created correctly
func (r *NotionRenderer) renderListItem(w util.BufWriter, source []byte, n ast.Node, entering bool) (ast.WalkStatus, error) {
	// if entering {
	// 	s := `{"object": "block",
	// 		"type": "bulleted_list_item",
	// 		"bulleted_list_item": {
	// 			"rich_text": [
	// 				{
	// 					"type": "text",
	// 					"text": {
	// 					"content": "`

	// 	_, _ = w.WriteString(
	// 		strings.ReplaceAll(strings.ReplaceAll(s, "\n", ""), "\t", ""),
	// 	)
	// } else {
	// 	_, _ = w.WriteString(
	// 		`"}}]}}`,
	// 	)
	// }

	return ast.WalkContinue, nil
}

func (r *NotionRenderer) renderTextBlock(w util.BufWriter, source []byte, n ast.Node, entering bool) (ast.WalkStatus, error) {
	if !entering {
		if _, ok := n.NextSibling().(ast.Node); ok && n.FirstChild() != nil {
			_ = w.WriteByte('\n')
		}
	}
	return ast.WalkContinue, nil
}

func (r *NotionRenderer) renderLink(w util.BufWriter, source []byte, node ast.Node, entering bool) (ast.WalkStatus, error) {
	n := node.(*ast.Link)
	if entering {

	} else {
		s := `", "link": {"url": "`
		_, _ = w.WriteString(s)

		if r.Unsafe || !html.IsDangerousURL(n.Destination) {
			_, _ = w.Write(util.EscapeHTML(util.URLEscape(n.Destination, true)))
		}
		_, _ = w.WriteString(
			`"}`,
		)
	}
	return ast.WalkContinue, nil
}

func (r *NotionRenderer) renderText(w util.BufWriter, source []byte, node ast.Node, entering bool) (ast.WalkStatus, error) {
	if !entering {
		return ast.WalkContinue, nil
	}
	n := node.(*ast.Text)
	segment := n.Segment
	if n.IsRaw() {
		r.Writer.RawWrite(w, segment.Value(source))
	} else {
		value := segment.Value(source)
		r.Writer.Write(w, value)
		if n.HardLineBreak() || (n.SoftLineBreak() && r.HardWraps) {
			if r.XHTML {
				_, _ = w.WriteString("<br />\n")
			} else {
				_, _ = w.WriteString("<br>\n")
			}
		} else if n.SoftLineBreak() {
			if r.EastAsianLineBreaks && len(value) != 0 {
				sibling := node.NextSibling()
				if sibling != nil && sibling.Kind() == ast.KindText {
					if siblingText := sibling.(*ast.Text).Text(source); len(siblingText) != 0 {
						thisLastRune := util.ToRune(value, len(value)-1)
						siblingFirstRune, _ := utf8.DecodeRune(siblingText)
						if !(util.IsEastAsianWideRune(thisLastRune) &&
							util.IsEastAsianWideRune(siblingFirstRune)) {
							_ = w.WriteByte('\n')
						}
					}
				}
			} else {
				_ = w.WriteByte('\n')
			}
		}
	}
	return ast.WalkContinue, nil
}

func (r *NotionRenderer) renderString(w util.BufWriter, source []byte, node ast.Node, entering bool) (ast.WalkStatus, error) {
	if !entering {
		return ast.WalkContinue, nil
	}
	n := node.(*ast.String)
	if n.IsCode() {
		_, _ = w.Write(n.Value)
	} else {
		if n.IsRaw() {
			r.Writer.RawWrite(w, n.Value)
		} else {
			r.Writer.Write(w, n.Value)
		}
	}
	return ast.WalkContinue, nil
}

type notionExtension struct {
}

var NotionExtension = &notionExtension{}

func (e *notionExtension) Extend(m goldmark.Markdown) {
	m.Renderer().AddOptions(renderer.WithNodeRenderers(
		util.Prioritized(NewNotionRenderer(), 500),
	))
}

func (r *NotionRenderer) writeLines(w util.BufWriter, source []byte, n ast.Node) {
	l := n.Lines().Len()
	for i := 0; i < l; i++ {
		line := n.Lines().At(i)
		r.Writer.RawWrite(w, line.Value(source))
	}
}
