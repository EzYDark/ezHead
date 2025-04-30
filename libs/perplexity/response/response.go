package response

import (
	"encoding/json"
	"fmt"

	"github.com/ezydark/ezHead/libs/perplexity/response/types"
	"github.com/rs/zerolog/log"
	"github.com/ysmood/gson"
)

func ProcessStreamChunk(chunk gson.JSON) (*types.Response, error) {
	log.Debug().Msg("Processing stream chunk")

	chunkJSON, err := chunk.MarshalJSON()
	if err != nil {
		return nil, fmt.Errorf("failed to marshal chunk JSON: %v", err)
	}

	// log.Debug().Msgf("Raw Chunk JSON: %s", string(chunkJSON))

	var response types.Response
	if err := json.Unmarshal(chunkJSON, &response); err != nil {
		return nil, fmt.Errorf("failed to unmarshal SSE chunk: %v", err)
	}

	if !response.FinalSSEMessage {
		for _, block := range response.Blocks {
			if block.IntendedUsage == "ask_text" {
				if block.MarkdownBlock != nil && len(block.MarkdownBlock.Chunks) > 0 {
					ChunkStorage = append(ChunkStorage, block.MarkdownBlock.Chunks...)
				}
			}
		}
	} else {
		ChunkStorage = append(ChunkStorage, "[END]")
	}

	// Log specific aspects of the response for debugging
	logResponseDetails(&response)

	return &response, nil
}

// logResponseDetails logs important aspects of the response
func logResponseDetails(response *types.Response) {
	log.Debug().
		Str("uuid", response.UUID).
		Bool("final_message", response.FinalSSEMessage).
		Int("block_count", len(response.Blocks)).
		Msg("Response details")

	// Process each block according to its type
	for i, block := range response.Blocks {
		log.Debug().
			Int("block_index", i).
			Str("intended_usage", block.IntendedUsage).
			Msg(" |-- Block details")

		switch block.IntendedUsage {
		case "ask_text":
			if block.MarkdownBlock != nil {
				log.Debug().
					Str("answer_preview", truncateString(block.MarkdownBlock.Answer, 50)).
					Str("progress", block.MarkdownBlock.Progress).
					Msg("   |-- Markdown block content")
			}
		case "web_results":
			if block.WebResultBlock != nil {
				log.Debug().
					Int("result_count", len(block.WebResultBlock.WebResults)).
					Str("progress", block.WebResultBlock.Progress).
					Msg("   |-- Web results block")
			}
		case "reasoning_plan":
			if block.ReasoningPlanBlock != nil {
				log.Debug().
					Int("goals_count", len(block.ReasoningPlanBlock.Goals)).
					Str("progress", block.ReasoningPlanBlock.Progress).
					Msg("   |-- Reasoning plan block")
			}
		case "pro_search_steps":
			if block.PlanBlock != nil {
				log.Debug().
					Int("steps_count", len(block.PlanBlock.Steps)).
					Str("progress", block.PlanBlock.Progress).
					Msg("   |-- Plan block")
			}
		}
	}
}

// Helper function to truncate long strings for logging
func truncateString(s string, maxLength int) string {
	if len(s) <= maxLength {
		return s
	}
	return s[:maxLength] + "..."
}
