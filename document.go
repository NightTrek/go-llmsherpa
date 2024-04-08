package chipper

import "strings"

type Document struct {
	reader   *LayoutReader
	rootNode BlockInterface
	json     []interface{}
}

func NewDocument(blocksJSON []interface{}) *Document {
	reader := &LayoutReader{}
	rootNode := reader.Read(blocksJSON)
	return &Document{
		reader:   reader,
		rootNode: rootNode,
		json:     blocksJSON,
	}
}

func (d *Document) Chunks() []BlockInterface {
	return d.rootNode.Chunks()
}

func (d *Document) Tables() []BlockInterface {
	return d.rootNode.Tables()
}

func (d *Document) Sections() []BlockInterface {
	return d.rootNode.Sections()
}

func (d *Document) ToText() string {
	text := ""
	for _, section := range d.Sections() {
		text += section.ToText(true, true) + "\n"
	}
	return strings.TrimSpace(text)
}

func (d *Document) ToHTML() string {
	htmlStr := "<html>"
	for _, section := range d.Sections() {
		htmlStr += section.ToHTML(true, true)
	}
	htmlStr += "</html>"
	return htmlStr
}
