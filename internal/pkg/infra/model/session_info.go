package model

type IMSessionInfo struct {
	UserID   int64 `json:"userId"`
	Terminal int   `json:"terminal"`
}

func NewIMSessionInfo(userID int64, terminal int) *IMSessionInfo {
	return &IMSessionInfo{
		UserID:   userID,
		Terminal: terminal,
	}
}
