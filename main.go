package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/jung-kurt/gofpdf"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

type LanguageData struct {
	English string
	Pinyin  []string
	Chinese []string
}

type OpenRouterRequest struct {
	Model    string    `json:"model"`
	Messages []Message `json:"messages"`
}

type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type OpenRouterResponse struct {
	Choices []struct {
		Message struct {
			Content string `json:"content"`
		} `json:"message"`
	} `json:"choices"`
}

func main() {
	useCalligraphyFont := true
	pdf := setupPDF(useCalligraphyFont)

	// Get predefined data
	data := getLanguageData()

	// Get user input
	fmt.Print("Enter an English phrase to translate (or press Enter to skip): ")
	var userInput string
	scanner := bufio.NewScanner(os.Stdin)
	if scanner.Scan() {
		userInput = scanner.Text()
	}

	// If user provided input, get AI translation
	if userInput != "" {
		fmt.Println("Getting translation...")
		aiTranslation, err := getAITranslation(userInput)
		if err != nil {
			fmt.Printf("Error getting translation: %v\n", err)
		} else {
			data = append(data, aiTranslation)
			fmt.Println("Translation successful!")
		}
	}

	x, y := 20.0, 20.0
	boxSize := 20.0

	for _, item := range data {
		addEnglishText(pdf, item.English, x, y)
		addPinyinAndChinese(pdf, item.Pinyin, item.Chinese, x, y, boxSize, useCalligraphyFont)
		y += boxSize + 35
	}

	if err := pdf.OutputFileAndClose("output.pdf"); err != nil {
		fmt.Printf("Error creating PDF: %v\n", err)
		return
	}

	fmt.Println("PDF created successfully!")
}

func setupPDF(useCalligraphyFont bool) *gofpdf.Fpdf {
	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.AddPage()

	fontDir := "/fonts"
	defaultFontFile := "chinese.msyh.ttf"
	calligraphyFontFile := "simsun.ttf"

	// Load the default font
	pdf.AddUTF8Font("YaHei", "", filepath.Join(fontDir, defaultFontFile))

	// Load the calligraphy font if needed
	if useCalligraphyFont {
		pdf.AddUTF8Font("Calligraphy", "", filepath.Join(fontDir, calligraphyFontFile))
	}

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

func addPinyinAndChinese(pdf *gofpdf.Fpdf, pinyin, chinese []string, x, y, boxSize float64, useCalligraphyFont bool) {
	pageWidth, _ := pdf.GetPageSize()
	marginLeft, _, marginRight, _ := pdf.GetMargins()
	usableWidth := pageWidth - marginLeft - marginRight
	numBoxes := int(usableWidth / (boxSize + 2))

	for i := 0; i < numBoxes; i++ {
		xPos := x + float64(i)*(boxSize+2)

		// Add Pinyin text if available
		if i < len(pinyin) {
			pdf.SetTextColor(0, 0, 0)
			pdf.SetFont("YaHei", "", 12)
			pdf.SetXY(xPos, y+12)
			pdf.CellFormat(boxSize, 10, pinyin[i], "1", 0, "C", false, 0, "")
		}

		// Add Chinese characters if available
		if i < len(chinese) {
			pdf.SetTextColor(200, 200, 200)
			if useCalligraphyFont {
				pdf.SetFont("Calligraphy", "", 40)
			} else {
				pdf.SetFont("YaHei", "", 40)
			}
			pdf.SetXY(xPos, y+22)
			pdf.CellFormat(boxSize, boxSize, chinese[i], "1", 0, "C", false, 0, "")
		} else {
			// Draw empty grid
			pdf.SetTextColor(200, 200, 200)
			pdf.SetFont("YaHei", "", 40)
			pdf.SetXY(xPos, y+22)
			pdf.CellFormat(boxSize, boxSize, "", "1", 0, "C", false, 0, "")
		}
	}
}

func getAITranslation(englishText string) (LanguageData, error) {
	apiKey := os.Getenv("OPENROUTER_API_KEY") // Get API key from environment variable
	if apiKey == "" {
		return LanguageData{}, fmt.Errorf("OPENROUTER_API_KEY environment variable not set")
	}

	prompt := `Translate the following English text to Chinese (Simplified) and Pinyin.
    Format the response exactly like this example:
    English: hello
    Pinyin: nǐ hǎo
    Chinese: 你好

    English: ` + englishText

	// Prepare the request
	reqBody := OpenRouterRequest{
		Model: "allenai/molmo-7b-d:free",
		Messages: []Message{
			{
				Role:    "user",
				Content: prompt,
			},
		},
	}

	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return LanguageData{}, fmt.Errorf("error marshaling request: %v", err)
	}

	fmt.Println("Request JSON:", string(jsonData))

	// Create the HTTP request
	req, err := http.NewRequest("POST", "https://openrouter.ai/api/v1/chat/completions",
		bytes.NewBuffer(jsonData))
	if err != nil {
		return LanguageData{}, fmt.Errorf("error creating request: %v", err)
	}

	// Set headers
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+apiKey)

	// Make the request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return LanguageData{}, fmt.Errorf("error making request: %v", err)
	}
	defer resp.Body.Close()

	// Read the response
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return LanguageData{}, fmt.Errorf("error reading response: %v", err)
	}

	fmt.Println("Response Status Code:", resp.StatusCode)
	fmt.Println("Raw Response Body:", string(body))

	// Parse the response
	var openRouterResp OpenRouterResponse
	if err := json.Unmarshal(body, &openRouterResp); err != nil {
		return LanguageData{}, fmt.Errorf("error parsing response: %v", err)
	}

	fmt.Println("Parsed OpenRouter Response:")
	fmt.Printf("%+v\n", openRouterResp)

	if len(openRouterResp.Choices) == 0 {
		return LanguageData{}, fmt.Errorf("no response from API")
	}

	// Parse the content
	response := openRouterResp.Choices[0].Message.Content
	fmt.Println("API Response Content:", response)

	lines := strings.Split(response, "\n")

	var result LanguageData
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, "English:") {
			result.English = strings.TrimSpace(strings.TrimPrefix(line, "English:"))
		} else if strings.HasPrefix(line, "Pinyin:") {
			pinyinStr := strings.TrimSpace(strings.TrimPrefix(line, "Pinyin:"))
			result.Pinyin = strings.Split(pinyinStr, " ")
		} else if strings.HasPrefix(line, "Chinese:") {
			chineseStr := strings.TrimSpace(strings.TrimPrefix(line, "Chinese:"))
			result.Chinese = strings.Split(chineseStr, "")
		}
	}

	fmt.Printf("Final Result: %+v\n", result)

	return result, nil
}
