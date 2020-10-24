package game

import (
	"database/sql"
	"errors"
	"projectx-server/model"
	rpc "projectx-server/rpc/generated"
)

func (s *SimpleLogic) PlaceBuilding(request *rpc.PlaceBuildingRequest) (*rpc.PlaceBuildingResponse, model.Error) {
	s.log.WithField("buildingID", request.GetBuildingID()).
		WithField("sessionID", request.GetSessionID()).
		WithField("location", *request.GetLocation()).
		Infof("PlaceBuilding request")

	session, authorized := s.sessions[request.GetSessionID()]
	if !authorized {
		return nil, model.ErrNotAuthorized
	}

	building, found := s.buildings[int(request.BuildingID)]
	if !found {
		return nil, model.ErrBuildingNotFound
	}

	if uint64(building.Cost) > session.SelectedCharacter.Gold {
		return nil, model.ErrNoEnoughMoney
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
	}); err != nil {
		s.log.WithError(err).Error("Failed to add building location")
		return nil, model.ErrInternalServerError
	}

	s.eventsChan <- model.NewPlaceBuildingEvent(building.ID, session.SelectedCharacter.ID, request.Location)

	session.SelectedCharacter.Gold -= uint64(building.Cost)
	if err := s.db.UpdateCharacter(*session.SelectedCharacter); err != nil {
		s.log.WithError(err).Error("Failed to decrease character's gold")
		return nil, model.ErrInternalServerError
	}

	return &rpc.PlaceBuildingResponse{}, nil
}
