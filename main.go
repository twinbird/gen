package main

import (
	"bufio"
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"github.com/adrg/xdg"
	"github.com/google/generative-ai-go/genai"
	openai "github.com/sashabaranov/go-openai"
	"google.golang.org/api/option"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
)

const (
	CONFIG_DIR_NAME  = "gen"
	CONFIG_FILE_NAME = "setting.json"
)

func usage() {
	fmt.Fprintf(os.Stderr, "Usage: %s [options...] script [file]\n", os.Args[0])
	flag.PrintDefaults()
}

func main() {
	flag.Usage = usage
	version := flag.Bool("version", false, "display version information")
	doConfigure := flag.Bool("configure", false, "output configuration file")

	flag.Parse()

	if *version {
		fmt.Println("gen version 0.0.1")
		return
	}

	if *doConfigure {
		if fileExists(getConfigPath()) {
			return
		}
		outputConfigurationFile()
		return
	}

	args := flag.Args()
	if len(args) == 0 {
		flag.Usage()
		os.Exit(1)
	}

	config, err := getConfig()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Configuration file load error.\nSee %s\nOr do `gen -configure`\n", getConfigPath())
		return
	}

	text := ""
	if len(args) > 1 {
		b, err := ioutil.ReadFile(args[1])
		if err != nil {
			log.Fatal("file read error", err)
		}
		text = string(b)
	} else {
		scanner := bufio.NewScanner(os.Stdin)
		scanner.Scan()
		text = scanner.Text()
	}

	if config.DefaultUseService == "chatgpt" {
		askChatGpt(config, args[0], string(text))
	} else {
		askGemini(config, args[0], string(text))
	}
}

func askChatGpt(config *Config, script string, text string) {
	s := createScript(config, script, text)
	client := openai.NewClient(config.ChatGpt.ApiKey)
	resp, err := client.CreateChatCompletion(
		context.Background(),
		openai.ChatCompletionRequest{
			Model: openai.GPT3Dot5Turbo,
			Messages: []openai.ChatCompletionMessage{
				{
					Role:    openai.ChatMessageRoleUser,
					Content: s,
				},
			},
		},
	)

	if err != nil {
		fmt.Printf("ChatCompletion error: %v\n", err)
		return
	}

	fmt.Println(resp.Choices[0].Message.Content)
}

func askGemini(config *Config, script string, text string) {
	s := createScript(config, script, text)

	ctx := context.Background()
	client, err := genai.NewClient(ctx, option.WithAPIKey(config.Gemini.ApiKey))
	if err != nil {
		log.Fatal("gemini client error", err)
	}
	defer client.Close()

	model := client.GenerativeModel("gemini-pro")
	resp, err := model.GenerateContent(ctx, genai.Text(s))
	if err != nil {
		log.Fatal("gemini generateContent error", err)
	}
	printGeminiResponse(resp)
}

func printGeminiResponse(resp *genai.GenerateContentResponse) {
	for _, cand := range resp.Candidates {
		if cand.Content != nil {
			for _, part := range cand.Content.Parts {
				fmt.Println(part)
			}
		}
	}
}

func createScript(config *Config, script string, fileText string) string {
	return script + "\n" + fileText
}

type GeminiConfig struct {
	ApiKey string
}

type ChatGptConfig struct {
	ApiKey string
}

type Config struct {
	DefaultUseService string
	Gemini            GeminiConfig
	ChatGpt           ChatGptConfig
}

func getConfig() (*Config, error) {
	path := getConfigPath()
	b, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	config := &Config{}
	if err := json.Unmarshal(b, config); err != nil {
		return nil, err
	}
	return config, nil
}

func getConfigDir() string {
	return filepath.Join(xdg.ConfigHome, CONFIG_DIR_NAME)
}

func getConfigPath() string {
	return filepath.Join(getConfigDir(), CONFIG_FILE_NAME)
}

func outputConfigurationFile() {
	dir := getConfigDir()

	if err := os.MkdirAll(dir, 0755); err != nil {
		log.Fatal(err)
	}

	config := &Config{}
	b, err := json.MarshalIndent(config, "", "\t")
	if err != nil {
		log.Fatal("config file marshal error")
	}
	err = ioutil.WriteFile(getConfigPath(), b, os.ModePerm)
	if err != nil {
		log.Fatal("config file write error")
	}
}

func fileExists(filename string) bool {
	_, err := os.Stat(filename)
	return err == nil
}
