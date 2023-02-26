package model

type AlfredResponse struct {
	Variables struct {
		AWSESSIONID string `json:"AW_SESSION_ID"`
	} `json:"variables"`
	Items []struct {
		Title        string `json:"title"`
		Subtitle     string `json:"subtitle"`
		Autocomplete string `json:"autocomplete"`
		Arg          string `json:"arg"`
		Uid          string `json:"uid"`
		Valid        bool   `json:"valid"`
		Type         string `json:"type"`
		Icon         struct {
			Path string `json:"path"`
			Type string `json:"type"`
		} `json:"icon"`
	} `json:"items"`
}
