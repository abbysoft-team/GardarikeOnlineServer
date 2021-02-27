package logic

import (
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
	//resourceIncremented := CheckRandomEventHappened(ResourceIncrementChance)
	//if resourceIncremented {
	//	r.logic.GameMapMutex.Lock()
	//	defer r.logic.GameMapMutex.Unlock()
	//
	//	resourceType := ResourceEvent(rand.Intn(3) + 1)
	//	switch resourceType {
	//	case TreeIncrementedEvent:
	//		r.logic.GameMap.Trees++
	//	case StoneIncrementedEvent:
	//		r.logic.GameMap.Stones++
	//	case AnimalIncrementedEvent:
	//		r.logic.GameMap.Animals++
	//	case PlantsIncrementedEvent:
	//		r.logic.GameMap.Plants++
	//	}
	//
	//	if err := r.logic.SaveGameMap(); err != nil {
	//		r.logger.WithError(err).Error("Failed to update resources: failed to save game map")
	//	}
	//
	//	r.logger.WithFields(log.Fields{
	//		"trees":   r.logic.GameMap.Trees,
	//		"stones":  r.logic.GameMap.Stones,
	//		"plants":  r.logic.GameMap.Plants,
	//		"animals": r.logic.GameMap.Animals,
	//	}).Info("Resources incremented on the map")
	//}
}
