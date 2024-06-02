package main

import (
	"encoding/json"
	"fmt"
	"groq-cli/types"
	"groq-cli/utils"
	"log"
	"os"
	"strings"

	"github.com/charmbracelet/huh"
	"github.com/charmbracelet/huh/spinner"
	"github.com/go-resty/resty/v2"
	"github.com/joho/godotenv"
)

var client *resty.Client

func main() {

	if err := godotenv.Load(); err != nil && !os.IsNotExist(err) {
		log.Fatalln("Error loading .env")
	}

	initRestClient()
	var output string
	var input string
	var prevInput string
	var response string

	for {
		input = ""
		fmt.Println()
		form := huh.NewForm(
			huh.NewGroup(
				huh.NewText().
					Title("How can I help you...?").
					Placeholder("Message...").
					Lines(3).
					Value(&input),
			),
		).WithTheme(huh.ThemeDracula()).WithShowHelp(false)

		err := form.Run()
		if err != nil {
			log.Fatal(err)
		}

		if prevInput != "" && (input == "s" || input == "S" || strings.ToLower(input) == "save") {
			utils.SaveToFile(prevInput, response)
			continue
		}

		if len([]byte(input)) > 1 {
			title := utils.RenderMarkDown("# "+utils.ToTitleCase(input), "pink")

			action := func() {
				response = getGrogResponse(input)
				output = utils.RenderMarkDown(response, "dracula")
			}
			prevInput = input
			_ = spinner.New().
				Title(" Processing...").
				Action(action).
				Type(utils.GetSpinnerStyle()).Run()

			fmt.Println(title)
			fmt.Println(output)
			fmt.Print(utils.RenderMarkDown("---", "pink"))
		}
	}
}

func initRestClient() {
	groqAPIKey, foundKey := os.LookupEnv("GROQ_API_KEY")

	if foundKey == false {
		panic("GROQ_API_KEY not found")
	}

	client = resty.New()
	client.SetHeader("Authorization", "Bearer "+groqAPIKey)
	client.SetHeader("Content-Type", "application/json")

}

func getGrogResponse(input string) string {
	resp, err := client.R().SetBody(getRequestBody(input)).Post("https://api.groq.com/openai/v1/chat/completions")

	if err != nil {
		log.Fatal(err)
	}

	response := types.MessageResponse{}
	json.Unmarshal(resp.Body(), &response)

	if resp.StatusCode() != 200 {
		fmt.Println(response, resp.StatusCode())
		os.Exit(1)
	}

	return response.Choices[0].Message.Content

}

func getRequestBody(input string) string {
	var grogModel = os.Getenv("GROQ_MODEL")
	if len(grogModel) == 0 {
		grogModel = "llama3-8b-8192"
	}
	return fmt.Sprintf(`{
		"model": "%s",
		"stream": false,
		"max_tokens": 8192,
		"messages": [
			{"role": "user", "content": "%s"}
		]
	}`, grogModel, input)
}
