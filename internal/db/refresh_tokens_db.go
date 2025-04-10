package db

import (
	"chirpy/internal/config"
	"chirpy/internal/database"
	"chirpy/internal/model"
	"context"
	"fmt"
)

func CreateRefreshToken(cfg *config.ApiConfig, refresh_token model.RefreshToken) error {
	ctx := context.Background()

	_, err := cfg.Db.CreateRefreshToken(ctx, database.CreateRefreshTokenParams{
		Token:     refresh_token.Token,
		UserID:    refresh_token.UserID,
		ExpiresAt: refresh_token.Expires_at,
	})
	if err != nil {
		return fmt.Errorf("error creating new refresh_token %w: ", err)
	}

	return nil
}

func GetRefreshTokenDB(cfg *config.ApiConfig, refresh_token model.RefreshToken) (model.RefreshToken, error) {
	ctx := context.Background()

	token, err := cfg.Db.GetRefreshToken(ctx, refresh_token.Token)
	if err != nil {
		return model.RefreshToken{}, fmt.Errorf("error retrieving token :%w ", err)
	}
	fmt.Println(token.RevokedAt)
	return model.RefreshToken{
		Token:      token.Token,
		UserID:     token.UserID,
		Expires_at: token.ExpiresAt,
		Revoked_at: token.RevokedAt,
	}, nil
}

func RevokeRefreshTokenDB(cfg *config.ApiConfig, refresh_token model.RefreshToken) error {
	ctx := context.Background()

	err := cfg.Db.UpdateRefreshToken(ctx, refresh_token.Token)
	if err != nil {
		return fmt.Errorf("error retrieving token :%w ", err)
	}
	return nil
}
