package chipper

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"testing"
)

// go test ./chipper -c -o ./bin/chippertest.test && ./bin/chippertest.test -test.run=TestChipper -test.v

// const testPDFURL = "https://raw.githubusercontent.com/run-llama/llama_index/main/docs/docs/examples/data/10q/uber_10q_march_2022.pdf"

func ReadPDFTest() (*Document, error) {
	var err error

	var response map[string]interface{}

	// pull response json from response.json file
	parserResponse, err := os.ReadFile("../response.json")
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(parserResponse, &response)
	if err != nil {
		return nil, err
	}

	blocks := response["return_dict"].(map[string]interface{})["result"].(map[string]interface{})["blocks"].([]interface{})
	return NewDocument(blocks), nil
}

func TestChipper(t *testing.T) {
	t.Run("TestReadPDFFromURL", func(t *testing.T) {

		// Run the Python test and capture its output
		cmd := exec.Command("python3.10", "../python_test/pdfreader_test.py")
		var outBuf, errBuf bytes.Buffer
		cmd.Stdout = &outBuf
		cmd.Stderr = &errBuf
		err := cmd.Run()
		if err != nil {
			t.Fatalf("Failed to run Python test: %v\nError: %s", err, errBuf.String())
		}
		pythonOutput := outBuf.String()

		// Initialize the LayoutPDFReader with a dummy parser API URL.
		// reader := NewLayoutPDFReader("testlink")

		// Call the ReadPDF method with the test PDF URL. Pass nil for the contents to indicate URL download.
		doc, err := ReadPDFTest()
		if err != nil {
			t.Fatalf("ReadPDF failed: %v", err)
		}

		// doc is a Document struct
		if doc == nil || len(doc.Chunks()) == 0 {
			t.Fatal("ReadPDF returned an invalid or empty document.")
		}

		sections := make([]string, 0)

		// iterate through doc.Sections()
		for _, section := range doc.Sections() {
			// t.Logf("Section tag: %s", section.Tag())
			sectionText := section.ToText(true, false)
			t.Logf("Section: %s", sectionText)
			// sectionHTML := section.ToHTML(false, true)
			// t.Logf("Section HTML: %s", sectionHTML)

			var iterChildren func(node BlockInterface, level int)
			iterChildren = func(node BlockInterface, level int) {
				switch nodeBlock := node.(type) {
				case *Paragraph:
					t.Logf("%s-Paragraph", strings.Repeat("  ", level))
					t.Logf("# Sentences: %d", len(nodeBlock.Sentences))
					t.Logf("PARAGRAPH SPOTTED: %s", strings.Join(nodeBlock.Sentences, "\n\n SENTENCE \n\n"))
				case *Block:
					for _, child := range nodeBlock.Children {
						switch childBlock := child.(type) {
						case *Section:
							t.Logf("%s-Section", strings.Repeat("  ", level))
							t.Logf("# Sentences: %d", len(childBlock.Sentences))
						}
						iterChildren(child, level+1)
					}
				case *Section:
					for _, child := range nodeBlock.Children {
						switch childBlock := child.(type) {
						case *Section:
							t.Logf("%s-Section", strings.Repeat("  ", level))
							t.Logf("# Sentences: %d", len(childBlock.Sentences))
						}
						iterChildren(child, level+1)
					}
					sections = append(sections, nodeBlock.ToContextText(true))
				case *ListItem:
					t.Logf("%s-Paragraph", strings.Repeat("  ", level))
					t.Logf("# Sentences: %d", len(nodeBlock.Sentences))
					t.Logf("LIST ITEM SPOTTED: %s", strings.Join(nodeBlock.Sentences, "\n\n LIST SENTENCE \n\n"))
					for _, child := range nodeBlock.Children {
						switch childBlock := child.(type) {
						case *Section:
							t.Logf("%s-Section", strings.Repeat("  ", level))
							t.Logf("# Sentences: %d", len(childBlock.Sentences))
						}
						iterChildren(child, level+1)
					}
				// Handle other block types as needed
				default:
					t.Logf("Unsupported block type: %T", node)
				}
			}

			iterChildren(section, 1)

		}

		sentences := make([]string, 0)
		// iterate through doc.Chunks()
		for _, chunk := range doc.Chunks() {
			// t.Logf("Chunk tag: %s", chunk.Tag())
			// chunkText := chunk.ToText(false, false)
			// t.Logf("Chunk: %s", chunkText)
			switch nodeBlock := chunk.(type) {
			case *Block:
				t.Logf("Number of sentences in block: %d", len(nodeBlock.Sentences))
				sentences = append(sentences, nodeBlock.Sentences...)
			case *Section:
				t.Logf("Number of sentences in section: %d", len(nodeBlock.Sentences))
				sentences = append(sentences, nodeBlock.Sentences...)
			case *ListItem:
				t.Logf("Number of sentences in listItem: %d", len(nodeBlock.Sentences))
				sentences = append(sentences, nodeBlock.Sentences...)
			case *Paragraph:
				t.Logf("Number of sentences in paragraph: %d", len(nodeBlock.Sentences))
				sentences = append(sentences, nodeBlock.Sentences...)
			}
			// // chunkHTML := chunk.ToHTML(false, false)
			// // t.Logf("Chunk HTML: %s", chunkHTML)
			// paragraphs := chunk.Paragraphs()
			// t.Logf("Number of paragraphs in chunk: %d", len(paragraphs))
			// for _, paragraph := range paragraphs {
			// 	t.Logf("Paragraph: %s", paragraph.ToText(false, false))
			// 	sentences := paragraph.(*Paragraph).Sentences
			// 	t.Logf("Number of sentences in paragraph: %d", len(sentences))

			// }

			//iterate through all the chunk children recursively and print out the lengths of the sentence array for every paragraph
		}

		// iterate and print each sentence with \n between
		for _, sentence := range sentences {
			t.Logf("\n\n%s\n\n", sentence)
		}

		t.Logf("number of sections: %d", len(doc.Sections()))
		t.Logf("number of chunks: %d", len(doc.Chunks()))
		t.Logf("number of tables: %d", len(doc.Tables()))
		t.Logf("number of sentences: %d", len(sentences))

		t.Logf("table example:\n\n%s\n", doc.Tables()[4].ToHTML(true, true))

		// Compare the results with the Python output
		goOutput := fmt.Sprintf("Number of sections: %d\nNumber of chunks: %d\nNumber of tables: %d\nNumber of sentences: %d",
			len(doc.Sections()), len(doc.Chunks()), len(doc.Tables()), len(sentences))

		if goOutput != pythonOutput {
			t.Errorf("Output mismatch:\nPython:\n%s\nGo:\n%s", pythonOutput, goOutput)
		}
		// Further checks can be added here based on the specifics of your implementation and what constitutes a successful read operation.
		// Examples might include checking specific text blocks or document properties to ensure the parsing was successful.
	})
}
