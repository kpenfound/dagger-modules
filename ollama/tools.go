package main

import (
	"encoding/json"

	"github.com/ollama/ollama/api"
)

func System() []api.Message {
	system := []api.Message{
		{
			Role:    "system",
			Content: "You are an assistant with access to the tool Dagger. You can use the Dagger tool to interact with the Dagger modules. Dagger modules can be used to lint and test code repositories on GitHub.",
		},
	}

	return system
}

const toolsJson = `
[
	{
		"type": "function",
		"function": {
	        "name": "getDaggerFunctions",
	        "description": "Gets usage information for a dagger module by listing its available functions",
	        "parameters": {
	            "type": "object",
	            "properties": {
		            "module": {
			            "type": "string",
			            "description": "The dagger module to list functions for"
		            }
		        },
		        "required": ["module"]
	        }
        }
	},
	{
		"type": "function",
		"function": {
	        "name": "getDaggerHelp",
	        "description": "Gets usage information for a dagger function finding information about the arguments and returned information",
	        "parameters": {
	            "type": "object",
	            "properties": {
		            "module": {
			            "type": "string",
			            "description": "The dagger module to get help for"
		            },
					"subcommand": {
			            "type": "string",
			            "description": "The subcommand to get help with"
		            }
		        },
		        "required": ["subcommand"]
	        }
        }
	},
	{
		"type": "function",
		"function": {
	        "name": "callDaggerFunction",
	        "description": "Call a dagger function in a dagger module with the dagger CLI",
	        "parameters": {
	            "type": "object",
	            "properties": {
		            "module": {
			            "type": "string",
			            "description": "The dagger module to call a function from"
		            },
					"function": {
			            "type": "string",
			            "description": "The dagger function to call"
		            },
					"arguments": {
						"type": "string",
				     	"description": "The arguments to pass to the dagger function"
				    }
		        },
		        "required": ["module", "function"]
	        }
        }
	}
]
`

func callDagger(module, function, arguments string) string {
	command := []string{"dagger", "-m", module}
	if function != "" {
		command = append(command, "call", function)
	}
	command = append(command, arguments)

	return ""
}

func CallTool(call api.ToolCall) api.Message {
	message := api.Message{
		Role: "tool",
	}
	switch call.Function.Name {
	case "getDaggerFunctions":
		message.Content = callDagger(
			call.Function.Arguments["module"].(string),
			"",
			"functions",
		)
	case "getDaggerHelp":
		message.Content = callDagger(
			call.Function.Arguments["module"].(string),
			call.Function.Arguments["subcommand"].(string),
			"--help",
		)
	case "callDaggerFunction":
		args := ""
		if val, ok := call.Function.Arguments["arguments"]; ok {
			args = val.(string)
		}
		message.Content = callDagger(
			call.Function.Arguments["module"].(string),
			call.Function.Arguments["subcommand"].(string),
			args,
		)
	}

	return message
}
func Tools() api.Tools {
	tools := []api.Tool{}
	err := json.Unmarshal([]byte(toolsJson), &tools)
	if err != nil {
		panic(err)
	}

	return tools
}
