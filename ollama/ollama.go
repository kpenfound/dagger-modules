package main

import (
	"context"
	"fmt"
	"net/http"

	"github.com/ollama/ollama/api"
)

const (
	MODEL = "llama3.1"
)

type OllamaClient struct {
	client *api.Client
}

func NewOllamaClient(endpoint string) *OllamaClient {
	// url, err := url.Parse(endpoint)
	// if err != nil {
	// 	panic(err)
	// }
	// c := api.NewClient(
	// 	url,
	// 	&http.Client{Transport: NewAddHeaderTransport(nil)},
	// )
	c, err := api.ClientFromEnvironment()
	if err != nil {
		panic(err)
	}

	return &OllamaClient{
		client: c,
	}
}

func (c *OllamaClient) Ask(ctx context.Context, prompt string) (string, error) {
	messages := append(
		System(),
		api.Message{
			Role:    "user",
			Content: prompt,
		},
	)

	// Dont stream
	stream := false
	req := &api.ChatRequest{
		Model:    MODEL,
		Messages: messages,
		Stream:   &stream,
		Tools:    Tools(),
	}

	// Append messages to response
	var response string
	done := false
	respFunc := func(resp api.ChatResponse) error {
		fmt.Printf("\n\nResp: %+v\n\n", resp)
		response += resp.Message.Content
		done = resp.Done
		//req.Messages = append(req.Messages, resp.Message)

		if len(resp.Message.ToolCalls) > 0 {
			done = false
			for _, call := range resp.Message.ToolCalls {
				req.Messages = append(req.Messages, CallTool(call))
			}
		}
		return nil
	}

	// Wait for Done
	for !done {
		// Create chat request
		err := c.client.Chat(ctx, req, respFunc)
		if err != nil {
			return "", err
		}
	}

	return response, nil
}

// Set Host header for ollama
type AddHeaderTransport struct {
	T http.RoundTripper
}

func (adt *AddHeaderTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	req.Header.Add("Host", "localhost")
	req.Header.Add("Origin", "localhost")
	resp, err := adt.T.RoundTrip(req)
	fmt.Printf("\nReq: %+v\nResp: %+v\n%+v\n", req, resp, err)
	return resp, err
}

func NewAddHeaderTransport(T http.RoundTripper) *AddHeaderTransport {
	if T == nil {
		T = http.DefaultTransport
	}
	return &AddHeaderTransport{T}
}
