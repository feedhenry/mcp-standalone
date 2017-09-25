package metrics

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"regexp"
	"strconv"
	"time"

	"bytes"

	"github.com/Sirupsen/logrus"
	"github.com/feedhenry/mcp-standalone/pkg/mobile"
	"github.com/pkg/errors"
)

// TODO: Add multi-worker support?
type usage struct {
	Master struct {
		Current         string    `json:"current"`
		Max             string    `json:"max"`
		Min             string    `json:"min"`
		Average         string    `json:"average"`
		NumberOfRecords int64     `json:"numberOfRecords"`
		From            time.Time `json:"from"`
		End             time.Time `json:"end"`
	} `json:"master"`
}

type workerProcessTime struct {
	Current         string    `json:"current"`
	Max             string    `json:"max"`
	Min             string    `json:"min"`
	Average         string    `json:"average"`
	NumberOfRecords int64     `json:"numberOfRecords"`
	From            time.Time `json:"from"`
	End             time.Time `json:"end"`
}

type workerQueueSize struct {
	Current         int64     `json:"current"`
	Max             int64     `json:"max"`
	Min             int64     `json:"min"`
	Average         float64   `json:"average"`
	NumberOfRecords int64     `json:"numberOfRecords"`
	From            time.Time `json:"from"`
	End             time.Time `json:"end"`
}

type timingStat struct {
	Current         string    `json:"current"`
	Max             string    `json:"max"`
	Min             string    `json:"min"`
	Average         string    `json:"average"`
	NumberOfRecords int64     `json:"numberOfRecords"`
	From            time.Time `json:"from"`
	End             time.Time `json:"end"`
}

type statsResponse struct {
	Metrics struct {
		CPUUsage       usage `json:"CPU usage"`
		RSSMemoryUsage usage `json:"RSS Memory Usage"`
		JobProcessTime struct {
			AckWorker     workerProcessTime `json:"ack_worker"`
			SyncWorker    workerProcessTime `json:"sync_worker"`
			PendingWorker workerProcessTime `json:"pending_worker"`
		} `json:"Job Process Time"`
		JobQueueSize struct {
			PendingWorker workerQueueSize `json:"pending_worker"`
			AckWorker     workerQueueSize `json:"ack_worker"`
			SyncWorker    workerQueueSize `json:"sync_worker"`
		} `json:"Job Queue Size"`
		APIProcessTime struct {
			Sync        timingStat `json:"sync"`
			SyncRecords timingStat `json:"syncRecords"`
		} `json:"API Process Time"`
		MongodbOperationTime struct {
			DoFindAndDeleteUpdate                  timingStat `json:"doFindAndDeleteUpdate"`
			DoUpdateManyDatasetClients             timingStat `json:"doUpdateManyDatasetClients"`
			DoListDatasetClients                   timingStat `json:"doListDatasetClients"`
			DoUpdateDatasetClient                  timingStat `json:"doUpdateDatasetClient"`
			DoListUpdates                          timingStat `json:"doListUpdates"`
			DoReadDatasetClient                    timingStat `json:"doReadDatasetClient"`
			DoUpdateDatasetClientWithRecords       timingStat `json:"doUpdateDatasetClientWithRecords"`
			DoReadDatasetClientWithRecordsUseCache timingStat `json:"doReadDatasetClientWithRecordsUseCache"`
			DoSaveUpdate                           timingStat `json:"doSaveUpdate"`
		} `json:"Mongodb Operation Time"`
	} `json:"metrics"`
}

type FhSyncServer struct {
	requestBuilder     mobile.HTTPRequesterBuilder
	serviceRepoBuilder mobile.ServiceRepoBuilder
	ServiceName        string
	logger             *logrus.Logger
}

func NewFhSyncServer(rbuilder mobile.HTTPRequesterBuilder, serviceRepoBuilder mobile.ServiceRepoBuilder, l *logrus.Logger) *FhSyncServer {
	return &FhSyncServer{requestBuilder: rbuilder, serviceRepoBuilder: serviceRepoBuilder, ServiceName: "fh-sync-server", logger: l}
}

/*
	"1.2things" => 1
	"1.2"       => 1
	"1ms"       => 1
	"1"         => 1
*/
var stringToInt64Regex = regexp.MustCompile(`(\d+)\.?\d*\w*`)

// Gather will retrieve varous metrics from fh-sync-server
func (ss *FhSyncServer) Gather() ([]*metric, error) {
	svc, err := ss.serviceRepoBuilder.UseDefaultSAToken().Build()
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

	var parseErrors = []error{}
	var stringToInt64 = func(val string) int64 {
		numOnly := stringToInt64Regex.FindStringSubmatch(val)
		if len(numOnly) < 2 {
			err := fmt.Errorf("no numerial value found in string %s", val)
			parseErrors = append(parseErrors, err)
			return -1
		}
		pInt, err := strconv.ParseInt(numOnly[1], 10, 64)
		if err != nil {
			parseErrors = append(parseErrors, err)
			return -1
		}
		return pInt
	}

	var ssMetrics = []*metric{}
	if nil != stats {
		// API Process Times
		if syncAPIProcessTime := stringToInt64(stats.Metrics.APIProcessTime.Sync.Current); syncAPIProcessTime != -1 {
			ssMetrics = append(ssMetrics, &metric{Type: "api_process_time_sync_ms", XValue: now.Format("2006-01-02 15:04:05"), YValue: syncAPIProcessTime})
		}
		if syncRecordsAPIProcessTime := stringToInt64(stats.Metrics.APIProcessTime.SyncRecords.Current); syncRecordsAPIProcessTime != -1 {
			ssMetrics = append(ssMetrics, &metric{Type: "api_process_time_syncRecords_ms", XValue: now.Format("2006-01-02 15:04:05"), YValue: syncRecordsAPIProcessTime})
		}

		// Mongodb Operation Times
		if doUpdateManyDatasetClients := stringToInt64(stats.Metrics.MongodbOperationTime.DoUpdateManyDatasetClients.Current); doUpdateManyDatasetClients != -1 {
			ssMetrics = append(ssMetrics, &metric{Type: "mongodb_operation_time_doUpdateManyDatasetClients_ms", XValue: now.Format("2006-01-02 15:04:05"), YValue: doUpdateManyDatasetClients})
		}
		if doListDatasetClients := stringToInt64(stats.Metrics.MongodbOperationTime.DoListDatasetClients.Current); doListDatasetClients != -1 {
			ssMetrics = append(ssMetrics, &metric{Type: "mongodb_operation_time_doListDatasetClients_ms", XValue: now.Format("2006-01-02 15:04:05"), YValue: doListDatasetClients})
		}
		if doSaveUpdate := stringToInt64(stats.Metrics.MongodbOperationTime.DoSaveUpdate.Current); doSaveUpdate != -1 {
			ssMetrics = append(ssMetrics, &metric{Type: "mongodb_operation_time_doSaveUpdate_ms", XValue: now.Format("2006-01-02 15:04:05"), YValue: doSaveUpdate})
		}
		if doUpdateDatasetClient := stringToInt64(stats.Metrics.MongodbOperationTime.DoUpdateDatasetClient.Current); doUpdateDatasetClient != -1 {
			ssMetrics = append(ssMetrics, &metric{Type: "mongodb_operation_time_doUpdateDatasetClient_ms", XValue: now.Format("2006-01-02 15:04:05"), YValue: doUpdateDatasetClient})
		}
		if doListUpdates := stringToInt64(stats.Metrics.MongodbOperationTime.DoListUpdates.Current); doListUpdates != -1 {
			ssMetrics = append(ssMetrics, &metric{Type: "mongodb_operation_time_doListUpdates_ms", XValue: now.Format("2006-01-02 15:04:05"), YValue: doListUpdates})
		}
		if doReadDatasetClient := stringToInt64(stats.Metrics.MongodbOperationTime.DoReadDatasetClient.Current); doReadDatasetClient != -1 {
			ssMetrics = append(ssMetrics, &metric{Type: "mongodb_operation_time_doReadDatasetClient_ms", XValue: now.Format("2006-01-02 15:04:05"), YValue: doReadDatasetClient})
		}
		if doUpdateDatasetClientWithRecords := stringToInt64(stats.Metrics.MongodbOperationTime.DoUpdateDatasetClientWithRecords.Current); doUpdateDatasetClientWithRecords != -1 {
			ssMetrics = append(ssMetrics, &metric{Type: "mongodb_operation_time_doUpdateDatasetClientWithRecords_ms", XValue: now.Format("2006-01-02 15:04:05"), YValue: doUpdateDatasetClientWithRecords})
		}
		if doReadDatasetClientWithRecordsUseCache := stringToInt64(stats.Metrics.MongodbOperationTime.DoReadDatasetClientWithRecordsUseCache.Current); doReadDatasetClientWithRecordsUseCache != -1 {
			ssMetrics = append(ssMetrics, &metric{Type: "mongodb_operation_time_doReadDatasetClientWithRecordsUseCache_ms", XValue: now.Format("2006-01-02 15:04:05"), YValue: doReadDatasetClientWithRecordsUseCache})
		}
		if doFindAndDeleteUpdate := stringToInt64(stats.Metrics.MongodbOperationTime.DoFindAndDeleteUpdate.Current); doFindAndDeleteUpdate != -1 {
			ssMetrics = append(ssMetrics, &metric{Type: "mongodb_operation_time_doFindAndDeleteUpdate_ms", XValue: now.Format("2006-01-02 15:04:05"), YValue: doFindAndDeleteUpdate})
		}

		// sync worker stats
		ssMetrics = append(ssMetrics, &metric{Type: "sync_worker_queue_count", XValue: now.Format("2006-01-02 15:04:05"), YValue: stats.Metrics.JobQueueSize.SyncWorker.Current})
		if syncProcessTime := stringToInt64(stats.Metrics.JobProcessTime.SyncWorker.Current); syncProcessTime != -1 {
			ssMetrics = append(ssMetrics, &metric{Type: "sync_worker_process_time_ms", XValue: now.Format("2006-01-02 15:04:05"), YValue: syncProcessTime})
		}
		if syncAvgProcessTime := stringToInt64(stats.Metrics.JobProcessTime.SyncWorker.Average); syncAvgProcessTime != -1 {
			ssMetrics = append(ssMetrics, &metric{Type: "sync_worker_process_time_ms_avg", XValue: now.Format("2006-01-02 15:04:05"), YValue: syncAvgProcessTime})
		}
		ssMetrics = append(ssMetrics, &metric{Type: "sync_worker_queue_count_total", XValue: now.Format("2006-01-02 15:04:05"), YValue: stats.Metrics.JobProcessTime.SyncWorker.NumberOfRecords})

		// pending worker stats
		ssMetrics = append(ssMetrics, &metric{Type: "pending_worker_queue_count", XValue: now.Format("2006-01-02 15:04:05"), YValue: stats.Metrics.JobQueueSize.PendingWorker.Current})
		if pendingProcessTime := stringToInt64(stats.Metrics.JobProcessTime.PendingWorker.Current); pendingProcessTime != -1 {
			ssMetrics = append(ssMetrics, &metric{Type: "pending_worker_process_time_ms", XValue: now.Format("2006-01-02 15:04:05"), YValue: pendingProcessTime})
		}
		if pendingAvgProcessTime := stringToInt64(stats.Metrics.JobProcessTime.PendingWorker.Average); pendingAvgProcessTime != -1 {
			ssMetrics = append(ssMetrics, &metric{Type: "pending_worker_process_time_ms_avg", XValue: now.Format("2006-01-02 15:04:05"), YValue: pendingAvgProcessTime})
		}
		ssMetrics = append(ssMetrics, &metric{Type: "pending_worker_queue_count_total", XValue: now.Format("2006-01-02 15:04:05"), YValue: stats.Metrics.JobProcessTime.PendingWorker.NumberOfRecords})

		// ack worker stats
		ssMetrics = append(ssMetrics, &metric{Type: "ack_worker_queue_count", XValue: now.Format("2006-01-02 15:04:05"), YValue: stats.Metrics.JobQueueSize.AckWorker.Current})
		if ackProcessTime := stringToInt64(stats.Metrics.JobProcessTime.AckWorker.Current); ackProcessTime != -1 {
			ssMetrics = append(ssMetrics, &metric{Type: "ack_worker_process_time_ms", XValue: now.Format("2006-01-02 15:04:05"), YValue: ackProcessTime})
		}
		if ackAvgProcessTime := stringToInt64(stats.Metrics.JobProcessTime.AckWorker.Average); ackAvgProcessTime != -1 {
			ssMetrics = append(ssMetrics, &metric{Type: "ack_worker_process_time_ms_avg", XValue: now.Format("2006-01-02 15:04:05"), YValue: ackAvgProcessTime})
		}
		ssMetrics = append(ssMetrics, &metric{Type: "ack_worker_queue_count_total", XValue: now.Format("2006-01-02 15:04:05"), YValue: stats.Metrics.JobProcessTime.AckWorker.NumberOfRecords})
	}

	if len(parseErrors) > 0 {
		logrus.Warn("Got the following errors when parsing sync metrics")
		for _, err := range parseErrors {
			logrus.Warn(err)
		}
	}

	return ssMetrics, nil
}

func (ss *FhSyncServer) getStats(host string) (*statsResponse, error) {
	statsURL := fmt.Sprintf("%s/sys/info/stats", host)
	requester := ss.requestBuilder.Insecure(true).Build()
	res, err := requester.Get(statsURL)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	data, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, errors.Wrap(err, "error reading stats response body")
	}
	ss.logger.Debugf("raw stats response %v", string(data))

	var stats statsResponse
	decoder := json.NewDecoder(bytes.NewBuffer(data))
	if err := decoder.Decode(&stats); err != nil {
		return nil, errors.Wrap(err, "failed to decode stats response")
	}
	ss.logger.Debugf("decoded stats response %v", stats)

	return &stats, nil
}
