package chipper

import (
	"encoding/json"
	"os"
	"testing"
)

// go test ./chipper -c -o ./bin/chippertest.test && ./bin/chippertest.test -test.run=TestChipper -test.v

// const testPDFURL = "https://raw.githubusercontent.com/run-llama/llama_index/main/docs/docs/examples/data/10q/uber_10q_march_2022.pdf"

func (r *LayoutPDFReader) ReadPDFTest() (*Document, error) {
	var err error

	var response map[string]interface{}

	// pull response json from response.json file
	parserResponse, err := os.ReadFile("response.json")
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
		// Initialize the LayoutPDFReader with a dummy parser API URL.
		reader := NewLayoutPDFReader("testlink")

		// Call the ReadPDF method with the test PDF URL. Pass nil for the contents to indicate URL download.
		doc, err := reader.ReadPDFTest()
		if err != nil {
			t.Fatalf("ReadPDF failed: %v", err)
		}

		// doc is a Document struct
		if doc == nil || len(doc.Chunks()) == 0 {
			t.Fatal("ReadPDF returned an invalid or empty document.")
		}

		// iterate through doc.Sections()
		for _, section := range doc.Sections() {
			// t.Logf("Section tag: %s", section.Tag())
			sectionText := section.ToText(true, false)
			t.Logf("Section: %s", sectionText)
			// sectionHTML := section.ToHTML(false, true)
			// t.Logf("Section HTML: %s", sectionHTML)
		}

		// iterate through doc.Chunks()
		for _, chunk := range doc.Chunks() {
			// t.Logf("Chunk tag: %s", chunk.Tag())
			chunkText := chunk.ToText(false, false)
			t.Logf("Chunk: %s", chunkText)
			// chunkHTML := chunk.ToHTML(false, false)
			// t.Logf("Chunk HTML: %s", chunkHTML)
			paragraphs := chunk.Paragraphs()
			t.Logf("Number of paragraphs in chunk: %d", len(paragraphs))
			for _, paragraph := range paragraphs {
				t.Logf("Paragraph: %s", paragraph.ToText(false, false))
			}

			// sentences := chunk.Sentences()
			// t.Logf("Number of sentences in chunk: %d", len(sentences))
			// for _, sentence := range sentences {
			// 	t.Logf("Sentence: %s", sentence)
			// }

		}

		t.Logf("number of sections: %d", len(doc.Sections()))
		t.Logf("number of chunks: %d", len(doc.Chunks()))
		t.Logf("number of tables: %d", len(doc.Tables()))

		// Further checks can be added here based on the specifics of your implementation and what constitutes a successful read operation.
		// Examples might include checking specific text blocks or document properties to ensure the parsing was successful.
	})
}
