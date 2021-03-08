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
	tx, err := r.logic.db.BeginTransaction(true, true)
	if err != nil {
		r.logger.WithError(err).Error("Failed to begin transaction")
		return
	}

	if err := tx.IncrementMapResources(resourceIncrementValue); err != nil {
		r.logger.WithError(err).Error("Failed to increment map resources")
	}
}
