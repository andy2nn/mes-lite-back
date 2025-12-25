package user

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type AuthClaims struct {
	ID   int64  `json:"id"`
	Role string `json:"role"`
	jwt.RegisteredClaims
}

type AuthService struct {
	repo      Repository
	rtRepo    RefreshTokenRepository
	jwtSecret string
	tokenTTL  time.Duration
}

func NewAuthService(
	repo Repository,
	rtRepo RefreshTokenRepository,
	secret string,
	ttl time.Duration,
) *AuthService {
	return &AuthService{
		repo:      repo,
		rtRepo:    rtRepo,
		jwtSecret: secret,
		tokenTTL:  ttl,
	}
}

func (a *AuthService) Authenticate(username, password string) (string, string, *User, error) {
	u, err := a.repo.GetByUsername(username)
	if err != nil || u == nil {
		return "", "", nil, ErrInvalidCreds
	}

	if bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password)) != nil {
		return "", "", nil, ErrInvalidCreds
	}

	access, err := a.newAccessToken(u)
	if err != nil {
		return "", "", nil, err
	}

	refresh, err := a.newRefreshToken(u.ID)
	if err != nil {
		return "", "", nil, err
	}

	return access, refresh, u, nil
}

func (a *AuthService) Refresh(refreshToken string) (string, string, error) {
	rt, err := a.rtRepo.Get(refreshToken)
	if err != nil {
		return "", "", fmt.Errorf("invalid refresh token")
	}

	if rt.ExpiresAt.Before(time.Now()) {
		_ = a.rtRepo.Delete(refreshToken)
		return "", "", fmt.Errorf("refresh token expired")
	}

	u, err := a.repo.GetByID(rt.UserID)
	if err != nil {
		return "", "", err
	}

	access, err := a.newAccessToken(u)
	if err != nil {
		return "", "", err
	}

	_ = a.rtRepo.Delete(refreshToken)
	newRefresh, err := a.newRefreshToken(u.ID)
	if err != nil {
		return "", "", err
	}

	return access, newRefresh, nil
}

func (a *AuthService) newAccessToken(u *User) (string, error) {
	now := time.Now()
	claims := AuthClaims{
		ID:   u.ID,
		Role: resolveRoleName(u.RoleID),
		RegisteredClaims: jwt.RegisteredClaims{
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(a.tokenTTL)),
		},
	}

	return jwt.NewWithClaims(jwt.SigningMethodHS256, claims).
		SignedString([]byte(a.jwtSecret))
}

func (a *AuthService) newRefreshToken(userID int64) (string, error) {
	token := uuid.New().String()

	rt := &RefreshToken{
		UserID:    userID,
		Token:     token,
		ExpiresAt: time.Now().Add(30 * 24 * time.Hour),
	}

	if err := a.rtRepo.Save(rt); err != nil {
		return "", err
	}

	return token, nil
}

func resolveRoleName(roleID int64) string {
	switch roleID {
	case 1:
		return "admin"
	case 2:
		return "operator"
	case 3:
		return "viewer"
	default:
		return "user"
	}
}
