package logic

import (
	"abbysoft/gardarike-online/model"
	"time"
)

const (
	gameLoopTps = 1.0
)

func (s *SimpleLogic) updateSessions() {
	sessionsCount := len(s.sessions)
	finishChan := make(chan bool, sessionsCount)

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
		for _ = range time.Tick(time.Second) {
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

func (s *SimpleLogic) updateSession(session *PlayerSession) {
	if time.Now().Sub(session.LastRequestTime) > s.config.AFKTimeout {
		s.log.WithField("sessionID", session.SessionID).
			Info("Session AFK timeout, delete session")
		delete(s.sessions, session.SessionID)
		return
	}

	character := session.SelectedCharacter

	populationGrownEvent := CheckRandomEventHappened(PopulationGrownEventChance)
	if populationGrownEvent && character.MaxPopulation != character.CurrentPopulation {
		s.characterPopulationGrownEvent(session)
	}

	resourcesGrownEvent := CheckRandomEventHappened(PlayerResourcesGrownEventChance)
	if resourcesGrownEvent {
		character.Resources.Add(model.Resources{
			Wood:    2,
			Food:    4,
			Stone:   1,
			Leather: 3,
		})

		s.log.WithField("sessionID", session.SessionID).
			WithField("character", session.SelectedCharacter.Name).
			WithField("resources", character.Resources).
			Debugf("Character resources have grown")

		if err := session.Tx.AddResourcesOrUpdate(character.Resources); err != nil {
			s.log.WithError(err).Error("Failed to update resources")
		}
	}
}
