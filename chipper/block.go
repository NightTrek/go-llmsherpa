package chipper

import (
	"fmt"
	"strings"
)

type BlockInterface interface {
	AddChild(node BlockInterface)
	ToHTML(includeChildren, recurse bool) string
	ToText(includeChildren, recurse bool) string
	ParentChain() []BlockInterface
	ParentText() string
	ToContextText(includeSectionInfo bool) string
	IterChildren(node BlockInterface, level int, nodeVisitor func(BlockInterface))
	Paragraphs() []BlockInterface
	Chunks() []BlockInterface
	Tables() []BlockInterface
	Sections() []BlockInterface
}

type Block struct {
	Tag       string    `json:"tag"`
	Level     int       `json:"level"`
	PageIdx   int       `json:"page_idx"`
	BlockIdx  int       `json:"block_idx"`
	Top       float64   `json:"top"`
	Left      float64   `json:"left"`
	Bbox      []float64 `json:"bbox"`
	Sentences []string  `json:"sentences"`
	Children  []BlockInterface
	Parent    BlockInterface
	BlockJSON map[string]interface{}
}

func NewBlock(blockJSON map[string]interface{}) *Block {
	block := &Block{
		BlockJSON: blockJSON,
	}

	if tag, ok := blockJSON["tag"].(string); ok {
		block.Tag = tag
	}
	if level, ok := blockJSON["level"].(float64); ok {
		block.Level = int(level)
	} else {
		block.Level = -1
	}
	if pageIdx, ok := blockJSON["page_idx"].(float64); ok {
		block.PageIdx = int(pageIdx)
	} else {
		block.PageIdx = -1
	}
	if blockIdx, ok := blockJSON["block_idx"].(float64); ok {
		block.BlockIdx = int(blockIdx)
	} else {
		block.BlockIdx = -1
	}
	if top, ok := blockJSON["top"].(float64); ok {
		block.Top = top
	} else {
		block.Top = -1
	}
	if left, ok := blockJSON["left"].(float64); ok {
		block.Left = left
	} else {
		block.Left = -1
	}
	if bbox, ok := blockJSON["bbox"].([]interface{}); ok {
		for _, val := range bbox {
			if f, ok := val.(float64); ok {
				block.Bbox = append(block.Bbox, f)
			}
		}
	}
	if sentences, ok := blockJSON["sentences"].([]interface{}); ok {
		for _, val := range sentences {
			if s, ok := val.(string); ok {
				block.Sentences = append(block.Sentences, s)
			}
		}
	}

	return block
}

func (b *Block) AddChild(node BlockInterface) {
	b.Children = append(b.Children, node)
	switch childNode := node.(type) {
	case *Block:
		childNode.Parent = b
	case *Paragraph:
		childNode.Parent = b
	case *Section:
		childNode.Parent = b
	case *ListItem:
		childNode.Parent = b
	case *Table:
		childNode.Parent = b
	}
}

func (b *Block) ToHTML(includeChildren, recurse bool) string {
	// Implement the ToHTML method for the Block struct
	return ""
}

func (b *Block) ToText(includeChildren, recurse bool) string {
	// Implement the ToText method for the Block struct
	return ""
}

func (b *Block) ParentChain() []BlockInterface {
	var chain []BlockInterface
	parent := b.Parent
	for parent != nil {
		chain = append([]BlockInterface{parent}, chain...)
		switch parentBlock := parent.(type) {
		case *Block:
			parent = parentBlock.Parent
		case *Paragraph:
			parent = parentBlock.Parent
		case *Section:
			parent = parentBlock.Parent
		case *ListItem:
			parent = parentBlock.Parent
		case *Table:
			parent = parentBlock.Parent
		default:
			parent = nil
		}
	}
	return chain
}

func (b *Block) ParentText() string {
	parentChain := b.ParentChain()
	var headerTexts, paraTexts []string
	for _, p := range parentChain {
		if p.(*Block).Tag == "header" {
			headerTexts = append(headerTexts, p.ToText(false, false))
		} else if p.(*Block).Tag == "list_item" || p.(*Block).Tag == "para" {
			paraTexts = append(paraTexts, p.ToText(false, false))
		}
	}
	text := ""
	if len(headerTexts) > 0 {
		text += " > " + strings.Join(headerTexts, " > ")
	}
	if len(paraTexts) > 0 {
		text += "\n" + strings.Join(paraTexts, "\n")
	}
	return text
}

func (b *Block) ToContextText(includeSectionInfo bool) string {
	text := ""
	if includeSectionInfo {
		text += b.ParentText() + "\n"
	}
	if b.Tag == "list_item" || b.Tag == "para" || b.Tag == "table" {
		text += b.ToText(true, true)
	} else {
		text += b.ToText(false, false)
	}
	return text
}

func (b *Block) IterChildren(node BlockInterface, level int, nodeVisitor func(BlockInterface)) {
	switch nodeBlock := node.(type) {
	case *Block:
		for _, child := range nodeBlock.Children {
			nodeVisitor(child)
			b.IterChildren(child, level+1, nodeVisitor)
		}
	case *Paragraph:
		for _, child := range nodeBlock.Children {
			nodeVisitor(child)
			b.IterChildren(child, level+1, nodeVisitor)
		}
	case *Section:
		for _, child := range nodeBlock.Children {
			nodeVisitor(child)
			b.IterChildren(child, level+1, nodeVisitor)
		}
	case *ListItem:
		for _, child := range nodeBlock.Children {
			nodeVisitor(child)
			b.IterChildren(child, level+1, nodeVisitor)
		}
	case *Table:
		for _, child := range nodeBlock.Children {
			nodeVisitor(child)
			b.IterChildren(child, level+1, nodeVisitor)
		}
	}
}

func (b *Block) Paragraphs() []BlockInterface {
	var paragraphs []BlockInterface
	paraCollector := func(node BlockInterface) {
		switch nodeBlock := node.(type) {
		case *Paragraph:
			paragraphs = append(paragraphs, nodeBlock)
		}
	}
	b.IterChildren(b, 0, paraCollector)
	return paragraphs
}

func (b *Block) Chunks() []BlockInterface {
	var chunks []BlockInterface
	chunkCollector := func(node BlockInterface) {
		switch nodeBlock := node.(type) {
		case *Paragraph, *ListItem, *Table:
			chunks = append(chunks, nodeBlock)
		}
	}
	b.IterChildren(b, 0, chunkCollector)
	return chunks
}

func (b *Block) Tables() []BlockInterface {
	var tables []BlockInterface
	tableCollector := func(node BlockInterface) {
		switch nodeBlock := node.(type) {
		case *Table:
			tables = append(tables, nodeBlock)
		}
	}
	b.IterChildren(b, 0, tableCollector)
	return tables
}

func (b *Block) Sections() []BlockInterface {
	var sections []BlockInterface
	sectionCollector := func(node BlockInterface) {
		switch nodeBlock := node.(type) {
		case *Section:
			sections = append(sections, nodeBlock)
		}
	}
	b.IterChildren(b, 0, sectionCollector)
	return sections
}

type Paragraph struct {
	*Block
}

func NewParagraph(paraJSON map[string]interface{}) *Paragraph {
	return &Paragraph{
		Block: NewBlock(paraJSON),
	}
}

func (p *Paragraph) ToText(includeChildren, recurse bool) string {
	paraText := strings.Join(p.Sentences, "\n")
	if includeChildren {
		for _, child := range p.Children {
			paraText += "\n" + child.ToText(recurse, recurse)
		}
	}
	return paraText
}

func (p *Paragraph) ToHTML(includeChildren, recurse bool) string {
	htmlStr := "<p>"
	htmlStr += strings.Join(p.Sentences, "<br>")
	if includeChildren && len(p.Children) > 0 {
		htmlStr += "<ul>"
		for _, child := range p.Children {
			htmlStr += child.ToHTML(recurse, recurse)
		}
		htmlStr += "</ul>"
	}
	htmlStr += "</p>"
	return htmlStr
}

type Section struct {
	*Block
	Title string `json:"title"`
}

func NewSection(sectionJSON map[string]interface{}) *Section {
	section := &Section{
		Block: NewBlock(sectionJSON),
	}
	section.Title = strings.Join(section.Sentences, "\n")
	return section
}

func (s *Section) ToText(includeChildren, recurse bool) string {
	text := s.Title
	if includeChildren {
		for _, child := range s.Children {
			text += "\n" + child.ToText(recurse, recurse)
		}
	}
	return text
}

func (s *Section) ToHTML(includeChildren, recurse bool) string {
	htmlStr := fmt.Sprintf("<h%d>", s.Level+1)
	htmlStr += s.Title
	htmlStr += fmt.Sprintf("</h%d>", s.Level+1)
	if includeChildren {
		for _, child := range s.Children {
			htmlStr += child.ToHTML(recurse, recurse)
		}
	}
	return htmlStr
}

type ListItem struct {
	*Block
}

func NewListItem(listJSON map[string]interface{}) *ListItem {
	return &ListItem{
		Block: NewBlock(listJSON),
	}
}

func (li *ListItem) ToText(includeChildren, recurse bool) string {
	text := strings.Join(li.Sentences, "\n")
	if includeChildren {
		for _, child := range li.Children {
			text += "\n" + child.ToText(recurse, recurse)
		}
	}
	return text
}

func (li *ListItem) ToHTML(includeChildren, recurse bool) string {
	htmlStr := "<li>"
	htmlStr += strings.Join(li.Sentences, "<br>")
	if includeChildren && len(li.Children) > 0 {
		htmlStr += "<ul>"
		for _, child := range li.Children {
			htmlStr += child.ToHTML(recurse, recurse)
		}
		htmlStr += "</ul>"
	}
	htmlStr += "</li>"
	return htmlStr
}

type TableCell struct {
	*Block
	ColSpan   int         `json:"col_span"`
	CellValue interface{} `json:"cell_value"`
	CellNode  *Paragraph
}

func NewTableCell(cellJSON map[string]interface{}) *TableCell {
	cell := &TableCell{
		Block: NewBlock(cellJSON),
	}

	if colSpan, ok := cellJSON["col_span"].(float64); ok {
		cell.ColSpan = int(colSpan)
	} else {
		cell.ColSpan = 1
	}

	cell.CellValue = cellJSON["cell_value"]

	if _, ok := cell.CellValue.(string); !ok {
		cell.CellNode = NewParagraph(cell.CellValue.(map[string]interface{}))
	}

	return cell
}

func (tc *TableCell) ToText() string {
	cellText := ""
	switch value := tc.CellValue.(type) {
	case string:
		cellText = value
	default:
		if tc.CellNode != nil {
			cellText = tc.CellNode.ToText(false, false)
		}
	}
	return cellText
}

func (tc *TableCell) ToHTML() string {
	cellHTML := ""
	switch value := tc.CellValue.(type) {
	case string:
		cellHTML = value
	default:
		if tc.CellNode != nil {
			cellHTML = tc.CellNode.ToHTML(false, false)
		}
	}

	htmlStr := "<td"
	if tc.ColSpan > 1 {
		htmlStr += fmt.Sprintf(" colspan=\"%d\"", tc.ColSpan)
	}
	htmlStr += ">" + cellHTML + "</td>"

	return htmlStr
}

type TableRow struct {
	*Block
	Cells []*TableCell
}

func NewTableRow(rowJSON map[string]interface{}) *TableRow {
	row := &TableRow{
		Block: NewBlock(rowJSON),
		Cells: make([]*TableCell, 0),
	}

	if rowType, ok := rowJSON["type"].(string); ok && rowType == "full_row" {
		cell := NewTableCell(rowJSON)
		row.Cells = append(row.Cells, cell)
	} else {
		if cells, ok := rowJSON["cells"].([]interface{}); ok {
			for _, cellJSON := range cells {
				cell := NewTableCell(cellJSON.(map[string]interface{}))
				row.Cells = append(row.Cells, cell)
			}
		}
	}

	return row
}

func (tr *TableRow) ToText(includeChildren, recurse bool) string {
	cellText := ""
	for _, cell := range tr.Cells {
		cellText += " | " + cell.ToText()
	}
	return cellText
}

func (tr *TableRow) ToHTML(includeChildren, recurse bool) string {
	htmlStr := "<tr>"
	for _, cell := range tr.Cells {
		htmlStr += cell.ToHTML()
	}
	htmlStr += "</tr>"
	return htmlStr
}

type TableHeader struct {
	*Block
	Cells []*TableCell
}

func NewTableHeader(rowJSON map[string]interface{}) *TableHeader {
	header := &TableHeader{
		Block: NewBlock(rowJSON),
		Cells: make([]*TableCell, 0),
	}

	if cells, ok := rowJSON["cells"].([]interface{}); ok {
		for _, cellJSON := range cells {
			cell := NewTableCell(cellJSON.(map[string]interface{}))
			header.Cells = append(header.Cells, cell)
		}
	}

	return header
}

func (th *TableHeader) ToText(includeChildren, recurse bool) string {
	cellText := ""
	for _, cell := range th.Cells {
		cellText += " | " + cell.ToText()
	}
	cellText += "\n"
	for range th.Cells {
		cellText += " | ---"
	}
	return cellText
}

func (th *TableHeader) ToHTML(includeChildren, recurse bool) string {
	htmlStr := "<th>"
	for _, cell := range th.Cells {
		htmlStr += cell.ToHTML()
	}
	htmlStr += "</th>"
	return htmlStr
}

type Table struct {
	*Block
	Rows    []*TableRow
	Headers []*TableHeader
	Name    string `json:"name"`
}

func NewTable(tableJSON map[string]interface{}, parent BlockInterface) *Table {
	table := &Table{
		Block:   NewBlock(tableJSON),
		Rows:    make([]*TableRow, 0),
		Headers: make([]*TableHeader, 0),
		Name:    tableJSON["name"].(string),
	}

	if tableRows, ok := tableJSON["table_rows"].([]interface{}); ok {
		for _, rowJSON := range tableRows {
			rowData := rowJSON.(map[string]interface{})
			if rowType, ok := rowData["type"].(string); ok && rowType == "table_header" {
				header := NewTableHeader(rowData)
				table.Headers = append(table.Headers, header)
			} else {
				row := NewTableRow(rowData)
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
	return strings.TrimSpace(text)
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

func (lr *LayoutReader) Debug(pdfRoot BlockInterface) {
	var iterChildren func(node BlockInterface, level int)
	iterChildren = func(node BlockInterface, level int) {
		switch nodeBlock := node.(type) {
		case *Block:
			for _, child := range nodeBlock.Children {
				fmt.Printf("%s%s (%d) %s\n", strings.Repeat("-", level), child.(*Block).Tag, len(child.(*Block).Children), child.ToText(false, false))
				iterChildren(child, level+1)
			}
		case *Paragraph, *Section, *ListItem, *Table:
			for _, child := range nodeBlock.(*Block).Children {
				fmt.Printf("%s%s (%d) %s\n", strings.Repeat("-", level), child.(*Block).Tag, len(child.(*Block).Children), child.ToText(false, false))
				iterChildren(child, level+1)
			}
		}
	}
	iterChildren(pdfRoot, 0)
}

func (lr *LayoutReader) Read(blocksJSON []interface{}) BlockInterface {
	rootNode := &Block{}
	var parent BlockInterface = rootNode
	parentStack := []BlockInterface{rootNode}
	var prevNode BlockInterface = rootNode
	var listStack []BlockInterface

	for _, blockData := range blocksJSON {
		blockMap := blockData.(map[string]interface{})
		tag := blockMap["tag"].(string)

		var node BlockInterface
		switch tag {
		case "para":
			node = NewParagraph(blockMap)
		case "table":
			node = NewTable(blockMap, prevNode)
		case "list_item":
			node = NewListItem(blockMap)
		case "header":
			node = NewSection(blockMap)
		default:
			node = NewBlock(blockMap)
		}

		currentLevel := -1
		if level, ok := blockMap["level"].(float64); ok {
			currentLevel = int(level)
		}

		// Handling list items with hierarchy and sections
		if tag == "list_item" {
			// Check if the last node was a list item and manage list stack accordingly
			if len(listStack) > 0 {
				lastListItem := listStack[len(listStack)-1].(*ListItem)
				if currentLevel > lastListItem.Level {
					listStack = append(listStack, node)
				} else {
					// Pop from stack until a node with less or equal level is found
					for len(listStack) > 0 && listStack[len(listStack)-1].(*ListItem).Level >= currentLevel {
						listStack = listStack[:len(listStack)-1]
					}
				}
			}
			if len(listStack) > 0 {
				listStack[len(listStack)-1].AddChild(node)
			} else {
				parent.AddChild(node)
			}
			listStack = append(listStack, node)
		} else {
			// Handling sections with hierarchy
			if tag == "header" {
				// Push to parent stack or pop accordingly
				for len(parentStack) > 0 {
					topBlock, ok := parentStack[len(parentStack)-1].(*Block)
					if !ok {
						// Handle error: topBlock is not of type *Block, which should ideally never happen if the stack is managed correctly
						fmt.Println("Error: Non-Block type found in parentStack")
						break
					}
					if topBlock.Level < currentLevel {
						break
					}
					parentStack = parentStack[:len(parentStack)-1]
				}

				if len(parentStack) > 0 {
					parent = parentStack[len(parentStack)-1]
				} else {
					parent = rootNode
				}
				parent.AddChild(node)
				parentStack = append(parentStack, node)
				parent = node // Set new parent to the current node since it's a header and can have children
			} else {
				parent.AddChild(node) // Add the current node to the children of the current parent
			}
		}

		prevNode = node
	}

	return rootNode
}

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
