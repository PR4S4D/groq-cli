# GROQ-CLI

A simple GO CLI tool using GROQ API and Charm.sh libraries

## Installation

```sh
  go mod tidy
  go run .

```

OR use the binary in the repo

OR `go build .`

## Env variables

- Rename .env.example to .env and configure the values
- OR configure GROQ_API_KEY, GROQ_MODEL and GROQ_MODELS env variables

- For windows:
  
![Window Env](./media/grog-cli-env-windows.png)

## Usage

```sh
  grog-cli.exe
  groq-cli.exe -f File.md # reads the prompt from the file
  groq-cli.exe -A -f File.md # uses all the available models (GROQ_MODELS env var) and saves the responses to a file
```

https://github.com/PR4S4D/groq-cli/assets/20255076/701fb865-79d2-495a-b14b-832300446942



