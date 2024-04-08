package chipper

import (
	"fmt"
	"strings"
)

// Block is the base interface for all block types
type Block interface {
	Tag() string
	Level() int
	PageIdx() int
	BlockIdx() int
	Top() float64
	Left() float64
	Bbox() []float64
	Sentences() []string
	Children() []Block
	Parent() Block
	BlockJSON() map[string]interface{}
	AddChild(node Block)
	ToHTML(includeChildren, recurse bool) string
	ToText(includeChildren, recurse bool) string
	ParentChain() []Block
	ParentText() string
	ToContextText(includeSectionInfo bool) string
	IterChildren(node Block, level int, nodeVisitor func(Block))
	Paragraphs() []Block
	Chunks() []Block
	Tables() []Block
	Sections() []Block
}

// BaseBlock is the base struct for all block types
type BaseBlock struct {
	tag       string
	level     int
	pageIdx   int
	blockIdx  int
	top       float64
	left      float64
	bbox      []float64
	sentences []string
	children  []Block
	parent    Block
	blockJSON map[string]interface{}
}

// NewBaseBlock creates a new BaseBlock instance
func NewBaseBlock(blockJSON map[string]interface{}) *BaseBlock {
	tag, _ := blockJSON["tag"].(string)
	level, _ := blockJSON["level"].(int)
	pageIdx, _ := blockJSON["page_idx"].(int)
	blockIdx, _ := blockJSON["block_idx"].(int)
	top, _ := blockJSON["top"].(float64)
	left, _ := blockJSON["left"].(float64)
	bbox, _ := blockJSON["bbox"].([]float64)
	sentences, _ := blockJSON["sentences"].([]string)

	return &BaseBlock{
		tag:       tag,
		level:     level,
		pageIdx:   pageIdx,
		blockIdx:  blockIdx,
		top:       top,
		left:      left,
		bbox:      bbox,
		sentences: sentences,
		children:  []Block{},
		parent:    nil,
		blockJSON: blockJSON,
	}
}

// Tag returns the tag of the block
func (b *BaseBlock) Tag() string {
	return b.tag
}

// Level returns the level of the block
func (b *BaseBlock) Level() int {
	return b.level
}

// PageIdx returns the page index of the block
func (b *BaseBlock) PageIdx() int {
	return b.pageIdx
}

// BlockIdx returns the block index of the block
func (b *BaseBlock) BlockIdx() int {
	return b.blockIdx
}

// Top returns the top position of the block
func (b *BaseBlock) Top() float64 {
	return b.top
}

// Left returns the left position of the block
func (b *BaseBlock) Left() float64 {
	return b.left
}

// Bbox returns the bounding box of the block
func (b *BaseBlock) Bbox() []float64 {
	return b.bbox
}

// Sentences returns the sentences of the block
func (b *BaseBlock) Sentences() []string {
	return b.sentences
}

// Children returns the immediate child blocks of the block
func (b *BaseBlock) Children() []Block {
	return b.children
}

// Parent returns the parent block of the block
func (b *BaseBlock) Parent() Block {
	return b.parent
}

// BlockJSON returns the JSON representation of the block
func (b *BaseBlock) BlockJSON() map[string]interface{} {
	return b.blockJSON
}

// AddChild adds a child block to the block and sets the parent of the child
func (b *BaseBlock) AddChild(node Block) {
	b.children = append(b.children, node)
	node.(*BaseBlock).parent = b
}

// ToHTML converts the block to HTML (not implemented in the base struct)
func (b *BaseBlock) ToHTML(includeChildren, recurse bool) string {
	return ""
}

// ToText converts the block to text (not implemented in the base struct)
func (b *BaseBlock) ToText(includeChildren, recurse bool) string {
	return ""
}

// ParentChain returns the parent chain of the block
func (b *BaseBlock) ParentChain() []Block {
	var chain []Block
	parent := b.Parent()
	for parent != nil {
		chain = append(chain, parent)
		parent = parent.Parent()
	}
	// Reverse the chain
	for i, j := 0, len(chain)-1; i < j; i, j = i+1, j-1 {
		chain[i], chain[j] = chain[j], chain[i]
	}
	return chain
}

// ParentText returns the text of the parent chain of the block
func (b *BaseBlock) ParentText() string {
	parentChain := b.ParentChain()
	var headerTexts, paraTexts []string
	for _, p := range parentChain {
		if p.Tag() == "header" {
			headerTexts = append(headerTexts, p.ToText(false, false))
		} else if p.Tag() == "list_item" || p.Tag() == "para" {
			paraTexts = append(paraTexts, p.ToText(false, false))
		}
	}
	text := ""
	if len(headerTexts) > 0 {
		text += strings.Join(headerTexts, " > ")
	}
	if len(paraTexts) > 0 {
		text += "\n" + strings.Join(paraTexts, "\n")
	}
	return text
}

// ToContextText returns the text of the block with section information
func (b *BaseBlock) ToContextText(includeSectionInfo bool) string {
	text := ""
	if includeSectionInfo {
		text += b.ParentText() + "\n"
	}
	if b.Tag() == "list_item" || b.Tag() == "para" || b.Tag() == "table" {
		text += b.ToText(true, true)
	} else {
		text += b.ToText(false, false)
	}
	return text
}

// IterChildren iterates over all the children of the block and calls the nodeVisitor function on each child
func (b *BaseBlock) IterChildren(node Block, level int, nodeVisitor func(Block)) {
	for _, child := range node.Children() {
		nodeVisitor(child)
		if child.Tag() != "list_item" && child.Tag() != "para" && child.Tag() != "table" {
			b.IterChildren(child, level+1, nodeVisitor)
		}
	}
}

// Paragraphs returns all the paragraphs in the block
func (b *BaseBlock) Paragraphs() []Block {
	var paragraphs []Block
	paraCollector := func(node Block) {
		if node.Tag() == "para" {
			paragraphs = append(paragraphs, node)
		}
	}
	b.IterChildren(b, 0, paraCollector)
	return paragraphs
}

// Chunks returns all the chunks in the block
func (b *BaseBlock) Chunks() []Block {
	var chunks []Block
	chunkCollector := func(node Block) {
		if node.Tag() == "para" || node.Tag() == "list_item" || node.Tag() == "table" {
			chunks = append(chunks, node)
		}
	}
	b.IterChildren(b, 0, chunkCollector)
	return chunks
}

// Tables returns all the tables in the block
func (b *BaseBlock) Tables() []Block {
	var tables []Block
	tableCollector := func(node Block) {
		if node.Tag() == "table" {
			tables = append(tables, node)
		}
	}
	b.IterChildren(b, 0, tableCollector)
	return tables
}

// Sections returns all the sections in the block
func (b *BaseBlock) Sections() []Block {
	var sections []Block
	sectionCollector := func(node Block) {
		if node.Tag() == "header" {
			sections = append(sections, node)
		}
	}
	b.IterChildren(b, 0, sectionCollector)
	return sections
}

// Paragraph represents a paragraph block
type Paragraph struct {
	*BaseBlock
}

// NewParagraph creates a new Paragraph instance
func NewParagraph(paraJSON map[string]interface{}) *Paragraph {
	return &Paragraph{
		BaseBlock: NewBaseBlock(paraJSON),
	}
}

// ToText converts the paragraph to text
func (p *Paragraph) ToText(includeChildren, recurse bool) string {
	paraText := strings.Join(p.Sentences(), "\n")
	if includeChildren {
		for _, child := range p.Children() {
			if recurse {
				paraText += "\n" + child.ToText(true, true)
			} else {
				paraText += "\n" + child.ToText(false, false)
			}
		}
	}
	return paraText
}

// ToHTML converts the paragraph to HTML
func (p *Paragraph) ToHTML(includeChildren, recurse bool) string {
	htmlStr := "<p>" + strings.Join(p.Sentences(), "\n")
	if includeChildren && len(p.Children()) > 0 {
		htmlStr += "<ul>"
		for _, child := range p.Children() {
			if recurse {
				htmlStr += child.ToHTML(true, true)
			} else {
				htmlStr += child.ToHTML(false, false)
			}
		}
		htmlStr += "</ul>"
	}
	htmlStr += "</p>"
	return htmlStr
}

// Section represents a section block
type Section struct {
	*BaseBlock
	Title string
}

// NewSection creates a new Section instance
func NewSection(sectionJSON map[string]interface{}) *Section {
	baseBlock := NewBaseBlock(sectionJSON)
	return &Section{
		BaseBlock: baseBlock,
		Title:     strings.Join(baseBlock.Sentences(), "\n"),
	}
}

// ToText converts the section to text
func (s *Section) ToText(includeChildren, recurse bool) string {
	text := s.Title
	if includeChildren {
		for _, child := range s.Children() {
			if recurse {
				text += "\n" + child.ToText(true, true)
			} else {
				text += "\n" + child.ToText(false, false)
			}
		}
	}
	return text
}

// ToHTML converts the section to HTML
func (s *Section) ToHTML(includeChildren, recurse bool) string {
	htmlStr := fmt.Sprintf("<h%d>%s</h%d>", s.Level()+1, s.Title, s.Level()+1)
	if includeChildren {
		for _, child := range s.Children() {
			if recurse {
				htmlStr += child.ToHTML(true, true)
			} else {
				htmlStr += child.ToHTML(false, false)
			}
		}
	}
	return htmlStr
}

// ListItem represents a list item block
type ListItem struct {
	*BaseBlock
}

// NewListItem creates a new ListItem instance
func NewListItem(listJSON map[string]interface{}) *ListItem {
	return &ListItem{
		BaseBlock: NewBaseBlock(listJSON),
	}
}

// ToText converts the list item to text
func (li *ListItem) ToText(includeChildren, recurse bool) string {
	text := strings.Join(li.Sentences(), "\n")
	if includeChildren {
		for _, child := range li.Children() {
			if recurse {
				text += "\n" + child.ToText(true, true)
			} else {
				text += "\n" + child.ToText(false, false)
			}
		}
	}
	return text
}

// ToHTML converts the list item to HTML
func (li *ListItem) ToHTML(includeChildren, recurse bool) string {
	htmlStr := "<li>" + strings.Join(li.Sentences(), "\n")
	if includeChildren && len(li.Children()) > 0 {
		htmlStr += "<ul>"
		for _, child := range li.Children() {
			if recurse {
				htmlStr += child.ToHTML(true, true)
			} else {
				htmlStr += child.ToHTML(false, false)
			}
		}
		htmlStr += "</ul>"
	}
	htmlStr += "</li>"
	return htmlStr
}

// TableCell represents a table cell block
type TableCell struct {
	*BaseBlock
	ColSpan   int
	CellValue interface{}
	CellNode  Block
}

// NewTableCell creates a new TableCell instance
func NewTableCell(cellJSON map[string]interface{}) *TableCell {
	baseBlock := NewBaseBlock(cellJSON)
	colSpan, _ := cellJSON["col_span"].(int)
	cellValue := cellJSON["cell_value"]
	var cellNode Block
	if _, ok := cellValue.(string); !ok {
		cellNode = NewParagraph(cellValue.(map[string]interface{}))
	}
	return &TableCell{
		BaseBlock: baseBlock,
		ColSpan:   colSpan,
		CellValue: cellValue,
		CellNode:  cellNode,
	}
}

// ToText returns the cell value as text
func (tc *TableCell) ToText() string {
	cellText := tc.CellValue.(string)
	if tc.CellNode != nil {
		cellText = tc.CellNode.ToText(false, false)
	}
	return cellText
}

// ToHTML returns the cell value as HTML
func (tc *TableCell) ToHTML() string {
	cellHTML := tc.CellValue.(string)
	if tc.CellNode != nil {
		cellHTML = tc.CellNode.ToHTML(false, false)
	}
	htmlStr := fmt.Sprintf("<td>%s</td>", cellHTML)
	if tc.ColSpan > 1 {
		htmlStr = fmt.Sprintf("<td colSpan=%d>%s</td>", tc.ColSpan, cellHTML)
	}
	return htmlStr
}

// TableRow represents a table row block
type TableRow struct {
	Cells []*TableCell
}

// NewTableRow creates a new TableRow instance
func NewTableRow(rowJSON map[string]interface{}) *TableRow {
	row := &TableRow{
		Cells: []*TableCell{},
	}
	if rowJSON["type"] == "full_row" {
		cell := NewTableCell(rowJSON)
		row.Cells = append(row.Cells, cell)
	} else {
		for _, cellJSON := range rowJSON["cells"].([]interface{}) {
			cell := NewTableCell(cellJSON.(map[string]interface{}))
			row.Cells = append(row.Cells, cell)
		}
	}
	return row
}

// ToText returns the text of the row with text from all the cells in the row delimited by '|'
func (tr *TableRow) ToText(includeChildren, recurse bool) string {
	var cellTexts []string
	for _, cell := range tr.Cells {
		cellTexts = append(cellTexts, cell.ToText())
	}
	return " | " + strings.Join(cellTexts, " | ")
}

// ToHTML returns the HTML for a <tr> with HTML from all the cells in the row as <td>
func (tr *TableRow) ToHTML(includeChildren, recurse bool) string {
	htmlStr := "<tr>"
	for _, cell := range tr.Cells {
		htmlStr += cell.ToHTML()
	}
	htmlStr += "</tr>"
	return htmlStr
}

type TableHeader struct {
	*BaseBlock
}

func NewTableHeader(headerJSON map[string]interface{}) *TableHeader {
	return &TableHeader{
		BaseBlock: NewBaseBlock(headerJSON),
	}
}

func (th *TableHeader) ToText(includeChildren, recurse bool) string {
	cellText := ""
	for _, cell := range th.Cells() {
		cellText += " | " + cell.ToText()
	}
	cellText += "\n"
	for range th.Cells() {
		cellText += " | ---"
	}
	return cellText
}

func (th *TableHeader) ToHTML(includeChildren, recurse bool) string {
	htmlStr := "<th>"
	for _, cell := range th.Cells() {
		htmlStr += cell.ToHTML()
	}
	htmlStr += "</th>"
	return htmlStr
}

func (th *TableHeader) Cells() []*TableCell {
	var cells []*TableCell
	for _, cellJSON := range th.BlockJSON()["cells"].([]interface{}) {
		cell := NewTableCell(cellJSON.(map[string]interface{}))
		cells = append(cells, cell)
	}
	return cells
}

type Table struct {
	*BaseBlock
	Rows    []*TableRow
	Headers []*TableHeader
	Name    string
}

func NewTable(tableJSON map[string]interface{}, parent Block) *Table {
	baseBlock := NewBaseBlock(tableJSON)
	table := &Table{
		BaseBlock: baseBlock,
		Rows:      []*TableRow{},
		Headers:   []*TableHeader{},
		Name:      tableJSON["name"].(string),
	}
	if rowsData, ok := tableJSON["table_rows"].([]interface{}); ok {
		for _, rowData := range rowsData {
			rowJSON := rowData.(map[string]interface{})
			if rowJSON["type"] == "table_header" {
				row := NewTableHeader(rowJSON)
				table.Headers = append(table.Headers, row)
			} else {
				row := NewTableRow(rowJSON)
				table.Rows = append(table.Rows, row)
			}
		}
	}
	return table
}

func (t *Table) ToText(includeChildren, recurse bool) string {
	text := ""
	for _, header := range t.Headers {
		text += header.ToText(false, false) + "\n"
	}
	for _, row := range t.Rows {
		text += row.ToText(false, false) + "\n"
	}
	return text
}

func (t *Table) ToHTML(includeChildren, recurse bool) string {
	htmlStr := "<table>"
	for _, header := range t.Headers {
		htmlStr += header.ToHTML(false, false)
	}
	for _, row := range t.Rows {
		htmlStr += row.ToHTML(false, false)
	}
	htmlStr += "</table>"
	return htmlStr
}

type LayoutReader struct{}

func (lr *LayoutReader) Debug(pdfRoot Block) {
	var iterChildren func(node Block, level int)
	iterChildren = func(node Block, level int) {
		indent := strings.Repeat("-", level)
		fmt.Printf("%s %s (%d) %s\n", indent, node.Tag(), len(node.Children()), node.ToText(false, false))
		for _, child := range node.Children() {
			iterChildren(child, level+1)
		}
	}
	iterChildren(pdfRoot, 0)
}

func (lr *LayoutReader) Read(blocksJSON []interface{}) Block {
	root := &BaseBlock{}
	parent := Block(root)
	parentStack := []Block{root}
	var prevNode Block = root
	listStack := []Block{}

	for _, blockData := range blocksJSON {
		blockJSON := blockData.(map[string]interface{})
		tag := blockJSON["tag"].(string)

		if tag != "list_item" && len(listStack) > 0 {
			listStack = []Block{}
		}

		var node Block
		switch tag {
		case "para":
			node = NewParagraph(blockJSON)
			parent.AddChild(node)
		case "table":
			node = NewTable(blockJSON, prevNode)
			parent.AddChild(node)
		case "list_item":
			node = NewListItem(blockJSON)
			if prevNode.Tag() == "para" && prevNode.Level() == node.Level() {
				listStack = append(listStack, prevNode)
			} else if prevNode.Tag() == "list_item" {
				if node.Level() > prevNode.Level() {
					listStack = append(listStack, prevNode)
				} else if node.Level() < prevNode.Level() {
					for len(listStack) > 0 && listStack[len(listStack)-1].Level() > node.Level() {
						listStack = listStack[:len(listStack)-1]
					}
				}
			}
			if len(listStack) > 0 {
				listStack[len(listStack)-1].AddChild(node)
			} else {
				parent.AddChild(node)
			}
		case "header":
			node = NewSection(blockJSON)
			if node.Level() > parent.Level() {
				parentStack = append(parentStack, node)
				parent.AddChild(node)
			} else {
				for len(parentStack) > 1 && parentStack[len(parentStack)-1].Level() > node.Level() {
					parentStack = parentStack[:len(parentStack)-1]
				}
				parentStack[len(parentStack)-1].AddChild(node)
				parentStack = append(parentStack, node)
			}
			parent = node
		}
		prevNode = node
	}

	return root
}

type Document struct {
	Reader   *LayoutReader
	RootNode Block
	JSON     []interface{}
}

func (d *Document) Chunks() []Block {
	return d.RootNode.Chunks()
}

func (d *Document) Tables() []Block {
	return d.RootNode.Tables()
}

func (d *Document) Sections() []Block {
	return d.RootNode.Sections()
}

func (d *Document) ToText() string {
	text := ""
	for _, section := range d.Sections() {
		text += section.ToText(true, true) + "\n"
	}
	return text
}

func (d *Document) ToHTML() string {
	htmlStr := "<html>"
	for _, section := range d.Sections() {
		htmlStr += section.ToHTML(true, true)
	}
	htmlStr += "</html>"
	return htmlStr
}

func NewDocument(blocksJSON []interface{}) *Document {
	reader := &LayoutReader{}
	rootNode := reader.Read(blocksJSON)
	return &Document{
		Reader:   reader,
		RootNode: rootNode,
		JSON:     blocksJSON,
	}
}
