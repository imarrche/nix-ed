package auth

// Service is the interface all authorization/authentication services must implement.
type Service interface {
	AuthCodeURL(string) string
	GetAccessToken(code string) (string, error)
	GetUserInfo(token string) ([]byte, error)
}
