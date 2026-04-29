import sys
from pypdf import PdfReader

def extract_text(pdf_path, output_path):
    try:
        reader = PdfReader(pdf_path)
        text = ""
        for i, page in enumerate(reader.pages):
            text += f"\n--- Page {i+1} ---\n"
            text += page.extract_text() or ""
        
        with open(output_path, "w", encoding="utf-8") as f:
            f.write(text)
        print(f"Extracted text to {output_path}")
    except Exception as e:
        print(f"Error: {e}")

if __name__ == "__main__":
    pdf_path = "D:\\desenvolvimento\\taxas-sejusp\\Docs\\Manual de Cadastro e Autenticação de sistemas com o Portal ICMS Transparente.pdf"
    output_path = "D:\\desenvolvimento\\taxas-sejusp\\Docs\\manual_extratado.txt"
    extract_text(pdf_path, output_path)
