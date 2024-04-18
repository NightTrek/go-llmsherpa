import json
from llmsherpa.readers import  Document, Paragraph, Block, Section, ListItem

def log(debug=False, *args):
    if debug:
        print(*args)
    

def test_layout_pdf_reader(debug=False):
    # # Initialize the LayoutPDFReader with a dummy parser API URL.
    # reader = LayoutPDFReader("testlink")

    # Read the response JSON from the response.json file.
    with open("../response.json", "r") as file:
        response_json = json.load(file)

    blocks = response_json["return_dict"]["result"]["blocks"]
    doc = Document(blocks)

    # Test the ReadPDF method
    assert doc is not None
    assert len(doc.chunks()) > 0

    sections = []

    sentences = []

    # Iterate through doc.chunks()
    for chunk in doc.chunks():
        if isinstance(chunk, Block):
            sentences.extend(chunk.sentences)
        elif isinstance(chunk, Section):
            sentences.extend(chunk.sentences)
        elif isinstance(chunk, ListItem):
            sentences.extend(chunk.sentences)
        elif isinstance(chunk, Paragraph):
            sentences.extend(chunk.sentences)

    print(f"Number of sections: {len(doc.sections())}")
    print(f"Number of chunks: {len(doc.chunks())}")
    print(f"Number of tables: {len(doc.tables())}")
    print(f"Number of sentences: {len(sentences)}")

    # print(f"Table example:\n\n{doc.tables()[4].to_html(True, True)}\n")

if __name__ == "__main__":
    test_layout_pdf_reader()
