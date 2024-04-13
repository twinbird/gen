package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/google/generative-ai-go/genai"
	"google.golang.org/api/option"
	"io/ioutil"
	"log"
	"os"
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
		fmt.Println("ased version 0.0.1")
		return
	}

	args := flag.Args()
	if len(args) == 0 {
		flag.Usage()
		os.Exit(1)
	}

	text, err := ioutil.ReadFile(args[1])
	if err != nil {
		log.Fatal("file read error", err)
	}
	withGemini(args[0], string(text))
}

func withGemini(script string, text string) {
	s := script + "\n" + text

	ctx := context.Background()
	client, err := genai.NewClient(ctx, option.WithAPIKey(os.Getenv("API_KEY")))
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
