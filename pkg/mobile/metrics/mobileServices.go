package metrics

import (
	"fmt"
	"sync"
	"time"

	"github.com/Sirupsen/logrus"
	"github.com/feedhenry/mcp-standalone/pkg/mobile"
)

// GathererScheduler schedules metrics Gathering Jobs
type GathererScheduler struct {
	ticker    *time.Ticker
	cancel    chan struct{}
	jobs      map[string]Gatherer
	logger    *logrus.Logger
	metrics   *metricsMap
	waitGroup *sync.WaitGroup
}

// NewGathererScheduler creates a default GathererScheduler
func NewGathererScheduler(ticker *time.Ticker, cancel chan struct{}, logger *logrus.Logger) *GathererScheduler {
	return &GathererScheduler{
		ticker:    ticker,
		cancel:    cancel,
		jobs:      map[string]Gatherer{},
		logger:    logger,
		metrics:   internalMetrics,
		waitGroup: &sync.WaitGroup{},
	}
}

type metricsMap struct {
	data map[string][]*mobile.Metric
	*sync.RWMutex
}

var internalMetrics = &metricsMap{
	RWMutex: &sync.RWMutex{},
	data:    map[string][]*mobile.Metric{},
}

func (mm *metricsMap) add(name string, m *mobile.Metric) {
	mm.Lock()
	defer mm.Unlock()
	mm.data[name] = append(mm.data[name], m)
}

func (mm *metricsMap) read(name string) []*mobile.Metric {
	mm.RLock()
	defer mm.RUnlock()
	return mm.data[name]
}

// Gatherer is something that knows how to Gather metrics
type Gatherer func() (*mobile.Metric, error)

// Add allows new Gatherers to be added
func (gs *GathererScheduler) Add(serviceName string, metricGatherer Gatherer) {}

// Run will start the jobs on a schedule
func (gs *GathererScheduler) Run() {
	for {
		select {
		case <-gs.ticker.C:
			gs.execute()
			//run gather jobs
		case <-gs.cancel:
			gs.logger.Info("stopping metrics gatherers")
			gs.ticker.Stop()
			gs.logger.Info("ticker stopped")
			return
		}
	}
}

func (gs *GathererScheduler) execute() {
	//wait for the previous group to be done. If all completed will continue on
	gs.logger.Debug("executing gatherers after previos set done")
	gs.waitGroup.Wait()
	gs.logger.Debug("executing gatherers previos complete")
	for s, g := range gs.jobs {
		go func(service string, gather Gatherer) {
			gs.waitGroup.Add(1)
			defer gs.waitGroup.Done()
			m, err := gather()
			if err != nil {
				gs.logger.Error("failed to gather metrics for service ", service, err)
				fmt.Println("gathered metrics for service ", service, m)
				gs.metrics.add(service, m)
			}
		}(s, g)
	}
}

type MetricsService struct{}

// Get will return the gathered metrics for a service
func (ms *MetricsService) Get(serviceName, metric string) *mobile.Metric {
	internalMetrics.RLock()
	defer internalMetrics.RUnlock()
	metrics := internalMetrics.data[serviceName]
	for _, m := range metrics {
		if m.Type == metric {
			return m
		}
	}
	return nil
}
