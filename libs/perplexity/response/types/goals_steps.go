package types

// Goal represents a goal in a plan
type Goal struct {
	Description    string `json:"description"`
	Final          bool   `json:"final,omitempty"`
	ID             string `json:"id"`
	TodoTaskStatus string `json:"todo_task_status,omitempty"`
}

// Step represents a step in a plan
type Step struct {
	Assets              []interface{}        `json:"assets"`
	StepType            string               `json:"step_type"`
	UUID                string               `json:"uuid"`
	InitialQueryContent *InitialQueryContent `json:"initial_query_content,omitempty"`
	SearchWebContent    *SearchWebContent    `json:"search_web_content,omitempty"`
	WebResultsContent   *WebResultsContent   `json:"web_results_content,omitempty"`
	TerminateContent    *TerminateContent    `json:"terminate_content,omitempty"`
}

// InitialQueryContent represents the initial query in a step
type InitialQueryContent struct {
	Query string `json:"query"`
}
