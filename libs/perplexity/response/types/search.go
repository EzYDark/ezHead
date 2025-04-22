package types

// WebResult represents a search result
type WebResult struct {
	IsAttachment      bool               `json:"is_attachment"`
	IsClientContext   bool               `json:"is_client_context,omitempty"`
	IsCodeInterpreter bool               `json:"is_code_interpreter,omitempty"`
	IsFocusedWeb      bool               `json:"is_focused_web"`
	IsImage           bool               `json:"is_image,omitempty"`
	IsKnowledgeCard   bool               `json:"is_knowledge_card,omitempty"`
	IsNavigational    bool               `json:"is_navigational"`
	IsWidget          bool               `json:"is_widget,omitempty"`
	MetaData          *WebResultMetaData `json:"meta_data,omitempty"`
	Name              string             `json:"name"`
	Snippet           string             `json:"snippet"`
	Timestamp         string             `json:"timestamp"`
	URL               string             `json:"url"`
}

// WebResultMetaData contains metadata for web results
type WebResultMetaData struct {
	Authors       []string `json:"authors,omitempty"`
	Client        string   `json:"client,omitempty"`
	Date          string   `json:"date,omitempty"`
	Description   string   `json:"description,omitempty"`
	DomainName    string   `json:"domain_name,omitempty"`
	Images        []string `json:"images,omitempty"`
	PublishedDate string   `json:"published_date,omitempty"`
}

type SearchWebContent struct {
	GoalID  string  `json:"goal_id"`
	Queries []Query `json:"queries"`
}

type Query struct {
	Engine string `json:"engine"`
	Limit  int    `json:"limit"`
	Query  string `json:"query"`
}

type WebResultsContent struct {
	GoalID     string      `json:"goal_id"`
	WebResults []WebResult `json:"web_results"`
}

type TerminateContent struct {
	GoalID string `json:"goal_id"`
}
