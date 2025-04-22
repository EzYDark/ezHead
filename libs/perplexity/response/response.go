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

	log.Debug().
		Str("uuid", response.UUID).
		Bool("final_message", response.FinalSSEMessage).
		Int("block_count", len(response.Blocks)).
		Msg("Response details")

	return &response, nil
}
