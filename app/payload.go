package app

type Payload struct {
	Comment struct {
		Content []struct {
			Content []struct {
				Text string `json:"text"`
				Type string `json:"type"`
			} `json:"content"`
			Type string `json:"type"`
		} `json:"content"`
		Type    string `json:"type"`
		Version int    `json:"version"`
	} `json:"comment"`
	Started          string `json:"started"`
	TimeSpentSeconds int    `json:"timeSpentSeconds"`
}
