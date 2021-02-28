package logic

import (
	"abbysoft/gardarike-online/model"
	log "github.com/sirupsen/logrus"
	"time"
)

var (
	resourceUpdateFreq = 1 * time.Minute
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

var resourceIncrementValue = model.ChunkResources{
	Trees:   5,
	Stones:  1,
	Animals: 8,
	Plants:  6,
}

func (r *ResourceManager) Update() {
	if err := r.logic.db.IncrementMapResources(resourceIncrementValue, true); err != nil {
		r.logger.WithError(err).Error("Failed to increment map resources")
	}
}
