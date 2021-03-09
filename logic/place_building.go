package logic

import (
	"abbysoft/gardarike-online/model"
	rpc "abbysoft/gardarike-online/rpc/generated"
	log "github.com/sirupsen/logrus"
)

func (s *SimpleLogic) PlaceBuilding(session *PlayerSession, request *rpc.PlaceBuildingRequest) (*rpc.PlaceBuildingResponse, model.Error) {
	s.log.WithFields(log.Fields{
		"sessionID":  session.SessionID,
		"buildingID": request.BuildingID,
		"townID":     request.TownID,
		"location":   request.Location,
	})

	building, found := model.Buildings[request.BuildingID]
	if !found {
		s.log.WithField("buildingID", request.BuildingID).Error("Failed to find building")
		return nil, model.ErrBadRequest
	}

	if err := session.Tx.AddTownBuilding(request.TownID, building); err != nil {
		s.log.WithError(err).Error("Failed to add town building")
		return nil, model.ErrInternalServerError
	}

	char := session.SelectedCharacter

	if !char.Resources.Subtract(building.Cost) {
		return nil, model.ErrNotEnoughResources
	}

	if err := session.Tx.AddResourcesOrUpdate(char.ID, char.Resources); err != nil {
		s.log.WithError(err).Error("Failed to update character resources")
		return nil, model.ErrInternalServerError
	}

	return &rpc.PlaceBuildingResponse{}, nil
}
