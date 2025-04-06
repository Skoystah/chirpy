package db

import (
	"chirpy/internal/config"
	"chirpy/internal/database"
	"chirpy/internal/model"
	"context"
	"fmt"
)

func CreateChirpDB(cfg *config.ApiConfig, chirp model.Chirp) (model.Chirp, error) {
	ctx := context.Background()

	newChirp, err := cfg.Db.CreateChirp(ctx, database.CreateChirpParams{Body: chirp.Body, UserID: chirp.UserID})
	if err != nil {
		return model.Chirp{}, fmt.Errorf("error creating new chirp %w: ", err)
	}

	return model.Chirp{
		ID:        newChirp.ID,
		CreatedAt: newChirp.CreatedAt,
		UpdatedAt: newChirp.UpdatedAt,
		Body:      newChirp.Body,
		UserID:    newChirp.UserID,
	}, nil
}

func GetChirpsDB(cfg *config.ApiConfig) ([]model.Chirp, error) {
	ctx := context.Background()

	chirps, err := cfg.Db.GetChirps(ctx)
	if err != nil {
		return []model.Chirp{}, fmt.Errorf("error creating new chirp %w: ", err)
	}

	//todo - dto 'mapper' for this?
	var fetchedChirps []model.Chirp

	for _, chirp := range chirps {
		fetchedChirps = append(fetchedChirps, model.Chirp{
			ID:        chirp.ID,
			CreatedAt: chirp.CreatedAt,
			UpdatedAt: chirp.UpdatedAt,
			Body:      chirp.Body,
			UserID:    chirp.UserID,
		})
	}

	return fetchedChirps, nil
}

func GetChirpDB(cfg *config.ApiConfig, chirp model.Chirp) (model.Chirp, error) {
	ctx := context.Background()

	fetchedChirp, err := cfg.Db.GetChirp(ctx, chirp.ID)
	if err != nil {
		return model.Chirp{}, fmt.Errorf("error creating new chirp %w: ", err)
	}

	return model.Chirp{
		ID:        fetchedChirp.ID,
		CreatedAt: fetchedChirp.CreatedAt,
		UpdatedAt: fetchedChirp.UpdatedAt,
		Body:      fetchedChirp.Body,
		UserID:    fetchedChirp.UserID,
	}, nil
}
