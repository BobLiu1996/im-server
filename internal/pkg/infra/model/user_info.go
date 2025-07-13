package model

type IMUserInfo struct {
	UserID   int64 `json:"userId"`
	Terminal int   `json:"terminal"`
}

func NewIMUserInfo(userID int64, terminal int) *IMUserInfo {
	return &IMUserInfo{
		UserID:   userID,
		Terminal: terminal,
	}
}
