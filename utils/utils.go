package utils

import (
	"fmt"
	"math/rand"
	"os"
	"regexp"
	"strings"

	"github.com/charmbracelet/glamour"
	"github.com/charmbracelet/huh/spinner"
)

func SaveToFile(input string, output string) {
	dir := "notes"
	err := os.MkdirAll(dir, 0644)
	if err != nil {
		fmt.Println(err)
		return
	}
	file := getFileName(input)
	fmt.Println(RenderMarkDown(fmt.Sprintf("Saving to `%s`", file), "dracula"))
	os.WriteFile(file, []byte(output), 0644)
}

func getFileName(s string) string {
	file := toSnakeCase(s)
	return "notes/" + cleanFilename(strings.TrimSpace(fmt.Sprintf("%.*s", 60, file))) + ".md"
}

func cleanFilename(filename string) string {
	var re = regexp.MustCompile(`[^\w\d_.-]+`)
	cleanedFilename := re.ReplaceAllString(filename, "")
	return cleanedFilename
}

func toSnakeCase(s string) string {
	return strings.ReplaceAll(strings.ToLower(s), " ", "_")
}

func RenderMarkDown(str string, theme string) string {
	output, _ := glamour.Render(str, theme)
	return output
}

var spinnerStyles = []spinner.Type{
	spinner.Line,
	spinner.Dots,
	spinner.MiniDot,
	spinner.Jump,
	spinner.Pulse,
	spinner.Points,
	spinner.Globe,
	spinner.Moon,
	spinner.Monkey,
	spinner.Meter,
	spinner.Hamburger,
	spinner.Ellipsis,
}

func GetSpinnerStyle() spinner.Type {
	return spinnerStyles[rand.Intn(12)]
}

func ToTitleCase(s string) string {
	s = strings.Trim(s, " ")
	if len([]byte(s)) < 1 {
		return s
	}
	words := strings.Split(s, " ")
	for i, word := range words {
		words[i] = strings.ToUpper(word[:1]) + strings.ToLower(word[1:])
	}
	return strings.Join(words, " ")
}
