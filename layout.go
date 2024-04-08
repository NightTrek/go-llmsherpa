package chipper

import (
	"fmt"
	"strings"
)

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
		if blockMap["tag"] != "list_item" && len(listStack) > 0 {
			listStack = []BlockInterface{}
		}

		var node BlockInterface
		switch blockMap["tag"] {
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

		switch nodeBlock := node.(type) {
		case *ListItem:
			if len(listStack) > 0 {
				lastListItem := listStack[len(listStack)-1].(*ListItem)
				if nodeBlock.Level > lastListItem.Level {
					listStack = append(listStack, prevNode)
				} else if nodeBlock.Level < lastListItem.Level {
					for len(listStack) > 0 {
						lastListItem = listStack[len(listStack)-1].(*ListItem)
						if lastListItem.Level >= nodeBlock.Level {
							listStack = listStack[:len(listStack)-1]
						} else {
							break
						}
					}
				}
			}
			if len(listStack) > 0 {
				listStack[len(listStack)-1].AddChild(node)
			} else {
				parent.AddChild(node)
			}
		case *Section:
			switch parentBlock := parent.(type) {
			case *Block:
				if nodeBlock.Level > parentBlock.Level {
					parentStack = append(parentStack, node)
					parent.AddChild(node)
				} else {
					for len(parentStack) > 1 {
						switch parentStackBlock := parentStack[len(parentStack)-1].(type) {
						case *Block:
							if parentStackBlock.Level > nodeBlock.Level {
								parentStack = parentStack[:len(parentStack)-1]
							}
						}
					}
					parentStack[len(parentStack)-1].AddChild(node)
					parentStack = append(parentStack, node)
				}
				parent = node
			}
		default:
			parent.AddChild(node)
		}
		prevNode = node
	}

	return rootNode
}
