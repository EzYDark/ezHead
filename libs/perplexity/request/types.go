package request

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

type Params struct {
	LastBackendUUID            *string  `json:"last_backend_uuid,omitempty"`
	ReadWriteToken             *string  `json:"read_write_token,omitempty"`
	Attachments                []string `json:"attachments"`
	Language                   string   `json:"language"`
	Timezone                   string   `json:"timezone"`
	SearchFocus                string   `json:"search_focus"`
	Sources                    Sources  `json:"sources"`
	SearchRecencyFilter        *string  `json:"search_recency_filter,omitempty"`
	FrontendUUID               string   `json:"frontend_uuid"`
	Mode                       string   `json:"mode"`
	TargetCollectionUUID       *string  `json:"target_collection_uuid,omitempty"`
	ModelPreference            Models   `json:"model_preference"`
	IsRelatedQuery             bool     `json:"is_related_query"`
	IsSponsored                bool     `json:"is_sponsored"`
	VisitorID                  string   `json:"visitor_id"`
	UserNextAuthID             string   `json:"user_nextauth_id"`
	FrontendContextUUID        string   `json:"frontend_context_uuid"`
	PromptSource               string   `json:"prompt_source"`
	QuerySource                string   `json:"query_source"`
	BrowserHistorySummary      []string `json:"browser_history_summary"`
	IsIncognito                bool     `json:"is_incognito"`
	UseSchematizedAPI          bool     `json:"use_schematized_api"`
	SendBackTextInStreamingAPI bool     `json:"send_back_text_in_streaming_api"`
	SupportedBlockUseCases     []string `json:"supported_block_use_cases"`
	ClientCoordinates          any      `json:"client_coordinates"`
	IsNavSuggestionsDisabled   bool     `json:"is_nav_suggestions_disabled"`
	FollowupSource             *string  `json:"followup_source,omitempty"`
	Version                    string   `json:"version"`
}

type Body struct {
	Params   Params `json:"params"`
	QueryStr string `json:"query_str"`
}

type Script string

// Perplexity`s provided models with names used in their API calls
type Models string

const (
	Claude_3_7_Thinking Models = "claude37sonnetthinking"
	Claude_3_7          Models = "claude2"
	Gemini_2_5_Pro      Models = "gemini2flash"
	Grok_3              Models = "grok"
	O4_Mini             Models = "o4mini"
	R1_1776             Models = "r1"
	GPT_4_1             Models = "gpt41"
	Sonar               Models = "experimental"
	Best                Models = "pplx_pro"
)

// TEMP: move to proper seperate file from here
// Supported LLM providers
type Providers string

const (
	Perplexity Providers = "perplexity"
	OpenAI     Providers = "openai" // TODO
)

// Perplexity`s provided search sources with names used in their API calls
type SourceType string
type Sources []SourceType

const (
	Web      SourceType = "web"
	Academic SourceType = "scholar"
	Social   SourceType = "social"
)
