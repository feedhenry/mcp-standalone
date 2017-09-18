package metrics

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"strconv"
	"strings"
	"time"

	"bytes"

	"github.com/Sirupsen/logrus"
	"github.com/feedhenry/mcp-standalone/pkg/mobile"
	"github.com/pkg/errors"
)

// TODO: Add multi-worker support?
type statsResponse struct {
	Metrics struct {
		CPUUsage struct {
			Master struct {
				Current         string    `json:"current"`
				Max             string    `json:"max"`
				Min             string    `json:"min"`
				Average         string    `json:"average"`
				NumberOfRecords int64     `json:"numberOfRecords"`
				From            time.Time `json:"from"`
				End             time.Time `json:"end"`
			} `json:"master"`
		} `json:"CPU usage"`
		RSSMemoryUsage struct {
			Master struct {
				Current         string    `json:"current"`
				Max             string    `json:"max"`
				Min             string    `json:"min"`
				Average         string    `json:"average"`
				NumberOfRecords int64     `json:"numberOfRecords"`
				From            time.Time `json:"from"`
				End             time.Time `json:"end"`
			} `json:"master"`
		} `json:"RSS Memory Usage"`
		JobProcessTime struct {
			SyncWorker struct {
				Current         string    `json:"current"`
				Max             string    `json:"max"`
				Min             string    `json:"min"`
				Average         string    `json:"average"`
				NumberOfRecords int64     `json:"numberOfRecords"`
				From            time.Time `json:"from"`
				End             time.Time `json:"end"`
			} `json:"sync_worker"`
			AckWorker struct {
				Current         string    `json:"current"`
				Max             string    `json:"max"`
				Min             string    `json:"min"`
				Average         string    `json:"average"`
				NumberOfRecords int64     `json:"numberOfRecords"`
				From            time.Time `json:"from"`
				End             time.Time `json:"end"`
			} `json:"ack_worker"`
			PendingWorker struct {
				Current         string    `json:"current"`
				Max             string    `json:"max"`
				Min             string    `json:"min"`
				Average         string    `json:"average"`
				NumberOfRecords int64     `json:"numberOfRecords"`
				From            time.Time `json:"from"`
				End             time.Time `json:"end"`
			} `json:"pending_worker"`
		} `json:"Job Process Time"`
		JobQueueSize struct {
			PendingWorker struct {
				Current         int64     `json:"current"`
				Max             int64     `json:"max"`
				Min             int64     `json:"min"`
				Average         float64   `json:"average"`
				NumberOfRecords int64     `json:"numberOfRecords"`
				From            time.Time `json:"from"`
				End             time.Time `json:"end"`
			} `json:"pending_worker"`
			AckWorker struct {
				Current         int64     `json:"current"`
				Max             int64     `json:"max"`
				Min             int64     `json:"min"`
				Average         float64   `json:"average"`
				NumberOfRecords int64     `json:"numberOfRecords"`
				From            time.Time `json:"from"`
				End             time.Time `json:"end"`
			} `json:"ack_worker"`
			SyncWorker struct {
				Current         int64     `json:"current"`
				Max             int64     `json:"max"`
				Min             int64     `json:"min"`
				Average         float64   `json:"average"`
				NumberOfRecords int64     `json:"numberOfRecords"`
				From            time.Time `json:"from"`
				End             time.Time `json:"end"`
			} `json:"sync_worker"`
		} `json:"Job Queue Size"`
		APIProcessTime struct {
			Sync struct {
				Current         string    `json:"current"`
				Max             string    `json:"max"`
				Min             string    `json:"min"`
				Average         string    `json:"average"`
				NumberOfRecords int64     `json:"numberOfRecords"`
				From            time.Time `json:"from"`
				End             time.Time `json:"end"`
			} `json:"sync"`
		} `json:"API Process Time"`
		MongodbOperationTime struct {
			DoUpdateManyDatasetClients struct {
				Current         string    `json:"current"`
				Max             string    `json:"max"`
				Min             string    `json:"min"`
				Average         string    `json:"average"`
				NumberOfRecords int64     `json:"numberOfRecords"`
				From            time.Time `json:"from"`
				End             time.Time `json:"end"`
			} `json:"doUpdateManyDatasetClients"`
			DoListDatasetClients struct {
				Current         string    `json:"current"`
				Max             string    `json:"max"`
				Min             string    `json:"min"`
				Average         string    `json:"average"`
				NumberOfRecords int64     `json:"numberOfRecords"`
				From            time.Time `json:"from"`
				End             time.Time `json:"end"`
			} `json:"doListDatasetClients"`
		} `json:"Mongodb Operation Time"`
	} `json:"metrics"`
}

type FhSyncServer struct {
	requestBuilder     mobile.HTTPRequesterBuilder
	tokenClientBuilder mobile.TokenScopedClientBuilder
	ServiceName        string
	logger             *logrus.Logger
}

func NewFhSyncServer(rbuilder mobile.HTTPRequesterBuilder, tokenCBuilder mobile.TokenScopedClientBuilder, l *logrus.Logger) *FhSyncServer {
	return &FhSyncServer{requestBuilder: rbuilder, tokenClientBuilder: tokenCBuilder, ServiceName: "fh-sync-server", logger: l}
}

// Gather will retrieve varous metrics from fh-sync-server
func (ss *FhSyncServer) Gather() ([]*metric, error) {
	svc, err := ss.tokenClientBuilder.UseDefaultSAToken().MobileServiceCruder("")
	if err != nil {
		return nil, errors.Wrap(err, "fh-sync-server gather failed to create svcruder")
	}
	ssServices, err := svc.List(func(attrs mobile.Attributer) bool {
		return attrs.GetName() == ss.ServiceName
	})
	if err != nil {
		return nil, errors.Wrap(err, "fh-sync-server gather failed to list existing services")
	}
	if len(ssServices) == 0 {
		return nil, errors.New(" no fh-sync-server service present")
	}
	ssService := ssServices[0] //TODO deal with more than one
	//TODO get protocol from secret
	host := ssService.Host
	now := time.Now()

	stats, err := ss.getStats(host)
	if err != nil {
		return nil, errors.Wrap(err, "fh-sync-server gather failed to list existing services")
	}

	var ssMetrics = []*metric{}
	if nil != stats {
		// sync worker stats
		ssMetrics = append(ssMetrics, &metric{Type: "sync_worker_queue_count", XValue: now.Format("2006-01-02 15:04:05"), YValue: stats.Metrics.JobQueueSize.SyncWorker.Current})
		syncProcessTimeStr := strings.Split(stats.Metrics.JobProcessTime.SyncWorker.Current, ".")[0]
		syncProcessTime, _ := strconv.ParseInt(syncProcessTimeStr, 10, 64)
		ssMetrics = append(ssMetrics, &metric{Type: "sync_worker_process_time_ms", XValue: now.Format("2006-01-02 15:04:05"), YValue: syncProcessTime})

		// pending worker stats
		ssMetrics = append(ssMetrics, &metric{Type: "pending_worker_queue_count", XValue: now.Format("2006-01-02 15:04:05"), YValue: stats.Metrics.JobQueueSize.PendingWorker.Current})
		pendingProcessTimeStr := strings.Split(stats.Metrics.JobProcessTime.PendingWorker.Current, ".")[0]
		pendingProcessTime, _ := strconv.ParseInt(pendingProcessTimeStr, 10, 64)
		ssMetrics = append(ssMetrics, &metric{Type: "pending_worker_process_time_ms", XValue: now.Format("2006-01-02 15:04:05"), YValue: pendingProcessTime})

		// ack worker stats
		ssMetrics = append(ssMetrics, &metric{Type: "ack_worker_queue_count", XValue: now.Format("2006-01-02 15:04:05"), YValue: stats.Metrics.JobQueueSize.AckWorker.Current})
		ackProcessTimeStr := strings.Split(stats.Metrics.JobProcessTime.AckWorker.Current, ".")[0]
		ackProcessTime, _ := strconv.ParseInt(ackProcessTimeStr, 10, 64)
		ssMetrics = append(ssMetrics, &metric{Type: "ack_worker_process_time_ms", XValue: now.Format("2006-01-02 15:04:05"), YValue: ackProcessTime})
	}

	return ssMetrics, nil
}

func (ss *FhSyncServer) getStats(host string) (*statsResponse, error) {
	statsURL := fmt.Sprintf("%s/sys/info/stats", host)
	requester := ss.requestBuilder.Insecure(true).Build()
	res, err := requester.Get(statsURL)
	if err != nil {
		return &statsResponse{}, err
	}
	defer res.Body.Close()

	data, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return &statsResponse{}, errors.Wrap(err, "error reading stats response body")
	}
	fmt.Printf("\n\ndata %s ", string(data))

	var stats statsResponse
	decoder := json.NewDecoder(bytes.NewBuffer(data))
	if err := decoder.Decode(&stats); err != nil {
		return &statsResponse{}, errors.Wrap(err, "failed to decode stats response")
	}

	fmt.Printf("\n\nstats %v ", stats)
	fmt.Printf("\n\nstats.JobQueueSize.SyncWorker.NumberOfRecords (%v) ", stats.Metrics.JobQueueSize.SyncWorker.NumberOfRecords)

	return &stats, nil
}
