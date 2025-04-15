package request

import (
	"encoding/json"
	"fmt"

	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
)

func (s *Script) Update(headers *Headers, body *Body) *Script {
	headersJSON, err := headers.IntoJSON()
	if err != nil {
		log.Fatal().Msgf("Failed to convert headers to JSON: %v", err)
	}
	bodyJSON, err := body.IntoJSON()
	if err != nil {
		log.Fatal().Msgf("Failed to convert body to JSON: %v", err)
	}

	jscript := fmt.Sprintf(`
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
	rs := Script(jscript)
	return &rs
}

func (h *Headers) Default() *Headers {
	h.Accept = "text/event-stream"
	h.AcceptLanguage = "cs,en-US;q=0.9,en;q=0.8"
	h.ContentType = "application/json"
	h.Priority = "u=1, i"
	h.SecChUa = "\"Microsoft Edge\";v=\"135\", \"Not-A.Brand\";v=\"8\", \"Chromium\";v=\"135\""
	h.SecChUaArch = "\"x86\""
	h.SecChUaBitness = "\"64\""
	h.SecChUaFullVersion = "\"135.0.3179.73\""
	h.SecChUaFullVersionList = "\"Microsoft Edge\";v=\"135.0.3179.73\", \"Not-A.Brand\";v=\"8.0.0.0\", \"Chromium\";v=\"135.0.7049.85\""
	h.SecChUaMobile = "?0"
	h.SecChUaModel = "\"\""
	h.SecChUaPlatform = "\"Windows\""
	h.SecChUaPlatformVersion = "\"19.0.0\""
	h.SecFetchDest = "empty"
	h.SecFetchMode = "cors"
	h.SecFetchSite = "same-origin"

	return h
}

func (b *Body) Default() *Body {
	b.Params = Params{
		Attachments:                []string{},
		Language:                   "en-US",
		Timezone:                   "Europe/Prague",
		SearchFocus:                "internet",
		Sources:                    []string{"web", "scholar", "social"},
		FrontendUUID:               "7eb6747b-9df3-4821-be2e-40c2fe5e7476", // Some random UUID (Unimportant)
		Mode:                       "copilot",
		ModelPreference:            "gemini2flash",
		SearchRecencyFilter:        nil,
		IsRelatedQuery:             false,
		IsSponsored:                false,
		VisitorID:                  "7ddb064a-37a6-4058-92f4-f94000fb00cf", // Some random UUID (Unimportant)
		UserNextAuthID:             "76732ce2-a124-40cb-a0ce-db657d4344b9", // Some random UUID (Unimportant)
		FrontendContextUUID:        "4dbd1a93-30c5-4169-b6df-4db9effe7465", // Some random UUID (Unimportant)
		PromptSource:               "user",
		QuerySource:                "home",
		BrowserHistorySummary:      []string{},
		IsIncognito:                false,
		UseSchematizedAPI:          true,
		SendBackTextInStreamingAPI: false,
		SupportedBlockUseCases: []string{
			"answer_modes", "media_items", "knowledge_cards", "inline_entity_cards",
			"place_widgets", "finance_widgets", "sports_widgets", "shopping_widgets",
			"jobs_widgets", "search_result_widgets", "entity_list_answer", "todo_list",
			"clarification_responses",
		},
		ClientCoordinates:        nil,
		IsNavSuggestionsDisabled: false,
		Version:                  "2.18",
	}
	b.QueryStr = "This important query is good :)" // [!] Message sent to the chat session

	return b
}

func (b *Body) ToFollowup(ChatUUID string) *Body {
	readWriteToken := uuid.NewString()
	targetCollectionUUID := uuid.NewString()
	followupSource := "link"

	b.Params.LastBackendUUID = &ChatUUID
	b.Params.ReadWriteToken = &readWriteToken
	b.Params.SearchRecencyFilter = nil
	b.Params.TargetCollectionUUID = &targetCollectionUUID
	b.Params.QuerySource = "followup"
	b.Params.FollowupSource = &followupSource

	return b
}

func (b *Body) IsFollowup() bool {
	// LastBackendUUID is main chat session identifier
	// Other parameters not so important to check
	return b.Params.LastBackendUUID != nil
}

func (h *Headers) IntoJSON() ([]byte, error) {
	JSON, err := json.Marshal(&h)
	if err != nil {
		log.Fatal().Msgf("Error marshaling headers: %v", err)
	}
	return JSON, nil
}

func (b *Body) IntoJSON() ([]byte, error) {
	JSON, err := json.Marshal(&b)
	if err != nil {
		log.Fatal().Msgf("Error marshaling body: %v", err)
	}
	return JSON, nil
}
