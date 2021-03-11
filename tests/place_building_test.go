package tests

import (
	rpc "abbysoft/gardarike-online/rpc/generated"
	"github.com/stretchr/testify/assert"
	"math/rand"
	"testing"
	"time"
)

func TestPlaceBuilding(t *testing.T) {
	TestSelectCharacter(t)

	rand.Seed(time.Now().UnixNano())

	var request rpc.Request
	request.Data = &rpc.Request_PlaceBuildingRequest{
		PlaceBuildingRequest: &rpc.PlaceBuildingRequest{
			SessionID:  sessionID,
			Location:   &rpc.Vector2D{X: rand.Float32() * 1000, Y: rand.Float32() * 1000},
			BuildingID: rpc.BuildingType(rand.Intn(int(rpc.BuildingType_QUARRY)) + 1),
			TownID:     4,
		},
	}

	resp, err := client.SendRequest(request)

	if !assert.NoError(t, err, "request error is not nil") {
		return
	}
	if !assert.NotNil(t, resp, "response is nil") {
		return
	}
	if !assert.NotNil(t, resp.GetPlaceBuildingResponse(), "response isn't a place town response") {
		return
	}
}
