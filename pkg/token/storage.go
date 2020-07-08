package token

import "context"

type Token struct {
	Token    string `json:"token"`
	Platform string `json:"platform"`
}

type Storage interface {
	SaveToken(ctx context.Context, userID uint64, token Token) error
	UserTokens(ctx context.Context, userID uint64) ([]Token, error)
	DeleteTokens(ctx context.Context, tokens []string) error
}
