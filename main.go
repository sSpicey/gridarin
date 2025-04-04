package main

import (
	"path/filepath"

	"github.com/jung-kurt/gofpdf"
)

func main() {
	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.AddPage()

	fontDir := "/fonts"
	fontFile := "chinese.msyh.ttf"

	pdf.AddUTF8Font("YaHei", "", filepath.Join(fontDir, fontFile))

	x, y := 20.0, 20.0
	boxSize := 30.0

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

	for _, item := range data {
		// Set text color to black for English text
		pdf.SetTextColor(0, 0, 0)
		pdf.SetFont("Arial", "B", 12)
		pdf.SetXY(x, y)
		pdf.CellFormat(0, 10, item.English, "", 1, "L", false, 0, "")

		for i, pinyin := range item.Pinyin {
			xPos := x + float64(i)*(boxSize+2)

			// Set text color to black for Pinyin text
			pdf.SetTextColor(0, 0, 0)
			pdf.SetFont("YaHei", "", 12)
			pdf.SetXY(xPos, y+12)
			pdf.CellFormat(boxSize, 10, pinyin, "1", 0, "C", false, 0, "")

			// Set light grey color for Chinese characters
			pdf.SetTextColor(200, 200, 200)
			pdf.SetFont("YaHei", "", 60)
			pdf.SetXY(xPos, y+22)
			pdf.CellFormat(boxSize, boxSize, item.Chinese[i], "1", 0, "C", false, 0, "")
		}

		y += boxSize + 35
	}

	err := pdf.OutputFileAndClose("output.pdf")
	if err != nil {
		panic(err)
	}
}
