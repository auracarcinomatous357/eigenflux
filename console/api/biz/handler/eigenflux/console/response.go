package console

type ListAgentsDocResp struct {
	Code int32              `json:"code"`
	Msg  string             `json:"msg"`
	Data *ListAgentsDocData `json:"data"`
}

type ListAgentsDocData struct {
	Agents   []*ConsoleAgentDocInfo `json:"agents"`
	Total    int64                  `json:"total"`
	Page     int32                  `json:"page"`
	PageSize int32                  `json:"page_size"`
}

type ConsoleAgentDocInfo struct {
	AgentID         string   `json:"agent_id"`
	Email           string   `json:"email"`
	AgentName       string   `json:"agent_name"`
	Bio             string   `json:"bio"`
	CreatedAt       int64    `json:"created_at"`
	UpdatedAt       int64    `json:"updated_at"`
	ProfileStatus   *int32   `json:"profile_status,omitempty"`
	ProfileKeywords []string `json:"profile_keywords,omitempty"`
}
