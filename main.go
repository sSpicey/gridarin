package main

import (
	"path/filepath"

	"github.com/jung-kurt/gofpdf"
)

type LanguageData struct {
	English string
	Pinyin  []string
	Chinese []string
}

func main() {
	pdf := setupPDF()
	data := getLanguageData()

	x, y := 20.0, 20.0
	boxSize := 30.0

	for _, item := range data {
		addEnglishText(pdf, item.English, x, y)
		addPinyinAndChinese(pdf, item.Pinyin, item.Chinese, x, y, boxSize)
		y += boxSize + 35
	}

	if err := pdf.OutputFileAndClose("output.pdf"); err != nil {
		panic(err)
	}
}

func setupPDF() *gofpdf.Fpdf {
	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.AddPage()

	fontDir := "/fonts"
	fontFile := "chinese.msyh.ttf"
	pdf.AddUTF8Font("YaHei", "", filepath.Join(fontDir, fontFile))

	return pdf
}

func getLanguageData() []LanguageData {
	return []LanguageData{
		{"hello", []string{"nǐ", "hǎo"}, []string{"你", "好"}},
		{"goodbye", []string{"zài", "jiàn"}, []string{"再", "见"}},
		{"Chinese, Chinese written language", []string{"zhōng", "wén"}, []string{"中", "文"}},
		{"to welcome", []string{"huān", "yíng"}, []string{"欢", "迎"}},
	}
}

func addEnglishText(pdf *gofpdf.Fpdf, text string, x, y float64) {
	pdf.SetTextColor(0, 0, 0)
	pdf.SetFont("Arial", "B", 12)
	pdf.SetXY(x, y)
	pdf.CellFormat(0, 10, text, "", 1, "L", false, 0, "")
}

func addPinyinAndChinese(pdf *gofpdf.Fpdf, pinyin, chinese []string, x, y, boxSize float64) {
	for i, p := range pinyin {
		xPos := x + float64(i)*(boxSize+2)

		// Add Pinyin text
		pdf.SetTextColor(0, 0, 0)
		pdf.SetFont("YaHei", "", 12)
		pdf.SetXY(xPos, y+12)
		pdf.CellFormat(boxSize, 10, p, "1", 0, "C", false, 0, "")

		// Add Chinese characters
		pdf.SetTextColor(200, 200, 200)
		pdf.SetFont("YaHei", "", 60)
		pdf.SetXY(xPos, y+22)
		pdf.CellFormat(boxSize, boxSize, chinese[i], "1", 0, "C", false, 0, "")
	}
}
