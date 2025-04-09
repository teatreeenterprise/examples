package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
)

type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type RequestPayload struct {
	Model    string    `json:"model"`
	Messages []Message `json:"messages"`
	Stream   bool      `json:"stream"`
}

type ResponseChunk struct {
	Message struct {
		Content string `json:"content"`
	} `json:"message"`
}

func SimpleChat(model, host string) {
	fmt.Printf("Starting chat with %s. Type 'exit' to end the conversation.\n", model)

	messages := []Message{}

	scanner := bufio.NewScanner(os.Stdin)

	for {
		fmt.Print("\nYou: ")
		if !scanner.Scan() {
			break
		}
		userInput := scanner.Text()

		if strings.ToLower(userInput) == "exit" || strings.ToLower(userInput) == "quit" || strings.ToLower(userInput) == "bye" {
			fmt.Println("Goodbye!")
			break
		}

		messages = append(messages, Message{
			Role:    "user",
			Content: userInput,
		})

		payload := RequestPayload{
			Model:    model,
			Messages: messages,
			Stream:   true,
		}

		payloadBytes, err := json.Marshal(payload)
		if err != nil {
			fmt.Printf("Error creating request: %v\n", err)
			continue
		}

		client := &http.Client{}
		req, err := http.NewRequest("POST", host+"/", bytes.NewBuffer(payloadBytes))
		if err != nil {
			fmt.Printf("Error creating request: %v\n", err)
			continue
		}
		req.Header.Set("Content-Type", "application/json")

		fmt.Print("\nAssistant: ")

		resp, err := client.Do(req)
		if err != nil {
			fmt.Printf("Error making request: %v\n", err)
			continue
		}
		defer resp.Body.Close()

		if resp.StatusCode != 200 {
			bodyBytes, _ := io.ReadAll(resp.Body)
			fmt.Printf("Error: HTTP %d - %s\n", resp.StatusCode, string(bodyBytes))
			continue
		}

		assistantResponse := ""
		reader := bufio.NewReader(resp.Body)

		for {
			line, err := reader.ReadBytes('\n')
			if err != nil {
				if err != io.EOF {
					fmt.Printf("Error reading response: %v\n", err)
				}
				break
			}

			if len(line) == 0 {
				continue
			}

			var chunk ResponseChunk
			if err := json.Unmarshal(line, &chunk); err != nil {
				continue
			}

			contentChunk := chunk.Message.Content
			assistantResponse += contentChunk
			fmt.Print(contentChunk)
		}

		fmt.Println()

		messages = append(messages, Message{
			Role:    "assistant",
			Content: assistantResponse,
		})
	}
}

func main() {
	model := "model"
	host := "https://deployment-id.region.enterprise.teatree.chat"

	if len(os.Args) > 1 {
		model = os.Args[1]
	}
	if len(os.Args) > 2 {
		host = os.Args[2]
	}

	SimpleChat(model, host)
}
