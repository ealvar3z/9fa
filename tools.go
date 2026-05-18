package main

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/openai/openai-go/v3"
)

var tools = []openai.ChatCompletionToolUnionParam{
	tool("run_command", "Execute an rc shell command on Plan 9 and return its output", map[string]any{
		"command": map[string]any{
			"type":        "string",
			"description": "The rc shell command to execute",
		},
	}, []string{"command"}),

	tool("read_file", "Read the contents of a file", map[string]any{
		"path": map[string]any{
			"type":        "string",
			"description": "The file path to read",
		},
	}, []string{"path"}),

	tool("write_file", "Write content to a file", map[string]any{
		"path": map[string]any{
			"type":        "string",
			"description": "The file path to write to",
		},
		"content": map[string]any{
			"type":        "string",
			"description": "The content to write",
		},
	}, []string{"path", "content"}),

	tool("list_directory", "List files in a directory", map[string]any{
		"path": map[string]any{
			"type":        "string",
			"description": "The directory path to list",
		},
	}, []string{"path"}),
}

func tool(name, description string, properties map[string]any, required []string) openai.ChatCompletionToolUnionParam {
	return openai.ChatCompletionFunctionTool(openai.FunctionDefinitionParam{
		Name:        name,
		Description: openai.String(description),
		Parameters: openai.FunctionParameters{
			"type":       "object",
			"properties": properties,
			"required":   required,
		},
	})
}

func executeTool(name, args string) string {
	var params map[string]string
	if err := json.Unmarshal([]byte(args), &params); err != nil {
		return fmt.Sprintf("Error: invalid tool arguments: %v", err)
	}

	switch name {
	case "run_command":
		cmd := exec.Command("rc", "-c", params["command"])
		out, err := cmd.CombinedOutput()
		if err != nil {
			return fmt.Sprintf("Error: %v\nOutput: %s", err, out)
		}
		return string(out)

	case "read_file":
		data, err := os.ReadFile(params["path"])
		if err != nil {
			return fmt.Sprintf("Error: %v", err)
		}
		return string(data)

	case "write_file":
		if err := os.WriteFile(params["path"], []byte(params["content"]), 0644); err != nil {
			return fmt.Sprintf("Error: %v", err)
		}
		return "File written successfully"

	case "list_directory":
		entries, err := os.ReadDir(params["path"])
		if err != nil {
			return fmt.Sprintf("Error: %v", err)
		}

		names := make([]string, 0, len(entries))
		for _, e := range entries {
			name := e.Name()
			if e.IsDir() {
				name += "/"
			}
			names = append(names, name)
		}
		return strings.Join(names, "\n")

	default:
		return fmt.Sprintf("Unknown tool: %s", name)
	}
}
