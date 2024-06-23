package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"groq-cli/types"
	"groq-cli/utils"
	"log"
	"os"
	"strings"
	"sync"

	"github.com/charmbracelet/huh"
	"github.com/charmbracelet/huh/spinner"
	"github.com/go-resty/resty/v2"
	"github.com/joho/godotenv"
)

// Message represents the structure of a chat message.
type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

var (
	client     *resty.Client
	grogModel  string
	grogModels []string
)

func initialize() {
	if err := godotenv.Load(); err != nil && !os.IsNotExist(err) {
		log.Fatalln("Error loading .env")
	}

	initRestClient()

	grogModel = os.Getenv("GROQ_MODEL")
	grogModels = strings.Split(os.Getenv("GROQ_MODELS"), " ")
	if len(grogModel) == 0 {
		grogModel = "llama3-8b-8192"
	}
	if len(grogModels) == 0 {
		grogModels = []string{"llama3-8b-8192"}
	}
}

func main() {
	initialize()

	var output, input, prevInput, response string

	filePtr := flag.String("f", "default", "Prompt file")
	useAllModels := flag.Bool("A", false, "Use all models")
	flag.Parse()

	// handle inputs from files
	if *filePtr != "default" {
		filePrompt, err := os.ReadFile(*filePtr)
		check(err)
		input = strings.TrimSpace(string(filePrompt))
		if *useAllModels {
			handleAllModels(input, *filePtr)
			os.Exit(0)
		} else {
			response = getGrogResponse(input, grogModel)
			output = utils.RenderMarkDown(response, "dracula")
			fmt.Print(output)
			os.Exit(0)
		}
	}

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
		check(err)

		if prevInput != "" && (input == "s" || input == "S" || strings.ToLower(input) == "save") {
			utils.SaveToFile(prevInput, response)
			continue
		}

		if len([]byte(input)) > 1 {
			title := utils.RenderMarkDown("# "+utils.ToTitleCase(input), "pink")

			action := func() {
				response = getGrogResponse(input, grogModel)
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

func getGrogResponse(input string, model string) string {
	resp, err := client.R().SetBody(getRequestBody(input, model)).Post("https://api.groq.com/openai/v1/chat/completions")

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

func handleAllModels(input, filePtr string) {
	outputFileName := strings.Split(filePtr, ".")[0] + "-output.md"

	var wg sync.WaitGroup

	f, err := os.Create(outputFileName)
	check(err)
	defer f.Close()

	for _, model := range grogModels {
		wg.Add(1)
		go func(model string) {
			defer wg.Done()
			response := getGrogResponse(input, model)
			f.WriteString(fmt.Sprintf("## %s \n\n", model))
			f.WriteString(response)
			f.WriteString("--- \n\n\n")
		}(model)
	}
	wg.Wait()
	fmt.Printf("Writing to file: %s\n", outputFileName)
	os.Exit(0)
}

func getRequestBody(input string, model string) string {
	if len(model) == 0 {
		model = "llama3-8b-8192"
	}

	var requestBody struct {
		Model    string    `json:"model"`
		Stream   bool      `json:"stream"`
		Messages []Message `json:"messages"`
	}

	requestBody.Model = model
	requestBody.Stream = false
	requestBody.Messages = []Message{{Role: "user", Content: input}}

	w := bytes.NewBuffer(nil)
	enc := json.NewEncoder(w)
	err := enc.Encode(&requestBody)
	if err != nil {
		panic(err)
	}
	jsonReq := w.Bytes()

	// fmt.Println(string(jsonReq))

	return string(jsonReq)
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}
