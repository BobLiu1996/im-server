package model

type IMLoginInfo struct {
	AccessToken string `json:"accessToken"`
}

func NewIMLoginInfo(accessToken string) *IMLoginInfo {
	return &IMLoginInfo{AccessToken: accessToken}
}
