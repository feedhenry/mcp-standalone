package metrics

import (
	"sync"
	"time"

	"runtime"

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

//might be able to make internal
type metric struct {
	Type   string
	XValue string
	YValue int64
}

type noServiceProvisionedErr struct {
	Message string
}

func (npe *noServiceProvisionedErr) Error() string {
	return npe.Message
}

func isNoServiceProvisionedErr(e error) bool {
	_, ok := e.(*noServiceProvisionedErr)
	return ok
}

// Gatherer is something that knows how to Gather metrics
type Gatherer func() ([]*metric, error)

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
	data map[string][]*mobile.GatheredMetric
	*sync.RWMutex
}

var internalMetrics = &metricsMap{
	RWMutex: &sync.RWMutex{},
	data:    map[string][]*mobile.GatheredMetric{},
}

func (mm *metricsMap) add(name string, m *metric) {
	mm.Lock()
	defer mm.Unlock()
	gathered := mm.data[name]
	if len(gathered) == 0 {
		gathered = []*mobile.GatheredMetric{{
			Type: m.Type,
			X:    []string{},
			Y:    map[string][]int64{},
		}}
		mm.data[name] = gathered
	}
	typeFound := false
	for i := range gathered {
		gm := gathered[i]
		if gm.Type == m.Type {
			typeFound = true
			gm.X = append(gm.X, m.XValue)
			if len(gm.X) > 30 {
				gm.X = gm.X[len(gm.X)-30:]
			}
			gm.Y[gm.Type] = append(gm.Y[gm.Type], m.YValue)
			if len(gm.Y[gm.Type]) > 30 {
				gm.Y[gm.Type] = gm.Y[gm.Type][len(gm.Y[gm.Type])-30:]
			}
			gathered[i] = gm
		}

	}
	if !typeFound {
		gathered = append(gathered, &mobile.GatheredMetric{
			Type: m.Type,
			X:    []string{m.XValue},
			Y: map[string][]int64{
				m.Type: {m.YValue},
			},
		})
	}
	mm.data[name] = gathered
}

func (mm *metricsMap) read(name string) []*mobile.GatheredMetric {
	mm.RLock()
	defer mm.RUnlock()
	return mm.data[name]
}

// Add allows new Gatherers to be added
func (gs *GathererScheduler) Add(serviceName string, metricGatherer Gatherer) {
	gs.jobs[serviceName] = metricGatherer
}

// Run will start the jobs on a schedule
func (gs *GathererScheduler) Run() {
	defer func() {
		gs.logger.Info("stopping metrics gatherers")
		gs.ticker.Stop()
		gs.logger.Info("ticker stopped")
	}()
	for {
		select {
		case <-gs.ticker.C:
			gs.execute()
		case <-gs.cancel:
			return
		}
	}
}

func (gs *GathererScheduler) execute() {

	//wait for the previous group to be done. If all completed will continue on
	gs.logger.Debug("executing gatherers after previous set done")
	gs.waitGroup.Wait()
	gs.logger.Debug("executing gatherers previous complete")
	for s, g := range gs.jobs {
		go func(service string, gather Gatherer) {
			defer func() {
				if err := recover(); err != nil {
					stack := make([]byte, 1024*8)
					stack = stack[:runtime.Stack(stack, false)]
					f := "PANIC: %s\n%s"
					gs.logger.Errorf(f, err, stack)
				}
			}()
			gs.waitGroup.Add(1)
			defer gs.waitGroup.Done()
			ms, err := gather()
			if err != nil && !isNoServiceProvisionedErr(err) {
				gs.logger.Error("unexpected error: failed to gather metrics for service ", service, err)
			}
			for _, m := range ms {
				gs.metrics.add(service, m)
			}
		}(s, g)
	}
}

type MetricsService struct{}

// Get will return the gathered metrics for a service
func (ms *MetricsService) GetAll(serviceName string) []*mobile.GatheredMetric {
	internalMetrics.RLock()
	defer internalMetrics.RUnlock()
	return internalMetrics.data[serviceName]
}

func (ms *MetricsService) GetOne(serviceName, metric string) *mobile.GatheredMetric {
	internalMetrics.RLock()
	defer internalMetrics.RUnlock()
	all := internalMetrics.data[serviceName]
	if len(all) == 0 {
		return nil
	}
	for _, m := range all {
		if m.Type == metric {
			return m
		}
	}
	return nil
}
