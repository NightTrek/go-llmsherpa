# Go LLM Sherpa

LLM Sherpa provides strategic APIs to accelerate large language model (LLM) use cases. We didnt like using the python version so we decided to convert it into go. This version should work well with the existing nlm-ingestor backend. See [https://github.com/nlmatics/nlm-ingestor](https://github.com/nlmatics/nlm-ingestor)


## LayoutPDFReader

Most PDF to text parsers do not provide layout information. Often times, even the sentences are split with arbritrary CR/LFs making it very difficult to find paragraph boundaries. This poses various challenges in chunking and adding long running contextual information such as section header to the passages while indexing/vectorizing PDFs for LLM applications such as retrieval augmented generation (RAG). 

LayoutPDFReader solves this problem by parsing PDFs along with hierarchical layout information such as:

1. Sections and subsections along with their levels.
2. Paragraphs - combines lines.
3. Links between sections and paragraphs.
4. Tables along with the section the tables are found in.
5. Lists and nested lists.
6. Join content spread across pages.
7. Removal of repeating headers and footers.
8. Watermark removal.

With LayoutPDFReader, developers can find optimal chunks of text to vectorize, and a solution for limited context window sizes of LLMs. 

Big thanks to the orginal python version which can be found here: [https://github.com/nlmatics/llmsherpa](https://github.com/nlmatics/llmsherpa)


# Usage 

To use the LayoutPDFReader:

1. Create a new instance of `LayoutPDFReader` by providing the URL of the PDF parser API:

```go
reader := NewLayoutPDFReader("https://example.com/parser-api")
```

2. Read a PDF file by providing the path or URL to the PDF:

```go
doc, err := reader.ReadPDF("path/to/file.pdf", nil)
if err != nil {
    // Handle error
}
```

Alternatively, you can provide the PDF contents directly:

```go
pdfData, err := os.ReadFile("path/to/file.pdf")
if err != nil {
    // Handle error
}
doc, err := reader.ReadPDF("file.pdf", pdfData)
if err != nil {
    // Handle error
}
```

3. The `ReadPDF` method returns a `Document` struct containing the parsed PDF information. You can access the various layout elements such as sections, paragraphs, tables, and lists from the `Document` struct.

Note: Make sure to provide a valid URL for the PDF parser API when creating the `LayoutPDFReader` instance.
