package types

// Block represents a content block which can be of different types
type Block struct {
	IntendedUsage      string              `json:"intended_usage"`
	PlanBlock          *PlanBlock          `json:"plan_block,omitempty"`
	ReasoningPlanBlock *ReasoningPlanBlock `json:"reasoning_plan_block,omitempty"`
	SourcesModeBlock   *SourcesModeBlock   `json:"sources_mode_block,omitempty"`
	MarkdownBlock      *MarkdownBlock      `json:"markdown_block,omitempty"`
	WebResultBlock     *WebResultBlock     `json:"web_result_block,omitempty"`
}

// Additional supporting types
type SourcesModeBlock struct {
	AnswerModeType string      `json:"answer_mode_type"`
	Progress       string      `json:"progress"`
	ResultCount    int         `json:"result_count"`
	WebResults     []WebResult `json:"web_results,omitempty"`
}

// PlanBlock represents planning steps
type PlanBlock struct {
	Final    bool   `json:"final"`
	Goals    []Goal `json:"goals"`
	Progress string `json:"progress"`
	Steps    []Step `json:"steps,omitempty"`
}

// ReasoningPlanBlock represents a reasoning plan
type ReasoningPlanBlock struct {
	Goals    []Goal `json:"goals"`
	Progress string `json:"progress"`
}

// MarkdownBlock represents the actual answer content
type MarkdownBlock struct {
	Answer              string   `json:"answer"`
	ChunkStartingOffset int      `json:"chunk_starting_offset"`
	Chunks              []string `json:"chunks"`
	Progress            string   `json:"progress"`
}

// WebResultBlock contains search results
type WebResultBlock struct {
	Progress   string      `json:"progress"`
	WebResults []WebResult `json:"web_results"`
}
