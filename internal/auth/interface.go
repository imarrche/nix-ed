package auth

//go:generate mockgen -source=interface.go -destination=mock/mock.go

// Service is the interface all authorization/authentication services must implement.
type Service interface {
	AuthCodeURL(string) string
	GetAccessToken(code string) (string, error)
	GetUserInfo(token string) ([]byte, error)
}
