package logic

import (
	"abbysoft/gardarike-online/model"
	log "github.com/sirupsen/logrus"
)

type ResourceManager struct {
	logic  *SimpleLogic
	logger *log.Entry
}

func NewResourceManager(l *SimpleLogic) ResourceManager {
	return ResourceManager{
		logic:  l,
		logger: log.WithField("module", "resource_manager"),
	}
}

func (r *ResourceManager) Update() {
	treeGrownEvent := CheckRandomEventHappened(TreeGrownEventChance)
	if treeGrownEvent {
		r.logic.GameMapMutex.Lock()
		defer r.logic.GameMapMutex.Unlock()

		r.logic.GameMap.TreesCount++
		if err := r.logic.SaveGameMap(); err != nil {
			r.logger.WithError(err).Error("Failed to update resources: failed to save game map")
		}

		r.logger.WithField("treesCount", r.logic.GameMap.TreesCount).Info("Trees incremented on the map")

		r.logic.EventsChan <- model.EventWrapper{
			Topic: model.GlobalTopic,
			Event: model.NewResourceUpdatedEvent(r.logic.GameMap.TreesCount),
		}
	}
}
