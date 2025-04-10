package db

import (
	"chirpy/internal/config"
	"chirpy/internal/database"
	"chirpy/internal/model"
	"context"
	"fmt"
)

// to-do input User or CreateUserRequest?
func CreateUserDB(cfg *config.ApiConfig, user model.User) (model.User, error) {
	ctx := context.Background()

	newUser, err := cfg.Db.CreateUser(ctx, database.CreateUserParams{
		HashedPassword: user.HashedPassword,
		Email:          user.Email,
	})

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

func UpdateUserDB(cfg *config.ApiConfig, user model.User) (model.User, error) {
	ctx := context.Background()

	updatedUser, err := cfg.Db.UpdateUser(ctx, database.UpdateUserParams{
		Email:          user.Email,
		HashedPassword: user.HashedPassword,
		ID:             user.ID,
	})

	if err != nil {
		return model.User{}, fmt.Errorf("error updating new user %w: ", err)
	}

	return model.User{
		ID:        updatedUser.ID,
		CreatedAt: updatedUser.CreatedAt,
		UpdatedAt: updatedUser.UpdatedAt,
		Email:     updatedUser.Email,
	}, nil
}

func GetUserDB(cfg *config.ApiConfig, user model.User) (model.User, error) {
	ctx := context.Background()

	fetchedUser, err := cfg.Db.GetUserByEmail(ctx, user.Email)
	if err != nil {
		return model.User{}, fmt.Errorf("error fetching user %w: ", err)
	}

	return model.User{
		ID:             fetchedUser.ID,
		CreatedAt:      fetchedUser.CreatedAt,
		UpdatedAt:      fetchedUser.UpdatedAt,
		Email:          fetchedUser.Email,
		HashedPassword: fetchedUser.HashedPassword,
	}, nil
}
