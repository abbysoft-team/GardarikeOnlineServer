package logic

import (
	"abbysoft/gardarike-online/model"
	"abbysoft/gardarike-online/model/consts"
	rpc "abbysoft/gardarike-online/rpc/generated"
	log "github.com/sirupsen/logrus"
	"math/rand"
)

func (s *SimpleLogic) PlaceTown(
	session *PlayerSession, request *rpc.PlaceTownRequest) (*rpc.PlaceTownResponse, model.Error) {
	s.log.WithFields(log.Fields{
		"sessionID": session.SessionID,
		"location":  request.Location,
		"name":      request.Name,
	}).Info("PlaceTown")

	if request.Location == nil {
		request.Location = &rpc.Vector2D{
			X: rand.Float32() * consts.MapChunkSize,
			Y: rand.Float32() * consts.MapChunkSize,
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

	if err := s.db.AddTown(town, true); err != nil {
		s.log.WithError(err).Error("Failed to add town")
		return nil, model.ErrInternalServerError
	}

	session.SelectedCharacter.Towns = append(session.SelectedCharacter.Towns, town)
	return &rpc.PlaceTownResponse{
		Location: request.Location,
	}, nil
}
