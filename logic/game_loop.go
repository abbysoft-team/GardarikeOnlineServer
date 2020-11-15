package logic

import (
	"abbysoft/gardarike-online/model"
	"time"
)

const (
	gameLoopTps = 1.0
)

// gameLoop - runs endless game loop
func (s *SimpleLogic) gameLoop() {
	sleepDuration := time.Duration(1000.0/gameLoopTps) * time.Millisecond

	for {
		sessionsCount := len(s.sessions)
		finishChan := make(chan bool, sessionsCount)

		for _, session := range s.sessions {
			session := session

			go func() {
				if session.SelectedCharacter == nil {
					finishChan <- true
					return
				}

				s.updateSession(session)
				finishChan <- true
			}()
		}

		for i := 0; i < sessionsCount; i++ {
			<-finishChan
		}

		s.resourceManager.Update()

		time.Sleep(sleepDuration)
	}
}

func (s *SimpleLogic) characterPopulationGrownEvent(session *PlayerSession) {
	session.Mutex.Lock()
	defer session.Mutex.Unlock()

	if session.SelectedCharacter.MaxPopulation != session.SelectedCharacter.CurrentPopulation {
		session.SelectedCharacter.CurrentPopulation++

		s.log.WithField("sessionID", session.SessionID).
			WithField("character", session.SelectedCharacter.Name).
			Debugf("Player's population grows")

		if err := s.db.UpdateCharacter(*session.SelectedCharacter); err != nil {
			s.log.WithError(err).Error("Failed to update character")
		} else {
			s.EventsChan <- model.EventWrapper{
				Topic: session.SessionID,
				Event: model.NewCharacterUpdatedEvent(session.SelectedCharacter),
			}
		}
	}
}

func (s *SimpleLogic) updateSession(session *PlayerSession) {
	if time.Now().Sub(session.LastRequestTime) > s.config.AFKTimeout {
		s.log.WithField("sessionID", session.SessionID).
			Info("Session AFK timeout, delete session")
		delete(s.sessions, session.SessionID)
	}

	populationGrownEvent := CheckRandomEventHappened(PopulationGrownEventChance)
	if populationGrownEvent {
		s.characterPopulationGrownEvent(session)
	}
}
