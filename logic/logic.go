package logic

import (
	db2 "abbysoft/gardarike-online/db"
	"abbysoft/gardarike-online/db/postgres"
	"abbysoft/gardarike-online/generation"
	"abbysoft/gardarike-online/model"
	"abbysoft/gardarike-online/model/consts"
	rpc "abbysoft/gardarike-online/rpc/generated"
	"fmt"
	"github.com/sirupsen/logrus"
	"time"
)

const (
	mapChunkSize = 100
)

type Logic interface {
	GetWorldMap(session *PlayerSession, request *rpc.GetWorldMapRequest) (*rpc.GetWorldMapResponse, model.Error)
	Login(request *rpc.LoginRequest) (*rpc.LoginResponse, model.Error)
	SelectCharacter(session *PlayerSession, request *rpc.SelectCharacterRequest) (*rpc.SelectCharacterResponse, model.Error)
	SendChatMessage(session *PlayerSession, request *rpc.SendChatMessageRequest) (*rpc.SendChatMessageResponse, model.Error)
	GetChatHistory(session *PlayerSession, request *rpc.GetChatHistoryRequest) (*rpc.GetChatHistoryResponse, model.Error)
	GetWorkDistribution(session *PlayerSession, request *rpc.GetWorkDistributionRequest) (*rpc.GetWorkDistributionResponse, model.Error)
	CreateAccount(request *rpc.CreateAccountRequest) (*rpc.CreateAccountResponse, model.Error)
	CreateCharacter(session *PlayerSession, request *rpc.CreateCharacterRequest) (*rpc.CreateCharacterResponse, model.Error)
	GetResources(session *PlayerSession, request *rpc.GetResourcesRequest) (*rpc.GetResourcesResponse, model.Error)
	PlaceTown(session *PlayerSession, request *rpc.PlaceTownRequest) (*rpc.PlaceTownResponse, model.Error)
	PlaceBuilding(session *PlayerSession, request *rpc.PlaceBuildingRequest) (*rpc.PlaceBuildingResponse, model.Error)
}

type SimpleLogic struct {
	db              db2.Database
	log             *logrus.Entry
	sessions        map[string]*PlayerSession
	EventsChan      chan model.EventWrapper
	config          Config
	resourceManager ResourceManager
	generator       generation.TerrainGenerator
}

type Config struct {
	AFKTimeout           time.Duration
	ChatMessageMaxLength int
	WaterLevel           float32
	ChunkSize            int
	AlwaysRegenerateMap  bool
}

func NewLogic(generator generation.TerrainGenerator, eventsChan chan model.EventWrapper, dbConfig postgres.Config, config Config) (*SimpleLogic, error) {
	database, err := postgres.NewDatabase(dbConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to init db: %w", err)
	}

	logic := &SimpleLogic{
		db:         database,
		log:        logrus.WithField("module", "logic"),
		sessions:   make(map[string]*PlayerSession),
		EventsChan: eventsChan,
		config:     config,
		generator:  generator,
	}

	logic.resourceManager = NewResourceManager(logic)

	logic.log.Info("Running game loop")
	logic.startGameLoop()

	return logic, nil
}

func (s *SimpleLogic) SelectCharacter(session *PlayerSession, request *rpc.SelectCharacterRequest) (*rpc.SelectCharacterResponse, model.Error) {
	s.log.WithField("characterID", request.GetCharacterID()).
		WithField("sessionID", request.GetSessionID()).
		Infof("SelectCharacter request")

	tx := session.Tx

	char, err := tx.GetCharacter(request.GetCharacterID())
	if err != nil {
		return nil, model.ErrCharacterNotFound
	}

	if char.AccountID != session.AccountID {
		return nil, model.ErrForbidden
	}

	towns, err := tx.GetTowns(char.Name)
	if err != nil {
		s.log.WithError(err).Error("Failed to get character's towns")
		return nil, model.ErrInternalServerError
	} else {
		char.Towns = towns
	}

	resources, err := tx.GetResources(char.ID)
	if err != nil {
		s.log.WithError(err).Error("Failed to get character's resources")
		return nil, model.ErrInternalServerError
	}

	char.Resources = resources

	session.SelectedCharacter = &char
	s.log.WithFields(logrus.Fields{
		"sessionID": request.GetSessionID(),
		"character": char,
	}).Info("User selected character")

	s.EventsChan <- model.NewSystemChatMessageEvent(consts.MessageCharacterAuthorized(char.Name))

	response := &rpc.SelectCharacterResponse{Resources: resources.ToRPC()}
	for _, town := range char.Towns {
		response.Towns = append(response.Towns, town.ToRPC())
	}

	return response, nil
}
