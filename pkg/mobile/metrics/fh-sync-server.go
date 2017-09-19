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
			AckWorker struct {
				Current         string    `json:"current"`
				Max             string    `json:"max"`
				Min             string    `json:"min"`
				Average         string    `json:"average"`
				NumberOfRecords int64     `json:"numberOfRecords"`
				From            time.Time `json:"from"`
				End             time.Time `json:"end"`
			} `json:"ack_worker"`
			SyncWorker struct {
				Current         string    `json:"current"`
				Max             string    `json:"max"`
				Min             string    `json:"min"`
				Average         string    `json:"average"`
				NumberOfRecords int64     `json:"numberOfRecords"`
				From            time.Time `json:"from"`
				End             time.Time `json:"end"`
			} `json:"sync_worker"`
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
			SyncRecords struct {
				Current         string    `json:"current"`
				Max             string    `json:"max"`
				Min             string    `json:"min"`
				Average         string    `json:"average"`
				NumberOfRecords int64     `json:"numberOfRecords"`
				From            time.Time `json:"from"`
				End             time.Time `json:"end"`
			} `json:"syncRecords"`
		} `json:"API Process Time"`
		MongodbOperationTime struct {
			DoFindAndDeleteUpdate struct {
				Current         string    `json:"current"`
				Max             string    `json:"max"`
				Min             string    `json:"min"`
				Average         string    `json:"average"`
				NumberOfRecords int64     `json:"numberOfRecords"`
				From            time.Time `json:"from"`
				End             time.Time `json:"end"`
			} `json:"doFindAndDeleteUpdate"`
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
			DoUpdateDatasetClient struct {
				Current         string    `json:"current"`
				Max             string    `json:"max"`
				Min             string    `json:"min"`
				Average         string    `json:"average"`
				NumberOfRecords int64     `json:"numberOfRecords"`
				From            time.Time `json:"from"`
				End             time.Time `json:"end"`
			} `json:"doUpdateDatasetClient"`
			DoListUpdates struct {
				Current         string    `json:"current"`
				Max             string    `json:"max"`
				Min             string    `json:"min"`
				Average         string    `json:"average"`
				NumberOfRecords int64     `json:"numberOfRecords"`
				From            time.Time `json:"from"`
				End             time.Time `json:"end"`
			} `json:"doListUpdates"`
			DoReadDatasetClient struct {
				Current         string    `json:"current"`
				Max             string    `json:"max"`
				Min             string    `json:"min"`
				Average         string    `json:"average"`
				NumberOfRecords int64     `json:"numberOfRecords"`
				From            time.Time `json:"from"`
				End             time.Time `json:"end"`
			} `json:"doReadDatasetClient"`
			DoUpdateDatasetClientWithRecords struct {
				Current         string    `json:"current"`
				Max             string    `json:"max"`
				Min             string    `json:"min"`
				Average         string    `json:"average"`
				NumberOfRecords int64     `json:"numberOfRecords"`
				From            time.Time `json:"from"`
				End             time.Time `json:"end"`
			} `json:"doUpdateDatasetClientWithRecords"`
			DoReadDatasetClientWithRecordsUseCache struct {
				Current         string    `json:"current"`
				Max             string    `json:"max"`
				Min             string    `json:"min"`
				Average         string    `json:"average"`
				NumberOfRecords int64     `json:"numberOfRecords"`
				From            time.Time `json:"from"`
				End             time.Time `json:"end"`
			} `json:"doReadDatasetClientWithRecordsUseCache"`
			DoSaveUpdate struct {
				Current         string    `json:"current"`
				Max             string    `json:"max"`
				Min             string    `json:"min"`
				Average         string    `json:"average"`
				NumberOfRecords int64     `json:"numberOfRecords"`
				From            time.Time `json:"from"`
				End             time.Time `json:"end"`
			} `json:"doSaveUpdate"`
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
		// API Process Times
		syncAPIProcessTime, err := stringToInt64(stats.Metrics.APIProcessTime.Sync.Current)
		if err == nil {
			ssMetrics = append(ssMetrics, &metric{Type: "api_process_time_sync_ms", XValue: now.Format("2006-01-02 15:04:05"), YValue: syncAPIProcessTime})
		}
		syncRecordsAPIProcessTime, err := stringToInt64(stats.Metrics.APIProcessTime.SyncRecords.Current)
		if err == nil {
			ssMetrics = append(ssMetrics, &metric{Type: "api_process_time_syncRecords_ms", XValue: now.Format("2006-01-02 15:04:05"), YValue: syncRecordsAPIProcessTime})
		}

		// Mongodb Operation Times
		doUpdateManyDatasetClients, err := stringToInt64(stats.Metrics.MongodbOperationTime.DoUpdateManyDatasetClients.Current)
		if err == nil {
			ssMetrics = append(ssMetrics, &metric{Type: "mongodb_operation_time_doUpdateManyDatasetClients_ms", XValue: now.Format("2006-01-02 15:04:05"), YValue: doUpdateManyDatasetClients})
		}
		doListDatasetClients, err := stringToInt64(stats.Metrics.MongodbOperationTime.DoListDatasetClients.Current)
		if err == nil {
			ssMetrics = append(ssMetrics, &metric{Type: "mongodb_operation_time_doListDatasetClients_ms", XValue: now.Format("2006-01-02 15:04:05"), YValue: doListDatasetClients})
		}
		doSaveUpdate, err := stringToInt64(stats.Metrics.MongodbOperationTime.DoSaveUpdate.Current)
		if err == nil {
			ssMetrics = append(ssMetrics, &metric{Type: "mongodb_operation_time_doSaveUpdate_ms", XValue: now.Format("2006-01-02 15:04:05"), YValue: doSaveUpdate})
		}
		doUpdateDatasetClient, err := stringToInt64(stats.Metrics.MongodbOperationTime.DoUpdateDatasetClient.Current)
		if err == nil {
			ssMetrics = append(ssMetrics, &metric{Type: "mongodb_operation_time_doUpdateDatasetClient_ms", XValue: now.Format("2006-01-02 15:04:05"), YValue: doUpdateDatasetClient})
		}
		doListUpdates, err := stringToInt64(stats.Metrics.MongodbOperationTime.DoListUpdates.Current)
		if err == nil {
			ssMetrics = append(ssMetrics, &metric{Type: "mongodb_operation_time_doListUpdates_ms", XValue: now.Format("2006-01-02 15:04:05"), YValue: doListUpdates})
		}
		doReadDatasetClient, err := stringToInt64(stats.Metrics.MongodbOperationTime.DoReadDatasetClient.Current)
		if err == nil {
			ssMetrics = append(ssMetrics, &metric{Type: "mongodb_operation_time_doReadDatasetClient_ms", XValue: now.Format("2006-01-02 15:04:05"), YValue: doReadDatasetClient})
		}
		doUpdateDatasetClientWithRecords, err := stringToInt64(stats.Metrics.MongodbOperationTime.DoUpdateDatasetClientWithRecords.Current)
		if err == nil {
			ssMetrics = append(ssMetrics, &metric{Type: "mongodb_operation_time_doUpdateDatasetClientWithRecords_ms", XValue: now.Format("2006-01-02 15:04:05"), YValue: doUpdateDatasetClientWithRecords})
		}
		doReadDatasetClientWithRecordsUseCache, err := stringToInt64(stats.Metrics.MongodbOperationTime.DoReadDatasetClientWithRecordsUseCache.Current)
		if err == nil {
			ssMetrics = append(ssMetrics, &metric{Type: "mongodb_operation_time_doReadDatasetClientWithRecordsUseCache_ms", XValue: now.Format("2006-01-02 15:04:05"), YValue: doReadDatasetClientWithRecordsUseCache})
		}
		doFindAndDeleteUpdate, err := stringToInt64(stats.Metrics.MongodbOperationTime.DoFindAndDeleteUpdate.Current)
		if err == nil {
			ssMetrics = append(ssMetrics, &metric{Type: "mongodb_operation_time_doFindAndDeleteUpdate_ms", XValue: now.Format("2006-01-02 15:04:05"), YValue: doFindAndDeleteUpdate})
		}

		// sync worker stats
		ssMetrics = append(ssMetrics, &metric{Type: "sync_worker_queue_count", XValue: now.Format("2006-01-02 15:04:05"), YValue: stats.Metrics.JobQueueSize.SyncWorker.Current})
		syncProcessTime, err := stringToInt64(stats.Metrics.JobProcessTime.SyncWorker.Current)
		if err == nil {
			ssMetrics = append(ssMetrics, &metric{Type: "sync_worker_process_time_ms", XValue: now.Format("2006-01-02 15:04:05"), YValue: syncProcessTime})
		}

		// pending worker stats
		ssMetrics = append(ssMetrics, &metric{Type: "pending_worker_queue_count", XValue: now.Format("2006-01-02 15:04:05"), YValue: stats.Metrics.JobQueueSize.PendingWorker.Current})
		pendingProcessTime, err := stringToInt64(stats.Metrics.JobProcessTime.PendingWorker.Current)
		if err == nil {
			ssMetrics = append(ssMetrics, &metric{Type: "pending_worker_process_time_ms", XValue: now.Format("2006-01-02 15:04:05"), YValue: pendingProcessTime})
		}

		// ack worker stats
		ssMetrics = append(ssMetrics, &metric{Type: "ack_worker_queue_count", XValue: now.Format("2006-01-02 15:04:05"), YValue: stats.Metrics.JobQueueSize.AckWorker.Current})
		ackProcessTime, err := stringToInt64(stats.Metrics.JobProcessTime.AckWorker.Current)
		if err == nil {
			ssMetrics = append(ssMetrics, &metric{Type: "ack_worker_process_time_ms", XValue: now.Format("2006-01-02 15:04:05"), YValue: ackProcessTime})
		}
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
	ss.logger.Debugf("raw stats response %v", string(data))

	var stats statsResponse
	decoder := json.NewDecoder(bytes.NewBuffer(data))
	if err := decoder.Decode(&stats); err != nil {
		return &statsResponse{}, errors.Wrap(err, "failed to decode stats response")
	}
	ss.logger.Debugf("decoded stats response %v", stats)

	return &stats, nil
}

func stringToInt64(val string) (int64, error) {
	// string expected in format "00.00ms"
	numOnly := strings.Split(val, ".")[0]
	return strconv.ParseInt(numOnly, 10, 64)
}
