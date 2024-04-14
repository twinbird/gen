package main

import (
	"bufio"
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"github.com/adrg/xdg"
	"github.com/google/generative-ai-go/genai"
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

	flag.Parse()

	if *version {
		fmt.Println("ged version 0.0.1")
		return
	}

	args := flag.Args()
	if len(args) == 0 {
		flag.Usage()
		os.Exit(1)
	}

	config, err := getConfig()
	if err != nil {
		log.Fatal("config file load error", err)
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
	askGemini(config, args[0], string(text))
}

func askGemini(config *Config, script string, text string) {
	s := script + "\n" + text

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

type GeminiConfig struct {
	ApiKey string
}

type ChatGptConfig struct {
	ApiKey string
}

type Config struct {
	Gemini  GeminiConfig
	ChatGpt ChatGptConfig
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

func getConfigPath() string {
	return filepath.Join(xdg.ConfigHome, CONFIG_DIR_NAME, CONFIG_FILE_NAME)
}
