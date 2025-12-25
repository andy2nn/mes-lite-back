package user

type RefreshTokenRepository interface {
	Save(token *RefreshToken) error
	Get(token string) (*RefreshToken, error)
	Delete(token string) error
}
