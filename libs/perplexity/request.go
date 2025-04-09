package perplexity

import (
	"encoding/json"
	"fmt"

	"github.com/rs/zerolog/log"
)

type Headers struct {
	Accept                 string `json:"accept"`
	AcceptLanguage         string `json:"accept-language"`
	ContentType            string `json:"content-type"`
	Priority               string `json:"priority"`
	SecChUa                string `json:"sec-ch-ua"`
	SecChUaArch            string `json:"sec-ch-ua-arch"`
	SecChUaBitness         string `json:"sec-ch-ua-bitness"`
	SecChUaFullVersion     string `json:"sec-ch-ua-full-version"`
	SecChUaFullVersionList string `json:"sec-ch-ua-full-version-list"`
	SecChUaMobile          string `json:"sec-ch-ua-mobile"`
	SecChUaModel           string `json:"sec-ch-ua-model"`
	SecChUaPlatform        string `json:"sec-ch-ua-platform"`
	SecChUaPlatformVersion string `json:"sec-ch-ua-platform-version"`
	SecFetchDest           string `json:"sec-fetch-dest"`
	SecFetchMode           string `json:"sec-fetch-mode"`
	SecFetchSite           string `json:"sec-fetch-site"`
}

// Define your request parameters structure
type Params struct {
	LastBackendUUID            string      `json:"last_backend_uuid"`
	ReadWriteToken             string      `json:"read_write_token"`
	Attachments                []string    `json:"attachments"`
	Language                   string      `json:"language"`
	Timezone                   string      `json:"timezone"`
	SearchFocus                string      `json:"search_focus"`
	Sources                    []string    `json:"sources"`
	FrontendUUID               string      `json:"frontend_uuid"`
	Mode                       string      `json:"mode"`
	ModelPreference            string      `json:"model_preference"`
	IsRelatedQuery             bool        `json:"is_related_query"`
	IsSponsored                bool        `json:"is_sponsored"`
	VisitorID                  string      `json:"visitor_id"`
	UserNextauthID             string      `json:"user_nextauth_id"`
	PromptSource               string      `json:"prompt_source"`
	QuerySource                string      `json:"query_source"`
	LocalSearchEnabled         bool        `json:"local_search_enabled"`
	BrowserHistorySummary      []string    `json:"browser_history_summary"`
	IsIncognito                bool        `json:"is_incognito"`
	UseSchematizedAPI          bool        `json:"use_schematized_api"`
	SendBackTextInStreamingAPI bool        `json:"send_back_text_in_streaming_api"`
	SupportedBlockUseCases     []string    `json:"supported_block_use_cases"`
	ClientCoordinates          interface{} `json:"client_coordinates"`
	IsNavSuggestionsDisabled   bool        `json:"is_nav_suggestions_disabled"`
	FollowupSource             string      `json:"followup_source"`
	Version                    string      `json:"version"`
}

// Define your request body structure
type RequestBody struct {
	Params   Params `json:"params"`
	QueryStr string `json:"query_str"`
}

func DefaultHeaders() Headers {
	return Headers{
		Accept:                 "text/event-stream",
		AcceptLanguage:         "cs,en-US;q=0.9,en;q=0.8",
		ContentType:            "application/json",
		Priority:               "u=1, i",
		SecChUa:                "\"Microsoft Edge\";v=\"135\", \"Not-A.Brand\";v=\"8\", \"Chromium\";v=\"135\"",
		SecChUaArch:            "\"x86\"",
		SecChUaBitness:         "\"64\"",
		SecChUaFullVersion:     "\"135.0.3179.54\"",
		SecChUaFullVersionList: "\"Microsoft Edge\";v=\"135.0.3179.54\", \"Not-A.Brand\";v=\"8.0.0.0\", \"Chromium\";v=\"135.0.7049.42\"",
		SecChUaMobile:          "?0",
		SecChUaModel:           "\"\"",
		SecChUaPlatform:        "\"Windows\"",
		SecChUaPlatformVersion: "\"19.0.0\"",
		SecFetchDest:           "empty",
		SecFetchMode:           "cors",
		SecFetchSite:           "same-origin",
	}
}

func DefaultBody() RequestBody {
	return RequestBody{
		Params: Params{
			LastBackendUUID:            "9f0a26f5-48c4-4c1d-b97c-49141d4c1946",
			ReadWriteToken:             "a25454be-e94d-4a63-b04e-05116abfa1d1",
			Attachments:                []string{},
			Language:                   "en-US",
			Timezone:                   "Europe/Prague",
			SearchFocus:                "internet",
			Sources:                    []string{"web"},
			FrontendUUID:               "0447638c-f6fe-4fe8-9cf6-d3b518634b45",
			Mode:                       "copilot",
			ModelPreference:            "claude37sonnetthinking",
			IsRelatedQuery:             false,
			IsSponsored:                false,
			VisitorID:                  "7ddb064a-37a6-4058-92f4-f94000fb00cf",
			UserNextauthID:             "76732ce2-a124-40cb-a0ce-db657d4344b9",
			PromptSource:               "user",
			QuerySource:                "followup",
			LocalSearchEnabled:         true,
			BrowserHistorySummary:      []string{},
			IsIncognito:                false,
			UseSchematizedAPI:          true,
			SendBackTextInStreamingAPI: false,
			SupportedBlockUseCases: []string{
				"answer_modes", "media_items", "knowledge_cards", "inline_entity_cards",
				"place_widgets", "finance_widgets", "sports_widgets", "shopping_widgets",
				"jobs_widgets", "search_result_widgets", "entity_list_answer", "todo_list",
			},
			ClientCoordinates:        nil,
			IsNavSuggestionsDisabled: false,
			FollowupSource:           "link",
			Version:                  "2.18",
		},
		QueryStr: "Default query: Give me some golang code examples",
	}
}

// Convert structs to JSON
func (h Headers) IntoJSON() ([]byte, error) {
	JSON, err := json.Marshal(h)
	if err != nil {
		log.Fatal().Msgf("Error marshaling headers: %v", err)
	}
	return JSON, nil
}

func (b RequestBody) IntoJSON() ([]byte, error) {
	JSON, err := json.Marshal(b)
	if err != nil {
		log.Fatal().Msgf("Error marshaling body: %v", err)
	}
	return JSON, nil
}

func DefaultJsScript(headersJSON, bodyJSON []byte) string {
	// backtick := "`"

	jsScript := fmt.Sprintf(`
async () => {
    // ADDED: Top-level try block to catch ANY error within the async function
    try {
        const url = "https://www.perplexity.ai/rest/sse/perplexity_ask";
        // Directly use the injected JSON strings. JS will parse them.
        const headers = %s;
        const body = %s;

        // Inner try specifically for the fetch/stream operation
        try {
            const response = await fetch(url, {
                method: "POST",
                headers: headers, // headers is now a JS object
                body: JSON.stringify(body), // body needs stringification for fetch
                credentials: "include", // Important for cookies/session
                mode: "cors",
            });

            // Check response status IN the inner try block
            if (!response.ok) {
                 // Throw detailed error to be caught by INNER catch
                throw new Error('HTTP error! Status: ' + response.status + ' ' + response.statusText);
            }

            const reader = response.body.getReader();
            const decoder = new TextDecoder();
            let finalMessage = null; // Keep track of the final message
            let buffer = "";
            // Removed 'messages' array since onlyFinal is true

            try { // Added try/finally specifically around reader loop
                while (true) {
                    const { done, value } = await reader.read();
                    if (done) break;

                    buffer += decoder.decode(value, { stream: true });
                    const lines = buffer.split("\n");
                    buffer = lines.pop() || ""; // Keep remaining partial line

                    for (const line of lines) {
                        if (line.startsWith("data: ")) {
                            const data = line.substring(6).trim();
                            if (data && data !== ":heartbeat" && data !== "[DONE]") {
                                try {
                                    const parsedData = JSON.parse(data);
                                    // Always track the latest potential final message
                                    if (parsedData.final_sse_message === true) {
                                        finalMessage = parsedData;
                                    }
                                    // You might want to process/store intermediate data here
                                    // if needed even when onlyFinal=true, e.g., for debugging.
                                } catch (parseError) {
                                    console.warn("Failed to parse SSE JSON:", parseError.message, "Data:", data);
                                    // Decide if a parse error should stop everything
                                    // throw new Error("SSE JSON parse error: " + parseError.message);
                                }
                            }
                        }
                    }
                }
            } finally {
                 // Ensure reader is released even if loop errors out
                 if (reader) {
                     reader.releaseLock();
                 }
            }


            // Return success object from the INNER try block
            return { success: true, finalMessage: finalMessage }; // Simplified for onlyFinal=true

        } catch (fetchOrStreamError) {
             // This catches errors from fetch, response.ok check, reader, or parsing
            console.error("INNER CATCH (Fetch/Stream Error):", fetchOrStreamError.name, fetchOrStreamError.message);
            // Re-throw to be caught by the outer catch block
            throw fetchOrStreamError;
        }

    // ADDED: Top-level catch block
    } catch (error) {
        // This catches setup errors (JSON.parse) or errors re-thrown from inner catch
        console.error("OUTER CATCH (Top-Level JS Error):", error.name, error.message);
        // Return a simple, serializable error object from the outer catch
        // Ensure error.message is included, provide default if it's missing
        return { success: false, error: error.message || 'Unknown JS execution error' };
    }
};
	`, headersJSON, bodyJSON)

	return jsScript
}
