// A module for Ollama and Dagger

package main

import (
	"context"
	"dagger/ollama/internal/dagger"
	"fmt"
)

type Ollama struct {
	Service *dagger.Service
}

func New(
	ctx context.Context,
	endpoint *dagger.Service,
) *Ollama {
	return &Ollama{
		Service: endpoint,
	}
}

// Returns a container that echoes whatever string argument is provided
func (m *Ollama) Ask(ctx context.Context, prompt string) (string, error) {
	endpoint, err := m.Service.Endpoint(ctx)
	if err != nil {
		return "", err
	}
	fmt.Println("Endpoint: ", endpoint)

	client := NewOllamaClient("http://" + endpoint)

	fmt.Println("Prompt: ", prompt)
	response, err := client.Ask(ctx, prompt)
	if err != nil {
		return "", err
	}

	return response, nil
}
