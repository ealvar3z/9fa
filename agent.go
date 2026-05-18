package main

import (
	"context"
	"fmt"
	"os"

	"github.com/openai/openai-go/v3"
)

type Agent struct {
	client openai.Client
	quiet  bool
}

func NewAgent(quiet bool) Agent {
	return Agent{
		client: newOpenAIClient(),
		quiet:  quiet,
	}
}

func (a Agent) Run(systemPrompt, userPrompt string) string {
	messages := []openai.ChatCompletionMessageParamUnion{
		openai.SystemMessage(systemPrompt),
		openai.UserMessage(userPrompt),
	}

	return a.runMessages(context.Background(), messages)
}

func (a Agent) runMessages(ctx context.Context, messages []openai.ChatCompletionMessageParamUnion) string {
	for range 10 {
		resp, err := a.client.Chat.Completions.New(ctx, openai.ChatCompletionNewParams{
			Model:    modelName(),
			Messages: messages,
			Tools:    tools,
		})
		if err != nil {
			return fmt.Sprintf("Error: %v", err)
		}
		if len(resp.Choices) == 0 {
			return "Error: No response from LLM"
		}

		msg := resp.Choices[0].Message
		messages = append(messages, msg.ToParam())

		if len(msg.ToolCalls) == 0 {
			return msg.Content
		}

		for _, tc := range msg.ToolCalls {
			if !a.quiet {
				fmt.Fprintf(os.Stderr, "[Tool: %s] %s\n", tc.Function.Name, tc.Function.Arguments)
			}

			result := executeTool(tc.Function.Name, tc.Function.Arguments)
			if len(result) > 2000 {
				result = result[:2000] + "\n... (truncated)"
			}

			messages = append(messages, openai.ToolMessage(result, tc.ID))
		}
	}

	return "Error: Max iterations reached"
}

func (a Agent) RunREPL() {
	messages := []openai.ChatCompletionMessageParamUnion{
		openai.SystemMessage(plan9Prompt),
	}

	runREPLLoop(a, messages)
}
