package request

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/google/uuid"
)

func (s *Script) Update(headers *Headers, body *Body) (*Script, error) {
	headersJSON, err := headers.IntoJSON()
	if err != nil {
		return nil, fmt.Errorf("failed to convert headers to JSON:\n%v", err)
	}
	bodyJSON, err := body.IntoJSON()
	if err != nil {
		return nil, fmt.Errorf("failed to convert body to JSON:\n%v", err)
	}

	reqScriptPath := "./libs/perplexity/request/req_script.js"
	reqScriptBytes, err := os.ReadFile(reqScriptPath)
	if err != nil {
		// Handle the error if the file cannot be read.
		return nil, fmt.Errorf("error reading file '%s':\n%v", reqScriptPath, err)
	}

	// Convert the byte slice to a string.
	reqScriptContent := string(reqScriptBytes)

	reqScript := fmt.Sprintf(reqScriptContent, headersJSON, bodyJSON)
	rs := Script(reqScript)
	return &rs, nil
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
		Sources:                    Sources{Web, Academic, Social},
		FrontendUUID:               "7eb6747b-9df3-4821-be2e-40c2fe5e7476", // Some random UUID (Unimportant)
		Mode:                       "copilot",
		ModelPreference:            Models(Claude_3_7_Thinking),
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
		return nil, fmt.Errorf("error marshaling headers: %v", err)
	}
	return JSON, nil
}

func (b *Body) IntoJSON() ([]byte, error) {
	JSON, err := json.Marshal(&b)
	if err != nil {
		return nil, fmt.Errorf("error marshaling body: %v", err)
	}
	return JSON, nil
}
