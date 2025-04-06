package db

import (
	"chirpy/internal/config"
	"chirpy/internal/model"
	"context"
	"fmt"
)

// todo input User or CreateUserRequest?
func CreateUserDB(cfg *config.ApiConfig, user model.User) (model.User, error) {
	ctx := context.Background()

	newUser, err := cfg.Db.CreateUser(ctx, user.Email)
	if err != nil {
		return model.User{}, fmt.Errorf("error creating new user %w: ", err)
	}

	return model.User{
		ID:        newUser.ID,
		CreatedAt: newUser.CreatedAt,
		UpdatedAt: newUser.UpdatedAt,
		Email:     newUser.Email,
	}, nil
}
