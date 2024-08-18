package response

type MetaVo struct {
	Title   string `json:"title,omitempty"`
	Icon    string `json:"icon,omitempty"`
	NoCache bool   `json:"noCache,omitempty"`
	Link    string `json:"link,omitempty"`
}
