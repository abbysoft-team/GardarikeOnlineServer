package logic

import (
	"abbysoft/gardarike-online/model"
	"math/rand"
	"time"
)

const (
	gameLoopTps                = 1.0
	populationGrownEventChance = 2.0
)

// checkRandomEventHappened - check if the random event of 'chance' percent freq has happened
// returns true if the event has happened and false otherwise
func checkRandomEventHappened(chance int) bool {
	rand.Seed(time.Now().UnixNano())
	n := rand.Int63n(10e10)
	chance *= 10e8

	if n >= int64(chance) {
		return false
	}

	return true
}

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
			s.eventsChan <- model.EventWrapper{
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

	populationGrownEvent := checkRandomEventHappened(populationGrownEventChance)
	if populationGrownEvent {
		s.characterPopulationGrownEvent(session)
	}
}
