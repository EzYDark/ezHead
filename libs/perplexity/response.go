package perplexity

import (
	"encoding/json"
	"errors"
	"time"

	"github.com/rs/zerolog/log"
)

// Response represents the top-level response structure from Perplexity API
type Response struct {
	FinalMessage FinalMessage `json:"finalMessage"` // The main response message content
	Success      bool         `json:"success"`      // Whether the request was successful
}

// FinalMessage contains the main content of the Perplexity response
type FinalMessage struct {
	Extras               Extras             `json:"_extras"`                // Additional metadata
	AccessLevel          string             `json:"access_level"`           // Access level
	AnswerModes          []AnswerMode       `json:"answer_modes"`           // Available answer modes
	Attachments          []json.RawMessage  `json:"attachments"`            // Any attached content
	AuthorID             string             `json:"author_id"`              // ID of the author
	AuthorImage          string             `json:"author_image"`           // URL to author's image
	AuthorUsername       string             `json:"author_username"`        // Username of the author
	BackendUUID          string             `json:"backend_uuid"`           // Backend identifier
	Blocks               []Block            `json:"blocks"`                 // Content blocks
	BookmarkState        string             `json:"bookmark_state"`         // Bookmark status
	BrowserTaskMessage   bool               `json:"browser_task_message"`   // Whether this is a browser task message
	ContextUUID          string             `json:"context_uuid"`           // Context identifier
	DisplayModel         string             `json:"display_model"`          // Model used for display
	EntryUpdatedDatetime string             `json:"entry_updated_datetime"` // When entry was updated
	ExpectSearchResults  string             `json:"expect_search_results"`  // Whether search results are expected
	Final                bool               `json:"final"`                  // Whether this is the final message
	FinalSSEMessage      bool               `json:"final_sse_message"`      // Whether this is the final SSE message
	FrontendContextUUID  string             `json:"frontend_context_uuid"`  // Frontend context identifier
	GPT4                 bool               `json:"gpt4"`                   // Whether GPT-4 was used
	ImageCompletions     []json.RawMessage  `json:"image_completions"`      // Image completion data
	IsProReasoningMode   bool               `json:"is_pro_reasoning_mode"`  // Whether pro reasoning mode was used
	MessageMode          string             `json:"message_mode"`           // Message mode
	Mode                 string             `json:"mode"`                   // Response mode
	NumSourcesDisplay    int                `json:"num_sources_display"`    // Number of sources to display
	Personalized         bool               `json:"personalized"`           // Whether the response is personalized
	Plan                 Plan               `json:"plan"`                   // Response plan
	PrivacyState         string             `json:"privacy_state"`          // Privacy state
	PromptSource         string             `json:"prompt_source"`          // Source of the prompt
	QuerySource          string             `json:"query_source"`           // Source of the query
	QueryStr             string             `json:"query_str"`              // The query string
	ReadWriteToken       string             `json:"read_write_token"`       // Read/write token
	ReasoningPlan        ReasoningPlan      `json:"reasoning_plan"`         // Reasoning plan
	Reconnectable        bool               `json:"reconnectable"`          // Whether connection can be reestablished
	RelatedQueries       []string           `json:"related_queries"`        // Related queries
	RelatedQueryItems    []RelatedQueryItem `json:"related_query_items"`    // Structured related query items
	S3SocialPreviewURL   string             `json:"s3_social_preview_url"`  // URL for social preview
	SearchFocus          string             `json:"search_focus"`           // Focus of the search
	ShouldUpsellPro      bool               `json:"should_upsell_pro"`      // Whether to upsell pro features
	Sources              Sources            `json:"sources"`                // Sources used for the response
	Status               string             `json:"status"`                 // Status of the response
	StepType             string             `json:"step_type"`              // Type of step
	Text                 string             `json:"text"`                   // Raw text content (contains JSON)
	ParsedSteps          []Step             `json:"-"`                      // Parsed steps (not in JSON)
	TextCompleted        bool               `json:"text_completed"`         // Whether text generation is complete
	ThreadAccess         int                `json:"thread_access"`          // Thread access level
	ThreadTitle          string             `json:"thread_title"`           // Title of the thread
	ThreadURLSlug        string             `json:"thread_url_slug"`        // URL slug for the thread
	UpdatedDatetime      string             `json:"updated_datetime"`       // When the response was updated
	UUID                 string             `json:"uuid"`                   // Unique identifier
}

// AnswerMode represents an answer mode configuration
type AnswerMode struct {
	AnswerModeType string `json:"answer_mode_type"` // Type of answer mode
}

// Block represents a content block in the response
type Block struct {
	IntendedUsage      string           `json:"intended_usage"`                 // The intended usage of the block
	PlanBlock          *json.RawMessage `json:"plan_block,omitempty"`           // Plan block content
	MarkdownBlock      *json.RawMessage `json:"markdown_block,omitempty"`       // Markdown block content
	WebResultBlock     *json.RawMessage `json:"web_result_block,omitempty"`     // Web result block content
	MediaBlock         *json.RawMessage `json:"media_block,omitempty"`          // Media block content
	SourcesModeBlock   *json.RawMessage `json:"sources_mode_block,omitempty"`   // Sources mode block content
	ReasoningPlanBlock *json.RawMessage `json:"reasoning_plan_block,omitempty"` // Reasoning plan block content
}

// Extras contains additional metadata about the request and response
type Extras struct {
	CoreElapsed   float64 `json:"core_elapsed"`    // Processing time in seconds
	Country       string  `json:"country"`         // Country code
	Next          *string `json:"next"`            // Next action
	ProSearchMode string  `json:"pro_search_mode"` // Pro search mode
	Subdomain     *string `json:"subdomain"`       // Subdomain
}

// Plan contains information about the response plan
type Plan struct {
	Final bool   `json:"final"` // Whether this is the final plan
	Goals []Goal `json:"goals"` // Goals in the plan
}

// Goal represents a single goal in the response plan
type Goal struct {
	Description    string `json:"description"`      // Description of the goal
	Final          bool   `json:"final"`            // Whether this is the final goal
	ID             string `json:"id"`               // Goal identifier
	TodoTaskStatus string `json:"todo_task_status"` // Todo task status
}

// ReasoningPlan contains information about the reasoning process
type ReasoningPlan struct {
	Final      bool            `json:"final"`                 // Whether this is the final reasoning plan
	Goals      []ReasoningGoal `json:"goals"`                 // Goals in the reasoning plan
	WebResults []WebResult     `json:"web_results,omitempty"` // Web results if present
}

// ReasoningGoal represents a single goal in the reasoning plan
type ReasoningGoal struct {
	Description string `json:"description"`  // Description of the reasoning goal
	ID          string `json:"id,omitempty"` // Goal ID if present
}

// WebResult represents a web search result
type WebResult struct {
	IsAttachment      bool     `json:"is_attachment"`       // Whether this is an attachment
	IsClientContext   bool     `json:"is_client_context"`   // Whether this is client context
	IsCodeInterpreter bool     `json:"is_code_interpreter"` // Whether this is code interpreter
	IsFocusedWeb      bool     `json:"is_focused_web"`      // Whether this is focused web
	IsImage           bool     `json:"is_image"`            // Whether this is an image
	IsKnowledgeCard   bool     `json:"is_knowledge_card"`   // Whether this is a knowledge card
	IsNavigational    bool     `json:"is_navigational"`     // Whether this is navigational
	IsWidget          bool     `json:"is_widget"`           // Whether this is a widget
	MetaData          MetaData `json:"meta_data"`           // Metadata about the result
	Name              string   `json:"name"`                // Name of the result
	Snippet           string   `json:"snippet"`             // Snippet from the result
	Timestamp         string   `json:"timestamp"`           // Timestamp
	URL               string   `json:"url"`                 // URL of the result
}

// MetaData contains metadata about a web result
type MetaData struct {
	Client        string   `json:"client"`           // Client
	Date          *string  `json:"date"`             // Date
	Description   *string  `json:"description"`      // Description
	DomainName    string   `json:"domain_name"`      // Domain name
	Images        []string `json:"images,omitempty"` // Images
	PublishedDate *string  `json:"published_date"`   // Published date
}

// RelatedQueryItem represents a suggested related query
type RelatedQueryItem struct {
	Text string `json:"text"` // Query text
	Type string `json:"type"` // Query type
}

// Sources represents sources used for the response
type Sources struct {
	Sources []string `json:"sources"` // Sources list
}

// Step represents a single step in the response processing
type Step struct {
	StepType string          `json:"step_type"` // Type of step (INITIAL_QUERY, SEARCH_WEB, etc.)
	Content  json.RawMessage `json:"content"`   // Content of the step (varies by type)
	UUID     string          `json:"uuid"`      // Unique identifier for the step
}

// Helper methods

// UnmarshalSteps parses the Text field into structured Step objects
func (fm *FinalMessage) UnmarshalSteps() error {
	err := json.Unmarshal([]byte(fm.Text), &fm.ParsedSteps)
	if err != nil {
		log.Fatal().Msgf("Failed to unmarshal steps:\n%v", err)
	}

	return nil
}

// GetFinalAnswer returns the answer text from the FINAL step if available
func (fm *FinalMessage) GetFinalAnswer() (string, error) {
	if len(fm.ParsedSteps) == 0 {
		if err := fm.UnmarshalSteps(); err != nil {
			log.Fatal().Msgf("Failed to unmarshal steps:\n%v", err)
		}
	}

	for _, step := range fm.ParsedSteps {
		if step.StepType == "FINAL" {
			var finalContent struct {
				Answer string `json:"answer"`
			}

			if err := json.Unmarshal(step.Content, &finalContent); err != nil {
				log.Fatal().Msgf("Failed to unmarshal final content:\n%v", err)
			}

			var parsedAnswer struct {
				Answer string `json:"answer"`
			}

			if err := json.Unmarshal([]byte(finalContent.Answer), &parsedAnswer); err != nil {
				log.Fatal().Msgf("Failed to unmarshal parsed answer:\n%v", err)
			}

			return parsedAnswer.Answer, nil
		}
	}

	return "", errors.New("No final step found")
}

// ParsedTime returns the UpdatedDatetime as a time.Time
func (fm *FinalMessage) ParsedTime() (time.Time, error) {
	parsedTime, err := time.Parse(time.RFC3339Nano, fm.UpdatedDatetime)
	if err != nil {
		log.Fatal().Msgf("Failed to parse time:\n%v", err)
	}

	return parsedTime, nil
}
