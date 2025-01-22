package main

import (
	"github.com/jung-kurt/gofpdf/v2"
)

func generatePDFWithCharacters(filename string) {
	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.AddPage()

	// Set font
	pdf.SetFont("Arial", "", 12)

	// Define the data
	data := []struct {
		English string
		Pinyin  []string
		Chinese []string
	}{
		{"hello", []string{"nǐ", "hǎo"}, []string{"你", "好"}},
		{"goodbye", []string{"zài", "jiàn"}, []string{"再", "见"}},
		{"Chinese, Chinese written language", []string{"zhōng", "wén"}, []string{"中", "文"}},
		{"to welcome", []string{"huān", "yíng"}, []string{"欢", "迎"}},
	}

	// Set initial position
	x, y := 10.0, 20.0
	boxSize := 20.0

	for _, item := range data {
		// Print English word
		pdf.SetXY(x, y)
		pdf.CellFormat(0, 10, item.English, "", 1, "L", false, 0, "")

		// Print Pinyin and Chinese characters
		for i, pinyin := range item.Pinyin {
			// Pinyin
			pdf.SetXY(x+float64(i)*boxSize, y+10)
			pdf.CellFormat(boxSize, 10, pinyin, "1", 0, "C", false, 0, "")

			// Chinese character
			pdf.SetXY(x+float64(i)*boxSize, y+20)
			pdf.SetFont("Arial", "", 24)
			pdf.CellFormat(boxSize, boxSize, item.Chinese[i], "1", 0, "C", false, 0, "")

			// Draw grid
			drawGrid(pdf, x+float64(i)*boxSize, y+20, boxSize)
		}

		// Move to next line
		y += 2*boxSize + 10
	}

	// Save the output file
	err := pdf.OutputFileAndClose(filename)
	if err != nil {
		panic(err)
	}
}

func drawGrid(pdf *gofpdf.Fpdf, x, y, size float64) {
	pdf.SetLineWidth(0.1)
	pdf.SetDrawColor(0, 255, 0) // Green color for grid

	// Draw diagonal lines
	pdf.Line(x, y, x+size, y+size)
	pdf.Line(x+size, y, x, y+size)

	// Draw horizontal and vertical lines
	pdf.Line(x, y+size/2, x+size, y+size/2)
	pdf.Line(x+size/2, y, x+size/2, y+size)
}

func main() {
	generatePDFWithCharacters("chinese_writing_grid.pdf")
}
