package main

import (
	"crypto/tls"
	"net/http"
	"os"

	"github.com/openai/openai-go/v3"
	"github.com/openai/openai-go/v3/option"
)

const defaultModel = "gpt-5.2"
const apiKey = "..."

func modelName() openai.ChatModel {
	model := os.Getenv("OPENAI_MODEL")
	if model == "" {
		model = defaultModel
	}
	return openai.ChatModel(model)
}

func newOpenAIClient() openai.Client {
	httpClient := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
	}

	return openai.NewClient(
		option.WithAPIKey(apiKey),
		option.WithHTTPClient(httpClient),
	)
}
