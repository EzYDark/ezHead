package types

// Response represents the main response structure
type Response struct {
	AnswerModes         []AnswerMode  `json:"answer_modes,omitempty"`
	Attachments         []interface{} `json:"attachments"`
	BackendUUID         string        `json:"backend_uuid"`
	Blocks              []Block       `json:"blocks,omitempty"`
	BrowserTaskMessage  bool          `json:"browser_task_message"`
	ContextUUID         string        `json:"context_uuid"`
	DisplayModel        string        `json:"display_model"`
	ExpectSearchResults string        `json:"expect_search_results"`
	FinalSSEMessage     bool          `json:"final_sse_message"`
	FrontendContextUUID string        `json:"frontend_context_uuid"`
	GPT4                bool          `json:"gpt4"`
	ImageCompletions    []interface{} `json:"image_completions"`
	IsProReasoningMode  bool          `json:"is_pro_reasoning_mode,omitempty"`
	MessageMode         string        `json:"message_mode"`
	Mode                string        `json:"mode"`
	NumSourcesDisplay   int           `json:"num_sources_display,omitempty"`
	Reconnectable       bool          `json:"reconnectable"`
	SearchFocus         string        `json:"search_focus"`
	ShouldUpsellPro     bool          `json:"should_upsell_pro,omitempty"`
	Status              string        `json:"status"`
	Text                string        `json:"text,omitempty"`
	TextCompleted       bool          `json:"text_completed"`
	ThreadURLSlug       string        `json:"thread_url_slug"`
	UUID                string        `json:"uuid"`
}
