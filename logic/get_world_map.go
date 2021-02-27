package logic

import (
	"abbysoft/gardarike-online/model"
	rpc "abbysoft/gardarike-online/rpc/generated"
	"database/sql"
	"errors"
	"fmt"
	log "github.com/sirupsen/logrus"
	"time"
)

func (s *SimpleLogic) saveChunk(chunk rpc.WorldMapChunk) error {
	modelChunk, err := model.NewWorldMapChunkFromRPC(chunk)
	if err != nil {
		return err
	}

	return s.db.SaveMapChunkOrUpdate(modelChunk, true)
}

func (s *SimpleLogic) generateAndSaveMapChunk(x, y int) (*rpc.WorldMapChunk, error) {
	s.log.WithFields(log.Fields{
		"x": x,
		"y": y,
	}).Info("Generating map chunk")

	terrain := s.generator.GenerateTerrain(
		s.config.ChunkSize,
		s.config.ChunkSize,
		float64(s.config.ChunkSize*x),
		float64(s.config.ChunkSize*y))

	chunk := rpc.WorldMapChunk{
		X:          int32(x),
		Y:          int32(y),
		Width:      int32(s.config.ChunkSize),
		Height:     int32(s.config.ChunkSize),
		Data:       terrain,
		Towns:      []*rpc.Town{},
		Trees:      0,
		Stones:     0,
		Animals:    0,
		Plants:     0,
		WaterLevel: s.config.WaterLevel,
	}

	if err := s.saveChunk(chunk); err != nil {
		return nil, fmt.Errorf("failed to save map chunk: %w", err)
	}

	return &chunk, nil
}

func (s *SimpleLogic) GetWorldMap(_ *PlayerSession, request *rpc.GetWorldMapRequest) (*rpc.GetWorldMapResponse, model.Error) {
	s.log.WithField("location", request.GetLocation()).
		WithField("sessionID", request.GetSessionID()).
		Infof("GetMap request")

	newChunk := func() (*rpc.GetWorldMapResponse, model.Error) {
		s.log.WithField("alwaysGenerate", s.config.AlwaysRegenerateMap).
			WithField("location", request.GetLocation()).Info("Generating chunk")

		if newChunk, err := s.generateAndSaveMapChunk(int(request.Location.X), int(request.Location.Y)); err != nil {
			s.log.WithError(err).Error("Failed to regenerate game map")
			return nil, model.ErrInternalServerError
		} else {
			return &rpc.GetWorldMapResponse{Map: newChunk}, nil
		}
	}

	if s.config.AlwaysRegenerateMap {
		s.generator.SetSeed(time.Now().UnixNano())
		return newChunk()
	}

	chunk, err := s.db.GetMapChunk(int64(request.Location.X), int64(request.Location.Y))
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		s.log.WithError(err).Error("Failed to get map chunk")
		return nil, model.ErrInternalServerError
	}

	if len(chunk.Data) == 0 {
		return newChunk()
	}

	rpcChunk, err := chunk.ToRPC()
	if err != nil {
		s.log.WithError(err).Error("Failed to convert map chunk to the rpc chunk")
		return nil, model.ErrInternalServerError
	}

	xStart := int(request.Location.X) * mapChunkSize
	xEnd := xStart + mapChunkSize
	yStart := int(request.Location.Y) * mapChunkSize
	yEnd := yStart + mapChunkSize

	towns, err := s.db.GetTownsForRect(xStart, xEnd, yStart, yEnd)
	if err != nil {
		s.log.WithError(err).Error("Failed to get chunk towns")
		return nil, model.ErrInternalServerError
	}

	for _, town := range towns {
		rpcChunk.Towns = append(rpcChunk.Towns, town.ToRPC())
	}

	return &rpc.GetWorldMapResponse{Map: rpcChunk}, nil
}
