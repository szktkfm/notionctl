package util

import (
	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/renderer"
	"github.com/yuin/goldmark/renderer/html"
	"github.com/yuin/goldmark/util"
)

// NewTaskCheckBoxParser returns a new  InlineParser that can parse
// checkboxes in list items.
// This parser must take precedence over the parser.LinkParser.

// NotionRenderer is a renderer.NodeRenderer implementation that
// renders checkboxes in list items.
type NotionRenderer struct {
	html.Config
}

// NewNotionRenderer returns a new TaskCheckBoxHTMLRenderer.
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
}

var GlobalAttributeFilter = util.NewBytesFilter(
	[]byte("accesskey"),
	[]byte("autocapitalize"),
	[]byte("autofocus"),
	[]byte("class"),
	[]byte("contenteditable"),
	[]byte("dir"),
	[]byte("draggable"),
	[]byte("enterkeyhint"),
	[]byte("hidden"),
	[]byte("id"),
	[]byte("inert"),
	[]byte("inputmode"),
	[]byte("is"),
	[]byte("itemid"),
	[]byte("itemprop"),
	[]byte("itemref"),
	[]byte("itemscope"),
	[]byte("itemtype"),
	[]byte("lang"),
	[]byte("part"),
	[]byte("slot"),
	[]byte("spellcheck"),
	[]byte("style"),
	[]byte("tabindex"),
	[]byte("title"),
	[]byte("translate"),
)

var HeadingAttributeFilter = GlobalAttributeFilter

// ParagraphAttributeFilter defines attribute names which paragraph elements can have.
var ParagraphAttributeFilter = GlobalAttributeFilter

func (r *NotionRenderer) renderHeading(w util.BufWriter, source []byte, node ast.Node, entering bool) (ast.WalkStatus, error) {
	n := node.(*ast.Heading)
	if entering {
		_, _ = w.WriteString(
			`{
			"object": "block",
			"type": "heading_2",
			"heading_2": {
				"rich_text": [{ "type": "text", "text": { "content": "`,
		)
		// _ = w.WriteByte("0123456"[n.Level])
		if n.Attributes() != nil {
			html.RenderAttributes(w, node, HeadingAttributeFilter)
		}
		// _ = w.WriteByte('>')
	} else {
		// _, _ = w.WriteString("</h")
		// _ = w.WriteByte("0123456"[n.Level])
		// _, _ = w.WriteString(">\n")
		_, _ = w.WriteString(
			`" } }]
			}
		},`,
		)
	}
	return ast.WalkContinue, nil
}

func (r *NotionRenderer) renderParagraph(w util.BufWriter, source []byte, n ast.Node, entering bool) (ast.WalkStatus, error) {
	// textの中の改行の扱いどうしようか
	if entering {
		_, _ = w.WriteString(
			`{"object": "block",
			"type": "paragraph",
			"paragraph": {
				"rich_text": [
					{
						"type": "text",
						"text": {
						"content": "`,
		)
	} else {
		_, _ = w.WriteString(
			`"}}]}},`,
		)
	}
	return ast.WalkContinue, nil
}

func (r *NotionRenderer) renderHeading2(w util.BufWriter, source []byte, node ast.Node, entering bool) (ast.WalkStatus, error) {
	n := node.(*ast.Heading)
	if entering {
		_, _ = w.WriteString("uwagaki2 <h")
		_ = w.WriteByte("0123456"[n.Level])
		if n.Attributes() != nil {
			html.RenderAttributes(w, node, HeadingAttributeFilter)
		}
		_ = w.WriteByte('>')
	} else {
		_, _ = w.WriteString("</h")
		_ = w.WriteByte("0123456"[n.Level])
		_, _ = w.WriteString(">\n")
	}
	return ast.WalkContinue, nil
}

type taskList struct {
}

// NotionExtension is an extension that allow you to use GFM task lists.
var NotionExtension = &taskList{}

func (e *taskList) Extend(m goldmark.Markdown) {
	// checkbox parserはDefaultParserに含まれていないので、このoptionは追加する必要があるが、、面倒だから飛ばそうかな。
	// m.Parser().AddOptions(parser.WithInlineParsers(
	// 	util.Prioritized(NewTaskCheckBoxParser(), 0),
	// ))
	m.Renderer().AddOptions(renderer.WithNodeRenderers(
		util.Prioritized(NewNotionRenderer(), 500),
	))
}
