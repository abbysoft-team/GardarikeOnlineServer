package logic

import (
	"abbysoft/gardarike-online/model"
	rpc "abbysoft/gardarike-online/rpc/generated"
	"database/sql"
	"errors"
)

func (s *SimpleLogic) PlaceBuilding(session *PlayerSession, request *rpc.PlaceBuildingRequest) (*rpc.PlaceBuildingResponse, model.Error) {
	s.log.WithField("buildingID", request.GetBuildingID()).
		WithField("sessionID", request.GetSessionID()).
		WithField("location", *request.GetLocation()).
		Infof("PlaceBuilding request")

	building, found := s.buildings[int(request.BuildingID)]
	if !found {
		return nil, model.ErrBuildingNotFound
	}

	var location [3]float32
	location[0] = request.Location.X
	location[1] = request.Location.Y
	location[2] = request.Location.Z

	_, err := s.db.GetBuildingLocation(location)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		s.log.WithError(err).Error("Failed to get building on location")
		return nil, model.ErrInternalServerError
	} else if err == nil {
		return nil, model.ErrBuildingSpotIsBusy
	}

	if err := s.db.AddBuildingLocation(model.BuildingLocation{
		BuildingID: int(request.BuildingID),
		OwnerID:    session.SelectedCharacter.ID,
		Location:   location,
	}, true); err != nil {
		s.log.WithError(err).Error("Failed to add building location")
		return nil, model.ErrInternalServerError
	}

	s.GameMap.Buildings = append(s.GameMap.Buildings, &rpc.Building{
		Id:       request.BuildingID,
		OwnerID:  int64(session.SelectedCharacter.ID),
		Location: request.Location,
	})

	session.SelectedCharacter.MaxPopulation += building.PopulationBonus

	if err := s.db.UpdateCharacter(*session.SelectedCharacter, true); err != nil {
		s.log.WithError(err).Error("Failed to update character")
		return nil, model.ErrInternalServerError
	}

	return &rpc.PlaceBuildingResponse{}, nil
}
