package domain

type Session struct {
	UserID       int
	AccessToken  string
	RefreshToken string
}
