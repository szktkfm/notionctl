package util

import (
	"bytes"
	"fmt"
	"strings"
	"unicode/utf8"

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

	reg.Register(ast.KindDocument, r.renderDocument)
	reg.Register(ast.KindHeading, r.renderHeading)
	reg.Register(ast.KindBlockquote, r.renderBlockquote)
	reg.Register(ast.KindCodeBlock, r.renderCodeBlock)
	reg.Register(ast.KindFencedCodeBlock, r.renderFencedCodeBlock)
	reg.Register(ast.KindHTMLBlock, r.renderHTMLBlock)
	reg.Register(ast.KindList, r.renderList)
	reg.Register(ast.KindListItem, r.renderListItem)
	reg.Register(ast.KindParagraph, r.renderParagraph)
	reg.Register(ast.KindTextBlock, r.renderTextBlock)
	reg.Register(ast.KindThematicBreak, r.renderThematicBreak)

	// inlines

	reg.Register(ast.KindAutoLink, r.renderAutoLink)
	reg.Register(ast.KindEmphasis, r.renderEmphasis)
	// reg.Register(ast.KindImage, r.renderImage)
	reg.Register(ast.KindLink, r.renderLink)
	reg.Register(ast.KindRawHTML, r.renderRawHTML)
	reg.Register(ast.KindText, r.renderText)
	reg.Register(ast.KindString, r.renderString)

}

func (r *NotionRenderer) renderDocument(w util.BufWriter, source []byte, node ast.Node, entering bool) (ast.WalkStatus, error) {
	// nothing to do
	return ast.WalkContinue, nil
}

var BlockquoteAttributeFilter = GlobalAttributeFilter.Extend(
	[]byte("cite"),
)

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
				"heading_1": {
					"rich_text": [{ "type": "text", "text": { "content": "`
		} else if n.Level == 4 {
			s = `{
				"object": "block",
				"type": "heading_3",
				"heading_1": {
					"rich_text": [{ "type": "text", "text": { "content": "`

		} else {
			s = `{
				"object": "block",
				"type": "heading_3",
				"heading_1": {
					"rich_text": [{ "type": "text", "text": { "content": "`
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

// TODO: blockがnestedしている場合は考える
func (r *NotionRenderer) renderBlockquote(w util.BufWriter, source []byte, n ast.Node, entering bool) (ast.WalkStatus, error) {
	if entering {
		s := `{
			"object": "block",
			"type": "quote",
			"quote": {
				"rich_text": [
					{"type": "text",
					"text": {"content": ""}
					}
				],
				"children": [`

		_, _ = w.WriteString(
			strings.ReplaceAll(strings.ReplaceAll(s, "\n", ""), "\t", ""),
		)

	} else {

		_, _ = w.WriteString(
			`]}}`,
		)
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
		s := fmt.Sprintf(`{
			"object": "block","type": "code",
			"code": {
				"language": "%s", "rich_text": [{"type": "text","text": {"content": "`, language)

		_, _ = w.WriteString(
			strings.ReplaceAll(strings.ReplaceAll(s, "\n", ""), "\t", ""),
		)

		// _ = w.WriteByte('>')

		r.writeLines(w, source, n)
	} else {
		_, _ = w.WriteString(
			`"}}]}}`,
		)
	}
	return ast.WalkContinue, nil
}

func (r *NotionRenderer) renderHTMLBlock(w util.BufWriter, source []byte, node ast.Node, entering bool) (ast.WalkStatus, error) {
	n := node.(*ast.HTMLBlock)
	if entering {
		if r.Unsafe {
			l := n.Lines().Len()
			for i := 0; i < l; i++ {
				line := n.Lines().At(i)
				r.Writer.SecureWrite(w, line.Value(source))
			}
		} else {
			_, _ = w.WriteString("<!-- raw HTML omitted -->\n")
		}
	} else {
		if n.HasClosure() {
			if r.Unsafe {
				closure := n.ClosureLine
				r.Writer.SecureWrite(w, closure.Value(source))
			} else {
				_, _ = w.WriteString("<!-- raw HTML omitted -->\n")
			}
		}
	}
	return ast.WalkContinue, nil
}

// ListAttributeFilter defines attribute names which list elements can have.
var ListAttributeFilter = GlobalAttributeFilter.Extend(
	[]byte("start"),
	[]byte("reversed"),
	[]byte("type"),
)

func (r *NotionRenderer) renderList(w util.BufWriter, source []byte, node ast.Node, entering bool) (ast.WalkStatus, error) {
	n := node.(*ast.List)
	tag := "ul"
	if n.IsOrdered() {
		tag = "ol"
	}
	if entering {
		_ = w.WriteByte('<')
		_, _ = w.WriteString(tag)
		if n.IsOrdered() && n.Start != 1 {
			fmt.Fprintf(w, " start=\"%d\"", n.Start)
		}
		if n.Attributes() != nil {
			html.RenderAttributes(w, n, ListAttributeFilter)
		}
		_, _ = w.WriteString(">\n")
	} else {
		_, _ = w.WriteString("</")
		_, _ = w.WriteString(tag)
		_, _ = w.WriteString(">\n")
	}
	return ast.WalkContinue, nil
}

// ListItemAttributeFilter defines attribute names which list item elements can have.
var ListItemAttributeFilter = GlobalAttributeFilter.Extend(
	[]byte("value"),
)

func (r *NotionRenderer) renderListItem(w util.BufWriter, source []byte, n ast.Node, entering bool) (ast.WalkStatus, error) {
	if entering {
		if n.Attributes() != nil {
			_, _ = w.WriteString("<li")
			html.RenderAttributes(w, n, ListItemAttributeFilter)
			_ = w.WriteByte('>')
		} else {
			_, _ = w.WriteString("<li>")
		}
		fc := n.FirstChild()
		if fc != nil {
			if _, ok := fc.(*ast.TextBlock); !ok {
				_ = w.WriteByte('\n')
			}
		}
	} else {
		_, _ = w.WriteString("</li>\n")
	}
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

// ThematicAttributeFilter defines attribute names which hr elements can have.
var ThematicAttributeFilter = GlobalAttributeFilter.Extend(
	[]byte("align"),   // [Deprecated]
	[]byte("color"),   // [Not Standardized]
	[]byte("noshade"), // [Deprecated]
	[]byte("size"),    // [Deprecated]
	[]byte("width"),   // [Deprecated]
)

func (r *NotionRenderer) renderThematicBreak(w util.BufWriter, source []byte, n ast.Node, entering bool) (ast.WalkStatus, error) {
	if !entering {
		return ast.WalkContinue, nil
	}
	_, _ = w.WriteString("<hr")
	if n.Attributes() != nil {
		html.RenderAttributes(w, n, ThematicAttributeFilter)
	}
	if r.XHTML {
		_, _ = w.WriteString(" />\n")
	} else {
		_, _ = w.WriteString(">\n")
	}
	return ast.WalkContinue, nil
}

// LinkAttributeFilter defines attribute names which link elements can have.
var LinkAttributeFilter = GlobalAttributeFilter.Extend(
	[]byte("download"),
	// []byte("href"),
	[]byte("hreflang"),
	[]byte("media"),
	[]byte("ping"),
	[]byte("referrerpolicy"),
	[]byte("rel"),
	[]byte("shape"),
	[]byte("target"),
)

func (r *NotionRenderer) renderAutoLink(w util.BufWriter, source []byte, node ast.Node, entering bool) (ast.WalkStatus, error) {
	n := node.(*ast.AutoLink)
	if !entering {
		return ast.WalkContinue, nil
	}
	_, _ = w.WriteString(`<a href="`)
	url := n.URL(source)
	label := n.Label(source)
	if n.AutoLinkType == ast.AutoLinkEmail && !bytes.HasPrefix(bytes.ToLower(url), []byte("mailto:")) {
		_, _ = w.WriteString("mailto:")
	}
	_, _ = w.Write(util.EscapeHTML(util.URLEscape(url, false)))
	if n.Attributes() != nil {
		_ = w.WriteByte('"')
		html.RenderAttributes(w, n, LinkAttributeFilter)
		_ = w.WriteByte('>')
	} else {
		_, _ = w.WriteString(`">`)
	}
	_, _ = w.Write(util.EscapeHTML(label))
	_, _ = w.WriteString(`</a>`)
	return ast.WalkContinue, nil
}

// CodeAttributeFilter defines attribute names which code elements can have.
var CodeAttributeFilter = GlobalAttributeFilter

// EmphasisAttributeFilter defines attribute names which emphasis elements can have.
var EmphasisAttributeFilter = GlobalAttributeFilter

func (r *NotionRenderer) renderEmphasis(w util.BufWriter, source []byte, node ast.Node, entering bool) (ast.WalkStatus, error) {
	n := node.(*ast.Emphasis)
	tag := "em"
	if n.Level == 2 {
		tag = "strong"
	}
	if entering {
		_ = w.WriteByte('<')
		_, _ = w.WriteString(tag)
		if n.Attributes() != nil {
			html.RenderAttributes(w, n, EmphasisAttributeFilter)
		}
		_ = w.WriteByte('>')
	} else {
		_, _ = w.WriteString("</")
		_, _ = w.WriteString(tag)
		_ = w.WriteByte('>')
	}
	return ast.WalkContinue, nil
}

func (r *NotionRenderer) renderLink(w util.BufWriter, source []byte, node ast.Node, entering bool) (ast.WalkStatus, error) {
	n := node.(*ast.Link)
	if entering {
		_, _ = w.WriteString("<a href=\"")
		if r.Unsafe || !html.IsDangerousURL(n.Destination) {
			_, _ = w.Write(util.EscapeHTML(util.URLEscape(n.Destination, true)))
		}
		_ = w.WriteByte('"')
		if n.Title != nil {
			_, _ = w.WriteString(` title="`)
			r.Writer.Write(w, n.Title)
			_ = w.WriteByte('"')
		}
		if n.Attributes() != nil {
			html.RenderAttributes(w, n, LinkAttributeFilter)
		}
		_ = w.WriteByte('>')
	} else {
		_, _ = w.WriteString("</a>")
	}
	return ast.WalkContinue, nil
}

// ImageAttributeFilter defines attribute names which image elements can have.
var ImageAttributeFilter = GlobalAttributeFilter.Extend(
	[]byte("align"),
	[]byte("border"),
	[]byte("crossorigin"),
	[]byte("decoding"),
	[]byte("height"),
	[]byte("importance"),
	[]byte("intrinsicsize"),
	[]byte("ismap"),
	[]byte("loading"),
	[]byte("referrerpolicy"),
	[]byte("sizes"),
	[]byte("srcset"),
	[]byte("usemap"),
	[]byte("width"),
)

// func (r *NotionRenderer) renderImage(w util.BufWriter, source []byte, node ast.Node, entering bool) (ast.WalkStatus, error) {
// 	if !entering {
// 		return ast.WalkContinue, nil
// 	}
// 	n := node.(*ast.Image)
// 	_, _ = w.WriteString("<img src=\"")
// 	if r.Unsafe || !html.IsDangerousURL(n.Destination) {
// 		_, _ = w.Write(util.EscapeHTML(util.URLEscape(n.Destination, true)))
// 	}
// 	_, _ = w.WriteString(`" alt="`)
// 	_, _ = w.Write(html.nodeToHTMLText(n, source))
// 	_ = w.WriteByte('"')
// 	if n.Title != nil {
// 		_, _ = w.WriteString(` title="`)
// 		r.Writer.Write(w, n.Title)
// 		_ = w.WriteByte('"')
// 	}
// 	if n.Attributes() != nil {
// 		html.RenderAttributes(w, n, ImageAttributeFilter)
// 	}
// 	if r.XHTML {
// 		_, _ = w.WriteString(" />")
// 	} else {
// 		_, _ = w.WriteString(">")
// 	}
// 	return ast.WalkSkipChildren, nil
// }

func (r *NotionRenderer) renderRawHTML(w util.BufWriter, source []byte, node ast.Node, entering bool) (ast.WalkStatus, error) {
	if !entering {
		return ast.WalkSkipChildren, nil
	}
	if r.Unsafe {
		n := node.(*ast.RawHTML)
		l := n.Segments.Len()
		for i := 0; i < l; i++ {
			segment := n.Segments.At(i)
			_, _ = w.Write(segment.Value(source))
		}
		return ast.WalkSkipChildren, nil
	}
	_, _ = w.WriteString("<!-- raw HTML omitted -->")
	return ast.WalkSkipChildren, nil
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

type notionExtension struct {
}

// NotionExtension is an extension that allow you to use GFM task lists.
var NotionExtension = &notionExtension{}

func (e *notionExtension) Extend(m goldmark.Markdown) {
	// checkbox parserはDefaultParserに含まれていないので、このoptionは追加する必要があるが、、面倒だから飛ばそうかな。
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
