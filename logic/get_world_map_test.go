package logic

import (
	"abbysoft/gardarike-online/model"
	rpc "abbysoft/gardarike-online/rpc/generated"
	"database/sql"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestSimpleLogic_GetWorldMap_GeneratedWhenNoMapIsFound(t *testing.T) {
	logic, db, session, generator := NewLogicMockWithTerrainGenerator()
	request := &rpc.GetWorldMapRequest{
		SessionID: "sessionID",
		Location: &rpc.IntVector2D{
			X: 0,
			Y: 0,
		},
	}

	logic.config.AlwaysRegenerateMap = false
	logic.config.ChunkSize = 10

	db.On("GetMapChunk", int64(0), int64(0)).Return(model.WorldMapChunk{}, sql.ErrNoRows)
	db.On("SaveMapChunkOrUpdate", mock.Anything, mock.Anything).Return(nil)

	generator.On("GenerateTerrain", mock.Anything, mock.Anything, mock.Anything, mock.Anything).
		Return([]float32{10.0, 10.0})

	response, err := logic.GetWorldMap(session, request)
	require.NoError(t, err)
	require.NotNil(t, response)
	require.NotNil(t, response.Map)
	assert.NotEmpty(t, response.Map.Data)
	assert.Equal(t, int32(0), response.Map.X)
	assert.Equal(t, int32(0), response.Map.Y)
}
