package server

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/rs/zerolog/log"
	openai "github.com/sashabaranov/go-openai"
)

func HandleChatCompletions(w http.ResponseWriter, r *http.Request) {
	// log.Debug().Msgf("Got request:\n%v", r)

	// if r.Method == http.MethodOptions {
	// 	w.WriteHeader(http.StatusOK)
	// 	return
	// }

	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Parse the request body into OpenAI's type
	var request openai.ChatCompletionRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, "Invalid request format", http.StatusBadRequest)
		return
	}

	pretty, err := json.MarshalIndent(request, "", "    ")
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to marshal request 'pretty'")
	}
	log.Debug().Msgf("Request parsed:\n%v", string(pretty))

	prettyMessages, err := json.MarshalIndent(request.Messages, "", "    ")
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to marshal request 'prettyMessages'")
	}
	log.Debug().Msgf("Request Messages:\n%v", string(prettyMessages))

	if request.Stream {
		handleStreamingResponse(w, request)
		return
	}

	// Process the request with your custom logic
	// This is where you'd connect to your actual LLM implementation
	response := ProcessChat(request)

	// Return the response using OpenAI's response type
	w.Header().Set("Content-Type", "application/json")
	// w.Header().Set("Access-Control-Allow-Origin", "*")
	// w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
	// w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		log.Error().Err(err).Msg("Failed to encode response")
	}
	log.Debug().Msgf("Sent response")
}

func ProcessChat(request openai.ChatCompletionRequest) openai.ChatCompletionResponse {
	// Use current timestamp (Unix time in seconds)
	currentTime := time.Now().Unix()

	// Generate a more realistic ID
	uniqueID := fmt.Sprintf("chatcmpl-%d", time.Now().UnixNano())

	response := openai.ChatCompletionResponse{
		ID:      uniqueID,
		Object:  "chat.completion",
		Created: currentTime,
		Model:   request.Model,
		Choices: []openai.ChatCompletionChoice{
			{
				Index: 0,
				Message: openai.ChatCompletionMessage{
					Role:    openai.ChatMessageRoleAssistant,
					Content: "This is a response from your custom implementation",
				},
				FinishReason: openai.FinishReasonStop,
			},
		},
		Usage: openai.Usage{
			PromptTokens:     100,
			CompletionTokens: 50,
			TotalTokens:      150,
		},
	}

	return response
}

func handleStreamingResponse(w http.ResponseWriter, request openai.ChatCompletionRequest) {
	// w.Header().Set("Content-Type", "text/event-stream")
	// w.Header().Set("Cache-Control", "no-cache")
	// w.Header().Set("Connection", "keep-alive")

	flusher, ok := w.(http.Flusher)
	if !ok {
		http.Error(w, "Streaming not supported", http.StatusInternalServerError)
		return
	}

	// Generate a response ID
	uniqueID := fmt.Sprintf("chatcmpl-%d", time.Now().UnixNano())

	// Split response into chunks
	// response := "This is a response from your custom implementation"
	chunks := []string{"This is ", "a response ", "from your ", "custom ", "implementation"}

	for i, chunk := range chunks {
		// Create a streaming chunk
		streamEvent := openai.ChatCompletionStreamResponse{
			ID:      uniqueID,
			Object:  "chat.completion.chunk",
			Created: time.Now().Unix(),
			Model:   request.Model,
			Choices: []openai.ChatCompletionStreamChoice{
				{
					Index: 0,
					Delta: openai.ChatCompletionStreamChoiceDelta{
						Content: chunk,
					},
					FinishReason: openai.FinishReasonNull,
				},
			},
		}

		// For the last chunk, set the finish reason
		if i == len(chunks)-1 {
			streamEvent.Choices[0].FinishReason = openai.FinishReasonStop
		}

		// Encode to JSON
		data, err := json.Marshal(streamEvent)
		if err != nil {
			log.Error().Err(err).Msg("Failed to marshal streaming response")
			continue
		}

		// Write SSE format
		fmt.Fprintf(w, "data: %s\n\n", data)
		flusher.Flush()

		// Add a small delay to simulate real streaming
		// time.Sleep(100 * time.Millisecond)
	}

	// Send the [DONE] message
	fmt.Fprintf(w, "data: [DONE]\n\n")
	flusher.Flush()
}
