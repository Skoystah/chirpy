package db

import (
	"chirpy/internal/config"
	"chirpy/internal/database"
	"chirpy/internal/model"
	"context"
	"fmt"
)

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
		ID:          newUser.ID,
		CreatedAt:   newUser.CreatedAt,
		UpdatedAt:   newUser.UpdatedAt,
		Email:       newUser.Email,
		IsChirpyRed: newUser.IsChirpyRed,
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
		ID:          updatedUser.ID,
		CreatedAt:   updatedUser.CreatedAt,
		UpdatedAt:   updatedUser.UpdatedAt,
		Email:       updatedUser.Email,
		IsChirpyRed: updatedUser.IsChirpyRed,
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
		IsChirpyRed:    fetchedUser.IsChirpyRed,
	}, nil
}

func UpdateUserRedDB(cfg *config.ApiConfig, user model.User) (model.User, error) {
	ctx := context.Background()

	updatedUser, err := cfg.Db.UpdateUserRed(ctx, database.UpdateUserRedParams{
		IsChirpyRed: user.IsChirpyRed,
		ID:          user.ID})

	if err != nil {
		return model.User{}, fmt.Errorf("error updating user red %w: ", err)
	}

	return model.User{
		ID:             updatedUser.ID,
		CreatedAt:      updatedUser.CreatedAt,
		UpdatedAt:      updatedUser.UpdatedAt,
		Email:          updatedUser.Email,
		HashedPassword: updatedUser.HashedPassword,
		IsChirpyRed:    updatedUser.IsChirpyRed,
	}, nil
}
