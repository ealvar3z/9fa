package main

import (
	"bufio"
	"context"
	"flag"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/openai/openai-go/v3"
)

func runAcmeMode() {
	data, err := io.ReadAll(os.Stdin)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	selection := string(data)
	if strings.TrimSpace(selection) == "" {
		fmt.Println("Error: No selection provided")
		return
	}

	agent := NewAgent(true)
	fmt.Print(agent.Run(acmePrompt, selection))
}

func runDirectMode(prompt string) {
	agent := NewAgent(false)
	fmt.Println(agent.Run(plan9Prompt, prompt))
}

func runREPLLoop(agent Agent, messages []openai.ChatCompletionMessageParamUnion) {
	scanner := bufio.NewScanner(os.Stdin)

	fmt.Println("Plan 9 AI Agent REPL")
	fmt.Println("Type 'quit' or 'exit' to exit, 'clear' to reset conversation")
	fmt.Println()

	for {
		fmt.Print("> ")

		if !scanner.Scan() {
			break
		}

		input := strings.TrimSpace(scanner.Text())
		if input == "" {
			continue
		}

		switch strings.ToLower(input) {
		case "quit", "exit":
			fmt.Println("Goodbye!")
			return
		case "clear":
			messages = []openai.ChatCompletionMessageParamUnion{
				openai.SystemMessage(plan9Prompt),
			}
			fmt.Println("Conversation cleared.")
			continue
		}

		messages = append(messages, openai.UserMessage(input))

		for range 10 {
			resp, err := agent.client.Chat.Completions.New(context.Background(), openai.ChatCompletionNewParams{
				Model:    modelName(),
				Messages: messages,
				Tools:    tools,
			})
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error: %v\n", err)
				break
			}
			if len(resp.Choices) == 0 {
				fmt.Println("Error: No response from LLM")
				break
			}

			msg := resp.Choices[0].Message
			messages = append(messages, msg.ToParam())

			if len(msg.ToolCalls) > 0 {
				for _, tc := range msg.ToolCalls {
					fmt.Fprintf(os.Stderr, "[Tool: %s] %s\n", tc.Function.Name, tc.Function.Arguments)

					result := executeTool(tc.Function.Name, tc.Function.Arguments)
					if len(result) > 2000 {
						result = result[:2000] + "\n... (truncated)"
					}

					messages = append(messages, openai.ToolMessage(result, tc.ID))
				}
				continue
			}

			if msg.Content != "" {
				fmt.Println()
				fmt.Println(msg.Content)
				fmt.Println()
			}
			break
		}
	}
}

func main() {
	acmeMode := flag.Bool("acme", false, "Acme editor integration mode")
	replMode := flag.Bool("repl", false, "Interactive REPL mode")
	flag.Parse()

	switch {
	case *acmeMode:
		runAcmeMode()

	case *replMode:
		NewAgent(false).RunREPL()

	default:
		args := flag.Args()
		if len(args) == 0 {
			fmt.Println("Usage: agent [-acme] [-repl] <prompt>")
			fmt.Println("  -acme    Acme editor mode, reads selection from stdin")
			fmt.Println("  -repl    Interactive REPL mode with conversation history")
			os.Exit(1)
		}

		runDirectMode(strings.Join(args, " "))
	}
}
