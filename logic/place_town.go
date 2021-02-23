package logic

import (
	"abbysoft/gardarike-online/model"
	"abbysoft/gardarike-online/model/consts"
	rpc "abbysoft/gardarike-online/rpc/generated"
	log "github.com/sirupsen/logrus"
	"math/rand"
)

// canPlaceTown - checks if the character can place one more town.
func canPlaceTown(character model.Character) bool {
	townCount := len(character.Towns)

	return character.Resources.IsEnough(model.ResourcesPlaceTown) &&
		character.CurrentPopulation >= uint64(townCount*500)
}

func (s *SimpleLogic) getMapHeightAt(x, y int) float32 {
	return s.GameMap.Data[y+x*s.config.ChunkSize]
}

func (s *SimpleLogic) PlaceTown(
	session *PlayerSession, request *rpc.PlaceTownRequest) (*rpc.PlaceTownResponse, model.Error) {
	s.log.WithFields(log.Fields{
		"sessionID": session.SessionID,
		"location":  request.Location,
		"name":      request.Name,
	}).Info("PlaceTown")

	// First town is free but other cost money
	if len(session.SelectedCharacter.Towns) > 0 {
		if !canPlaceTown(*session.SelectedCharacter) {
			return nil, model.ErrNotEnoughResources
		}
	}

	if request.Location == nil {
		request.Location = &rpc.Vector2D{
			X: rand.Float32() * float32(s.config.ChunkSize),
			Y: rand.Float32() * float32(s.config.ChunkSize),
		}
	} else {
		if request.Location.X > float32(s.config.ChunkSize) ||
			request.Location.Y > float32(s.config.ChunkSize) ||
			request.Location.X < 0 ||
			request.Location.Y < 0 {
			s.log.Error("PlaceTown: incorrect location")
			return nil, model.ErrBadRequest
		}

		if s.getMapHeightAt(int(request.Location.X), int(request.Location.Y)) < s.config.WaterLevel {
			s.log.Error("PlaceTown: trying to place town bellow the water level")
			return nil, model.ErrBadRequest
		}
	}

	if request.Name == "" {
		err := model.NewError("Field 'Name' must be filled", rpc.Error_BAD_REQUEST)
		s.log.WithError(err).Error("Failed to add town")
		return nil, err
	}

	town := model.Town{
		X:          int64(request.Location.X),
		Y:          int64(request.Location.Y),
		OwnerName:  session.SelectedCharacter.Name,
		Population: 0,
		Name:       request.Name,
	}

	session.SelectedCharacter.MaxPopulation += consts.TownPopulationBonus
	session.SelectedCharacter.Resources.Subtract(model.ResourcesPlaceTown)

	if err := s.db.AddTown(town, false); err != nil {
		s.log.WithError(err).Error("Failed to add town")
		return nil, model.ErrInternalServerError
	}

	if err := s.db.UpdateCharacter(*session.SelectedCharacter, false); err != nil {
		s.log.WithError(err).Error("Failed to update character")
		return nil, model.ErrInternalServerError
	}

	if err := s.db.AddResourcesOrUpdate(session.SelectedCharacter.Resources, true); err != nil {
		s.log.WithError(err).Error("Failed to update character resources")
		return nil, model.ErrInternalServerError
	}

	session.SelectedCharacter.Towns = append(session.SelectedCharacter.Towns, town)
	return &rpc.PlaceTownResponse{
		Location: request.Location,
	}, nil
}
