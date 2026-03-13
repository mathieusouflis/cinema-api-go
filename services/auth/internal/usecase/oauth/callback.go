package oauthCallbackUsecase

import (
	"context"
	"crypto/sha256"
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"time"

	"authService/internal/config"
	"authService/internal/domain"

	"filmserver/pkg/errors"
	"filmserver/pkg/jwt"
	"filmserver/pkg/logger"
)

var adjectives = []string{
	"bold", "calm", "cool", "dark", "deep", "epic", "fair", "fast", "free",
	"glad", "good", "gray", "grim", "hard", "keen", "kind", "lazy", "lean",
	"loud", "mild", "neat", "nice", "pale", "pure", "rare", "rich", "safe",
	"slim", "slow", "soft", "tall", "tiny", "true", "vast", "warm", "wide",
	"wise", "wild",
}

var nouns = []string{
	"bear", "bird", "bull", "crab", "deer", "duck", "eel", "elk", "fox",
	"frog", "hawk", "jay", "lamb", "lion", "lynx", "mink", "mole", "moth",
	"mule", "newt", "puma", "rat", "seal", "slug", "swan", "toad", "vole",
	"wasp", "wolf", "wren",
}

const (
	AccessTokenTTL  = 15 * time.Minute
	RefreshTokenTTL = 7 * 24 * time.Hour
)

type Input struct {
	Email    string
	Id       string
	Provider string
}

type Output struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

type Usecase struct {
	TokenRepository domain.TokenRepository
	UserRepository  domain.UserRepository
}

func New(tokenRepository domain.TokenRepository, userRepository domain.UserRepository) *Usecase {
	return &Usecase{TokenRepository: tokenRepository, UserRepository: userRepository}
}

func (u *Usecase) Execute(ctx context.Context, input Input) (*Output, error) {
	log := logger.New(config.Load().Env)
	if input.Email == "" || input.Id == "" || input.Provider == "" {
		log.Error("OAUTHUSECASE -> Invalid input")
		return nil, errors.ErrBadRequest
	}

	var user *domain.User

	switch input.Provider {
	case "google":
		existing, err := u.UserRepository.GetByGoogleID(ctx, input.Id)
		if err == nil {
			log.Info("OAUTHUSECASE -> Found existing user by Google ID", "email", existing.Email)

			user = existing
		} else {
			byEmail, err := u.UserRepository.GetByEmail(ctx, input.Email)
			if err == nil {
				log.Info("OAUTHUSECASE -> Found existing user by Google Email", "email", existing.Email)
				if linkErr := u.UserRepository.UpdateGoogleID(ctx, byEmail.ID, input.Id); linkErr != nil {
					return nil, linkErr
				}
				user = byEmail
			} else {
				log.Info("OAUTHUSECASE -> No user fount")
				username, err := u.generateUniqueUsername(ctx, input.Id)
				if err != nil {
					return nil, err
				}
				created, err := u.UserRepository.CreateWithOAuth(ctx, domain.CreateOAuthUserInput{
					Username: username,
					Email:    input.Email,
					GoogleID: input.Id,
				})
				if err != nil {
					return nil, err
				}
				user = created
			}
		}
	default:
		return nil, errors.ErrBadRequest
	}

	claims := jwt.JWTClaims{
		UserId: user.ID,
		Email:  user.Email,
		Role:   user.Role,
	}

	accessToken, err := jwt.SignToken(claims, AccessTokenTTL)
	if err != nil {
		return nil, err
	}

	refreshToken, err := jwt.SignToken(claims, RefreshTokenTTL)
	if err != nil {
		return nil, err
	}

	if err := u.TokenRepository.StoreRefreshToken(ctx, user.ID, refreshToken, RefreshTokenTTL); err != nil {
		return nil, err
	}

	return &Output{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

func (u *Usecase) generateUniqueUsername(ctx context.Context, seed string) (string, error) {
	for attempt := range 10 {
		h := sha256.Sum256([]byte(fmt.Sprintf("%s:%d", seed, attempt)))

		adjIdx := int(binary.BigEndian.Uint16(h[0:2])) % len(adjectives)
		nounIdx := int(binary.BigEndian.Uint16(h[2:4])) % len(nouns)
		suffix := hex.EncodeToString(h[4:6])

		username := fmt.Sprintf("%s_%s%s", adjectives[adjIdx], nouns[nounIdx], suffix)

		_, err := u.UserRepository.GetByUsername(ctx, username)
		if err != nil {
			return username, nil
		}
	}
	return "", fmt.Errorf("failed to generate a unique username after 10 attempts")
}
