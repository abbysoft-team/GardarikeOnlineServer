package logic

import (
	"abbysoft/gardarike-online/model"
	"abbysoft/gardarike-online/model/consts"
	rpc "abbysoft/gardarike-online/rpc/generated"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestSimpleLogic_PlaceTown_FirstTown(t *testing.T) {
	logic, db, session := NewLogicMock()
	request := &rpc.PlaceTownRequest{
		SessionID: "sessionID",
		Name:      "TestTown",
	}

	session.SelectedCharacter = &model.Character{
		ID:        1,
		AccountID: 1,
		Name:      "test",
	}

	db.On("AddTown", mock.MatchedBy(func(town model.Town) bool {
		return town.OwnerName == session.SelectedCharacter.Name &&
			town.Name == request.Name
	}), mock.Anything).Return(nil)

	db.On("UpdateCharacter", mock.MatchedBy(func(character model.Character) bool {
		return character.MaxPopulation == consts.TownPopulationBonus
	}), mock.Anything).Return(nil)

	resp, err := logic.PlaceTown(session, request)
	require.NoError(t, err)
	require.NotEmpty(t, resp)
	require.NotNil(t, resp.Location)
}

func TestSimpleLogic_PlaceTown_PlacingSecondTown(t *testing.T) {
	logic, db, session := NewLogicMock()
	request := &rpc.PlaceTownRequest{
		SessionID: "sessionID",
		Name:      "TestTown",
	}

	session.SelectedCharacter = &model.Character{
		ID:        1,
		AccountID: 1,
		Name:      "test",
		Towns: []model.Town{
			{
				X:          10,
				Y:          15,
				OwnerName:  "test",
				Population: 10,
				Name:       "TestTown",
			},
		},
		Resources: model.Resources{
			CharacterID: 1,
			Wood:        1000,
			Food:        1000,
			Stone:       1000,
			Leather:     1000,
		},
		MaxPopulation:     2000,
		CurrentPopulation: 1500,
	}

	db.On("AddTown", mock.MatchedBy(func(town model.Town) bool {
		return town.OwnerName == session.SelectedCharacter.Name &&
			town.Name == request.Name
	}), mock.Anything).Return(nil)

	db.On("UpdateCharacter", mock.MatchedBy(func(character model.Character) bool {
		return character.MaxPopulation == 2000+consts.TownPopulationBonus
	}), mock.Anything).Return(nil)

	resourcesAfterPlacing := model.Resources{
		CharacterID: 1,
		Wood:        1000,
		Food:        1000,
		Stone:       1000,
		Leather:     1000,
	}
	resourcesAfterPlacing.Subtract(model.ResourcesPlaceTown)

	db.On("AddResourcesOrUpdate", resourcesAfterPlacing, mock.Anything).Return(nil)

	resp, err := logic.PlaceTown(session, request)
	require.NoError(t, err)
	require.NotEmpty(t, resp)
	require.NotNil(t, resp.Location)

	// Trying to place another town. Should return an error because all resources were spent
	resp, err = logic.PlaceTown(session, request)
	require.EqualError(t, err, model.ErrNotEnoughResources.Error())
	require.Nil(t, resp)

	// Placing bellow the waterlevel
	session.SelectedCharacter.Resources.Add(model.ResourcesPlaceTown)
	logic.config.WaterLevel = 0.1
	logic.config.ChunkSize = 2

	chunk, convertErr := model.NewWorldMapChunkFromRPC(rpc.WorldMapChunk{
		Data: []float32{0.05, 0.04, 0.08, 0.09},
	})

	require.NoError(t, convertErr)

	db.On("GetMapChunk", int64(0), int64(0)).Return(chunk, nil)

	request.Location = &rpc.Vector2D{
		X: 1,
		Y: 1,
	}

	resp, err = logic.PlaceTown(session, request)
	require.EqualError(t, err, model.ErrBadRequest.Error())
	require.Nil(t, resp)
}
