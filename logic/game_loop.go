package logic

import (
	"abbysoft/gardarike-online/model"
	rpc "abbysoft/gardarike-online/rpc/generated"
	"time"
)

const (
	gameLoopTps = 1.0
)

func (s *SimpleLogic) updateSessions() {
	sessionsCount := len(s.sessions)
	finishChan := make(chan bool, sessionsCount)

	var buildings map[int64]model.CharacterBuildings

	tx, err := s.db.BeginTransaction(true, true)
	if err != nil {
		s.log.WithError(err).Error("Failed to begin transaction")
	} else if b, err := tx.GetAllBuildings(); err != nil {
		s.log.WithError(err).Error("Failed to get characters buildings")
	} else {
		buildings = b
	}

	if buildings == nil {
		buildings = make(map[int64]model.CharacterBuildings)
	}

	for _, session := range s.sessions {
		session := session

		go func() {
			session.Mutex.Lock()
			defer session.Mutex.Unlock()

			if session.SelectedCharacter == nil {
				finishChan <- true
				return
			}

			tx, err := s.db.BeginTransaction(false, true)
			if err != nil {
				s.log.WithError(err).Error("Failed to begin transaction")
				finishChan <- true
				return
			}

			session.Tx = tx

			s.updateSession(session)

			if buildings[session.SelectedCharacter.ID] != nil {
				s.updateSessionBuildings(session, buildings[session.SelectedCharacter.ID])
			}

			if !tx.IsCompleted() {
				if err := tx.EndTransaction(); err != nil {
					s.log.WithError(err).Error("Failed to commit transaction")
				}
			}

			finishChan <- true
		}()
	}

	for i := 0; i < sessionsCount; i++ {
		<-finishChan
	}
}

// startGameLoop - runs endless game loop
func (s *SimpleLogic) startGameLoop() {
	go func() {
		for _ = range time.Tick(5 * time.Second) {
			s.updateSessions()
		}
	}()

	go func() {
		for _ = range time.Tick(time.Minute) {
			s.resourceManager.Update()
		}
	}()
}

func (s *SimpleLogic) characterPopulationGrownEvent(session *PlayerSession) {
	session.SelectedCharacter.CurrentPopulation++
	session.WorkDistribution.IdleCount++

	s.log.WithField("sessionID", session.SessionID).
		WithField("character", session.SelectedCharacter.Name).
		Debugf("Player's population grows")

	if err := session.Tx.UpdateCharacter(*session.SelectedCharacter); err != nil {
		s.log.WithError(err).Error("Failed to update character")
	}
}

func (s *SimpleLogic) updateSessionBuildings(session *PlayerSession, buildings model.CharacterBuildings) {
	character := session.SelectedCharacter

	houseNumber := buildings[rpc.BuildingType_HOUSE]
	quarryNumber := buildings[rpc.BuildingType_QUARRY]

	character.Resources.Add(model.Resources{Food: houseNumber, Stone: quarryNumber})
	if err := session.Tx.AddResourcesOrUpdate(character.ID, character.Resources); err != nil {
		s.log.WithError(err).Error("Failed to update character resources")
	}
}

func (s *SimpleLogic) updateSession(session *PlayerSession) {
	if time.Now().Sub(session.LastRequestTime) > s.config.AFKTimeout {
		s.log.WithField("sessionID", session.SessionID).
			WithField("timeout", s.config.AFKTimeout).
			Info("Session AFK timeout, delete session")
		delete(s.sessions, session.SessionID)
		return
	}

	character := session.SelectedCharacter

	populationGrownEvent := CheckRandomEventHappened(PopulationGrownEventChance)
	if populationGrownEvent && character.MaxPopulation != character.CurrentPopulation {
		s.characterPopulationGrownEvent(session)
	}

	if !character.Resources.IsEnough(model.ResourcesLimit) {
		character.Resources.Add(model.Resources{
			Wood:    2,
			Food:    3,
			Stone:   1,
			Leather: 3,
		})

		if err := session.Tx.AddResourcesOrUpdate(character.ID, character.Resources); err != nil {
			s.log.WithError(err).Error("Failed to update resources")
		}
	}
}
