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

    # Iterate through doc.sections()
    for section in doc.sections():
        section_text = section.to_text(True, False)
        log(debug, f"Section: {section_text}")

        def iter_children(node, level):
            for child in node.children:
                if isinstance(child, Paragraph):
                    log(debug, f"{'  ' * level}-Paragraph")
                    log(debug, f"# Sentences: {len(child.sentences)}")
                    log(debug, f"PARAGRAPH SPOTTED: {' '.join(child.sentences)}")
                elif isinstance(child, Block):
                    for grandchild in child.children:
                        if isinstance(grandchild, Section):
                            log(debug, f"{'  ' * level}-Section")
                            log(debug, f"# Sentences: {len(grandchild.sentences)}")
                        iter_children(grandchild, level + 1)
                elif isinstance(child, Section):
                    for grandchild in child.children:
                        if isinstance(grandchild, Section):
                            log(debug, f"{'  ' * level}-Section")
                            log(debug, f"# Sentences: {len(grandchild.sentences)}")
                        iter_children(grandchild, level + 1)
                    sections.append(child.to_context_text(True))
                elif isinstance(child, ListItem):
                    log(debug, f"{'  ' * level}-ListItem")
                    log(debug, f"# Sentences: {len(child.sentences)}")
                    log(debug, f"LIST ITEM SPOTTED: {' '.join(child.sentences)}")
                    for grandchild in child.children:
                        if isinstance(grandchild, Section):
                            log(debug, f"{'  ' * level}-Section")
                            log(debug, f"# Sentences: {len(grandchild.sentences)}")
                        iter_children(grandchild, level + 1)
                else:
                    log(debug, f"Unsupported block type: {type(child)}")

        iter_children(section, 1)

    sentences = []

    # Iterate through doc.chunks()
    for chunk in doc.chunks():
        if isinstance(chunk, Block):
            log(debug, f"Number of sentences in block: {len(chunk.sentences)}")
            sentences.extend(chunk.sentences)
        elif isinstance(chunk, Section):
            log(debug, f"Number of sentences in section: {len(chunk.sentences)}")
            sentences.extend(chunk.sentences)
        elif isinstance(chunk, ListItem):
            log(debug, f"Number of sentences in list item: {len(chunk.sentences)}")
            sentences.extend(chunk.sentences)
        elif isinstance(chunk, Paragraph):
            log(debug, f"Number of sentences in paragraph: {len(chunk.sentences)}")
            sentences.extend(chunk.sentences)

    # Iterate and log each sentence with a newline between
    for sentence in sentences:
        log(debug, f"\n\n{sentence}\n\n")

    print(f"Number of sections: {len(doc.sections())}")
    print(f"Number of chunks: {len(doc.chunks())}")
    print(f"Number of tables: {len(doc.tables())}")
    print(f"Number of sentences: {len(sentences)}")

    # print(f"Table example:\n\n{doc.tables()[4].to_html(True, True)}\n")

if __name__ == "__main__":
    test_layout_pdf_reader()
