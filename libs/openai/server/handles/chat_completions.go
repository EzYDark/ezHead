package handles

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/ezydark/ezHead/libs/perplexity"
	"github.com/ezydark/ezHead/libs/perplexity/request"
	"github.com/rs/zerolog/log"
	"github.com/sashabaranov/go-openai"
)

var Chunks []string

func HandleChatCompletions(w http.ResponseWriter, r *http.Request) {
	// log.Debug().Msgf("Got request:\n%v", r)

	if r.Method == http.MethodOptions {
		log.Debug().Msg("Handling OPTIONS request")
		w.WriteHeader(http.StatusOK)
		return
	}

	if r.Method != http.MethodPost {
		log.Error().Msg("Method not allowed")
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Parse the req body into OpenAI's type
	var req openai.ChatCompletionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Error().Msgf("Invalid request format:\n%v", err)
		http.Error(w, "Invalid request format", http.StatusBadRequest)
		return
	}

	// Initialize Perplexity struct
	perplex, err := perplexity.Init()
	if err != nil {
		log.Fatal().Msgf("Could not initialize Perplexity struct:\n%v", err)
	}

	// "model:claude-3.7-thinking//provider:perplexity//search:web--academic--social"
	if req.Model == "" {
		msg := "Model is required (with optional `//` parameters)"
		log.Error().Msg(msg)
		http.Error(w, msg, http.StatusBadRequest)
		return
	}
	if strings.Contains(req.Model, "//") {
		sections := strings.SplitSeq(req.Model, "//")

		for section := range sections {
			if !strings.Contains(section, ":") {
				msg := "Provided model string has invalid format"
				log.Error().Msg(msg)
				http.Error(w, msg, http.StatusBadRequest)
				return
			}

			parts := strings.SplitN(section, ":", 2)
			key := parts[0]
			value := parts[1]

			switch key {
			case "model":
				perplex.ReqBody.Params.ModelPreference = getModelName(value)
			case "provider":
				switch value {
				case string(request.Perplexity):
					// TODO
				case string(request.OpenAI):
					// TODO
				}
			case "search":
				searchItems := strings.Split(value, "--")
				for _, item := range searchItems {
					switch item {
					case "web":
						perplex.ReqBody.Params.Sources = append(perplex.ReqBody.Params.Sources, request.Web)
					case "academic":
						perplex.ReqBody.Params.Sources = append(perplex.ReqBody.Params.Sources, request.Academic)
					case "social":
						perplex.ReqBody.Params.Sources = append(perplex.ReqBody.Params.Sources, request.Social)
					}
				}
			}
		}
	} else {
		// If no "//" in the model string, probably just provided the model name
		perplex.ReqBody.Params.ModelPreference = getModelName(req.Model)
	}

	err = perplex.SendRequest(perplex.RodPage, req.Messages[len(req.Messages)-1].Content)
	if err != nil {
		log.Error().Msgf("Failed to send request to Perplexity API:\n%v", err)
		return
	}

	// pretty, err := json.MarshalIndent(req, "", "    ")
	// if err != nil {
	// 	log.Fatal().Msgf("Failed to marshal request 'pretty':\n%v", err)
	// }
	// log.Debug().Msgf("Request parsed:\n%v", string(pretty))

	// prettyMessages, err := json.MarshalIndent(req.Messages, "", "    ")
	// if err != nil {
	// 	log.Fatal().Msgf("Failed to marshal request 'prettyMessages':\n%v", err)
	// }
	// log.Debug().Msgf("Request Messages:\n%v", string(prettyMessages))

	if req.Stream {
		log.Debug().Msg("Continuing streaming the response")
		handleStreamingResponse(w, req)
		return
	}

	// Process the request with your custom logic
	// This is where you'd connect to your actual LLM implementation
	response := ProcessChat(req)

	// Return the response using OpenAI's response type
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		log.Fatal().Msgf("Failed to encode response:\n%v", err)
	}

	log.Debug().Msgf("Sent successfully a response")
}

func getModelName(reqModel string) request.Models {
	switch reqModel {
	case "claude-3.7-thinking":
		return request.Claude_3_7_Thinking
	case "claude-3.7":
		return request.Claude_3_7
	case "gemini-2.5-pro":
		return request.Gemini_2_5_Pro
	case "grok-3":
		return request.Grok_3
	case "o4-mini":
		return request.O4_Mini
	case "r1-1776":
		return request.R1_1776
	case "gpt-4.1":
		return request.GPT_4_1
	case "sonar":
		return request.Sonar
	case "best":
		return request.Best
	case "auto":
		return request.Best // The same as the `best` model. Just for convenience.
	default:
		log.Info().Msg("No model specified. Using provider`s default model.")
		return request.Best
	}
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
					Content: "This is a response from your custom NON-streaming implementation",
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

	Chunks = nil

	return response
}

func handleStreamingResponse(w http.ResponseWriter, req openai.ChatCompletionRequest) {
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")

	flusher, ok := w.(http.Flusher)
	if !ok {
		log.Error().Msg("Streaming not supported")
		http.Error(w, "Streaming not supported", http.StatusInternalServerError)
		return
	}

	// Generate a response ID
	uniqueID := fmt.Sprintf("chatcmpl-%d", time.Now().UnixNano())

	chunkIndex := 0
	for {
		if len(Chunks) == 0 || chunkIndex >= len(Chunks) {
			log.Debug().Msg("No chunks available. Waiting for new chunks...")
			continue
		} else if chunkIndex == len(Chunks)-1 && Chunks[len(Chunks)-1] == "[END]" {
			log.Debug().Msg("All chunks processed")
			Chunks = nil
			break
		}

		streamEvent := openai.ChatCompletionStreamResponse{
			ID:      uniqueID,
			Object:  "chat.completion.chunk",
			Created: time.Now().Unix(),
			Model:   req.Model,
			Choices: []openai.ChatCompletionStreamChoice{
				{
					Index: 0,
					Delta: openai.ChatCompletionStreamChoiceDelta{
						Content: Chunks[chunkIndex],
					},
					FinishReason: openai.FinishReasonNull,
				},
			},
		}

		if chunkIndex == len(Chunks)-2 && Chunks[len(Chunks)-1] == "[END]" {
			streamEvent.Choices[0].FinishReason = openai.FinishReasonStop
		}

		data, err := json.Marshal(streamEvent)
		if err != nil {
			log.Fatal().Msgf("Failed to marshal streaming response:\n%v", err)
		}

		// Write SSE format
		_, err = fmt.Fprintf(w, "data: %s\n\n", data)
		if err != nil {
			log.Fatal().Msgf("Failed to write streaming response:\n%v", err)
		}
		flusher.Flush()

		chunkIndex++
	}

	// Send the [DONE] message
	_, err := fmt.Fprintf(w, "data: [DONE]\n\n")
	if err != nil {
		log.Fatal().Msgf("Failed to write streaming response:\n%v", err)
	}
	flusher.Flush()
}
