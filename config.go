package main

import (
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/openai/openai-go/v3"
	"github.com/openai/openai-go/v3/option"
)

const defaultModel = "gpt-5.2"
const apiKeyEnv = "OPENAI_API_KEY"

func modelName() openai.ChatModel {
	model := os.Getenv("OPENAI_MODEL")
	if model == "" {
		model = defaultModel
	}
	return openai.ChatModel(model)
}

func openAIAPIKey() string {
	apiKey := strings.TrimSpace(os.Getenv(apiKeyEnv))
	if apiKey == "" {
		fmt.Fprintf(os.Stderr, "Error: %s is not set\n", apiKeyEnv)
		os.Exit(1)
	}

	return apiKey
}

func newOpenAIClient() openai.Client {
	httpClient := &http.Client{}

	return openai.NewClient(
		option.WithAPIKey(openAIAPIKey()),
		option.WithHTTPClient(httpClient),
	)
}
