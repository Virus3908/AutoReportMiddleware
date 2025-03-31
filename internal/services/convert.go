package services

import "main/internal/database"

type ConvertService struct {
	DB database.Database
}

func NewConvertService(db database.Database) *ConvertService{
	return &ConvertService{
		DB: db,
	}
}

// func (s *ConvertService)