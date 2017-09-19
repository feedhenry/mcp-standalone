package metrics

import (
	"testing"

	"github.com/Sirupsen/logrus"
	"github.com/feedhenry/mcp-standalone/pkg/mobile"

	"net/http"

	"bytes"
	"io/ioutil"

	"github.com/feedhenry/mcp-standalone/pkg/mock"
	"github.com/pkg/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/kubernetes/fake"
	v1 "k8s.io/client-go/pkg/api/v1"
	ktesting "k8s.io/client-go/testing"
)

func TestFhSyncServer_Gather(t *testing.T) {
	cases := []struct {
		Name        string
		ExpectError bool
		Client      func() kubernetes.Interface
		Validate    func(t *testing.T, metrics []*metric)
		Requester   func(t *testing.T) mobile.ExternalHTTPRequester
	}{
		{
			Name: "test gather gathers as expected",
			Client: func() kubernetes.Interface {
				client := &fake.Clientset{}
				client.AddReactor("list", "secrets", func(action ktesting.Action) (handled bool, ret runtime.Object, err error) {
					return true, &v1.SecretList{
						Items: []v1.Secret{
							{
								ObjectMeta: metav1.ObjectMeta{
									Name: "fh-sync-server",
								},
								Data: map[string][]byte{
									"uri":  []byte("http://fh-sync-server.192.168.37.1.nip.io"),
									"name": []byte("fh-sync-server"),
									"type": []byte("fh-sync-server"),
								},
							},
						},
					}, nil
				})
				return client
			},
			Requester: func(t *testing.T) mobile.ExternalHTTPRequester {
				return &mock.Requester{
					Test: t,
					Responder: func(host string, path string, method string, t *testing.T) (*http.Response, error) {
						if path == "/sys/info/stats" {
							bod := bytes.NewReader([]byte(mockResponse))
							return &http.Response{StatusCode: 200, Body: ioutil.NopCloser(bod)}, nil
						}

						return nil, errors.New("unknown path " + path + " don't know how to respond")
					},
				}
			},
			Validate: func(t *testing.T, metrics []*metric) {
				if len(metrics) == 0 {
					t.Fatal("expected some metrics but got none")
				}
				expectedMetricsCount := 17
				metricsCount := 0
				for _, m := range metrics {
					if m.Type == "api_process_time_sync_ms" {
						if m.YValue != 2 {
							t.Fatalf("expected the value of api_process_time_sync_ms to be 2 but got %v", m.YValue)
						}
						metricsCount++
					}
					if m.Type == "api_process_time_syncRecords_ms" {
						if m.YValue != 5 {
							t.Fatalf("expected the value of api_process_time_syncRecords_ms to be 5 but got %v", m.YValue)
						}
						metricsCount++
					}
					if m.Type == "mongodb_operation_time_doUpdateManyDatasetClients_ms" {
						if m.YValue != 1 {
							t.Fatalf("expected the value of mongodb_operation_time_doUpdateManyDatasetClients_ms to be 1 but got %v", m.YValue)
						}
						metricsCount++
					}
					if m.Type == "mongodb_operation_time_doListDatasetClients_ms" {
						if m.YValue != 7 {
							t.Fatalf("expected the value of mongodb_operation_time_doListDatasetClients_ms to be 7 but got %v", m.YValue)
						}
						metricsCount++
					}
					if m.Type == "mongodb_operation_time_doSaveUpdate_ms" {
						if m.YValue != 8 {
							t.Fatalf("expected the value of mongodb_operation_time_doSaveUpdate_ms to be 8 but got %v", m.YValue)
						}
						metricsCount++
					}
					if m.Type == "mongodb_operation_time_doUpdateDatasetClient_ms" {
						if m.YValue != 2 {
							t.Fatalf("expected the value of mongodb_operation_time_doUpdateDatasetClient_ms to be 2 but got %v", m.YValue)
						}
						metricsCount++
					}
					if m.Type == "mongodb_operation_time_doListUpdates_ms" {
						if m.YValue != 1 {
							t.Fatalf("expected the value of mongodb_operation_time_doListUpdates_ms to be 1 but got %v", m.YValue)
						}
						metricsCount++
					}
					if m.Type == "mongodb_operation_time_doReadDatasetClient_ms" {
						if m.YValue != 15 {
							t.Fatalf("expected the value of mongodb_operation_time_doReadDatasetClient_ms to be 15 but got %v", m.YValue)
						}
						metricsCount++
					}
					if m.Type == "mongodb_operation_time_doUpdateDatasetClientWithRecords_ms" {
						if m.YValue != 4 {
							t.Fatalf("expected the value of mongodb_operation_time_doUpdateDatasetClientWithRecords_ms to be 4 but got %v", m.YValue)
						}
						metricsCount++
					}
					if m.Type == "mongodb_operation_time_doReadDatasetClientWithRecordsUseCache_ms" {
						if m.YValue != 9 {
							t.Fatalf("expected the value of mongodb_operation_time_doReadDatasetClientWithRecordsUseCache_ms to be 9 but got %v", m.YValue)
						}
						metricsCount++
					}
					if m.Type == "mongodb_operation_time_doFindAndDeleteUpdate_ms" {
						if m.YValue != 13 {
							t.Fatalf("expected the value of mongodb_operation_time_doFindAndDeleteUpdate_ms to be 13 but got %v", m.YValue)
						}
						metricsCount++
					}
					if m.Type == "sync_worker_queue_count" {
						if m.YValue != 121 {
							t.Fatalf("expected the value of sync_worker_queue_count to be 121 but got %v", m.YValue)
						}
						metricsCount++
					}
					if m.Type == "sync_worker_process_time_ms" {
						if m.YValue != 9 {
							t.Fatalf("expected the value of sync_worker_process_time_ms to be 9 but got %v", m.YValue)
						}
						metricsCount++
					}
					if m.Type == "pending_worker_queue_count" {
						if m.YValue != 99 {
							t.Fatalf("expected the value of pending_worker_queue_count to be 99 but got %v", m.YValue)
						}
						metricsCount++
					}
					if m.Type == "pending_worker_process_time_ms" {
						if m.YValue != 99 {
							t.Fatalf("expected the value of pending_worker_process_time_ms to be 99 but got %v", m.YValue)
						}
						metricsCount++
					}
					if m.Type == "ack_worker_queue_count" {
						if m.YValue != 55 {
							t.Fatalf("expected the value of ack_worker_queue_count to be 55 but got %v", m.YValue)
						}
						metricsCount++
					}
					if m.Type == "ack_worker_process_time_ms" {
						if m.YValue != 999 {
							t.Fatalf("expected the value of ack_worker_process_time_ms to be 999 but got %v", m.YValue)
						}
						metricsCount++
					}
				}
				if metricsCount != expectedMetricsCount {
					t.Fatalf("expected to have %v metrics but got %v", expectedMetricsCount, metricsCount)
				}
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.Name, func(t *testing.T) {
			httpClientBuilder := &mock.HttpClientBuilder{Requester: tc.Requester(t)}
			ss := NewFhSyncServer(httpClientBuilder, buildDefaultTestTokenClientBuilder(tc.Client()), logrus.StandardLogger())
			metrics, err := ss.Gather()
			if err == nil && tc.ExpectError {
				t.Fatal("expected an error but got none")
			}
			if err != nil && !tc.ExpectError {
				t.Fatalf("did not expect an error but got one %s", err)
			}
			if tc.Validate != nil {
				tc.Validate(t, metrics)
			}
		})
	}
}

var mockResponse = `{
  "metrics": {
    "CPU usage": {
      "master": {
        "current": "1.00%",
        "max": "7.99%",
        "min": "0.00%",
        "average": "2.97%",
        "numberOfRecords": 1000,
        "from": "2017-09-19T09:59:49.170Z",
        "end": "2017-09-19T11:27:22.293Z"
      }
    },
    "RSS Memory Usage": {
      "master": {
        "current": "52.37MB",
        "max": "89.99MB",
        "min": "49.50MB",
        "average": "61.64MB",
        "numberOfRecords": 1000,
        "from": "2017-09-19T09:59:48.168Z",
        "end": "2017-09-19T11:27:21.290Z"
      }
    },
    "Job Process Time": {
      "sync_worker": {
        "current": "9.00ms",
        "max": "12.00ms",
        "min": "4.00ms",
        "average": "8.21ms",
        "numberOfRecords": 19,
        "from": "2017-09-19T11:24:01.595Z",
        "end": "2017-09-19T11:27:19.852Z"
      },
      "pending_worker": {
        "current": "99.00ms",
        "max": "12.00ms",
        "min": "4.00ms",
        "average": "8.21ms",
        "numberOfRecords": 19,
        "from": "2017-09-19T11:24:01.595Z",
        "end": "2017-09-19T11:27:19.852Z"
      },
      "ack_worker": {
        "current": "999.00ms",
        "max": "12.00ms",
        "min": "4.00ms",
        "average": "8.21ms",
        "numberOfRecords": 19,
        "from": "2017-09-19T11:24:01.595Z",
        "end": "2017-09-19T11:27:19.852Z"
      }
    },
    "Job Queue Size": {
      "pending_worker": {
        "current": 99,
        "max": 0,
        "min": 0,
        "average": 0,
        "numberOfRecords": 40,
        "from": "2017-09-19T11:24:06.424Z",
        "end": "2017-09-19T11:27:21.455Z"
      },
      "ack_worker": {
        "current": 55,
        "max": 0,
        "min": 0,
        "average": 0,
        "numberOfRecords": 40,
        "from": "2017-09-19T11:24:06.424Z",
        "end": "2017-09-19T11:27:21.455Z"
      },
      "sync_worker": {
        "current": 121,
        "max": 1,
        "min": 0,
        "average": 0.05,
        "numberOfRecords": 40,
        "from": "2017-09-19T11:24:06.424Z",
        "end": "2017-09-19T11:27:21.454Z"
      }
    },
    "API Process Time": {
      "sync": {
        "current": "2.00ms",
        "max": "116.00ms",
        "min": "1.00ms",
        "average": "5.57ms",
        "numberOfRecords": 114,
        "from": "2017-09-19T09:23:53.346Z",
        "end": "2017-09-19T11:27:12.247Z"
      },
      "syncRecords": {
        "current": "5.00ms",
        "max": "8.00ms",
        "min": "3.00ms",
        "average": "4.40ms",
        "numberOfRecords": 5,
        "from": "2017-09-19T09:23:56.519Z",
        "end": "2017-09-19T11:15:45.249Z"
      }
    },
    "Mongodb Operation Time": {
      "doUpdateManyDatasetClients": {
        "current": "1.00ms",
        "max": "11.00ms",
        "min": "0.00ms",
        "average": "0.60ms",
        "numberOfRecords": 613,
        "from": "2017-09-19T11:24:47.226Z",
        "end": "2017-09-19T11:27:22.330Z"
      },
      "doListDatasetClients": {
        "current": "7.00ms",
        "max": "9.00ms",
        "min": "0.00ms",
        "average": "0.78ms",
        "numberOfRecords": 306,
        "from": "2017-09-19T11:24:47.729Z",
        "end": "2017-09-19T11:27:22.328Z"
      },
			"doSaveUpdate": {
        "current": "8.00ms",
        "max": "9.00ms",
        "min": "0.00ms",
        "average": "0.78ms",
        "numberOfRecords": 306,
        "from": "2017-09-19T11:24:47.729Z",
        "end": "2017-09-19T11:27:22.328Z"
			},
      "doUpdateDatasetClient": {
        "current": "2.00ms",
        "max": "3.00ms",
        "min": "0.00ms",
        "average": "1.46ms",
        "numberOfRecords": 41,
        "from": "2017-09-19T11:24:49.249Z",
        "end": "2017-09-19T11:27:19.852Z"
      },
      "doUpdateDatasetClientWithRecords": {
        "current": "4.00ms",
        "max": "7.00ms",
        "min": "2.00ms",
        "average": "3.50ms",
        "numberOfRecords": 14,
        "from": "2017-09-19T11:24:56.927Z",
        "end": "2017-09-19T11:27:19.850Z"
      },
      "doListUpdates": {
        "current": "1.00ms",
        "max": "2.00ms",
        "min": "0.00ms",
        "average": "0.77ms",
        "numberOfRecords": 13,
        "from": "2017-09-19T11:24:49.248Z",
        "end": "2017-09-19T11:27:12.247Z"
			},
      "doReadDatasetClientWithRecordsUseCache": {
        "current": "9.00ms",
        "max": "2.00ms",
        "min": "0.00ms",
        "average": "0.77ms",
        "numberOfRecords": 13,
        "from": "2017-09-19T11:24:49.248Z",
        "end": "2017-09-19T11:27:12.247Z"
      },
      "doReadDatasetClient": {
        "current": "15.00ms",
        "max": "2.00ms",
        "min": "0.00ms",
        "average": "1.00ms",
        "numberOfRecords": 13,
        "from": "2017-09-19T11:24:49.247Z",
        "end": "2017-09-19T11:27:12.246Z"
      },
      "doFindAndDeleteUpdate": {
        "current": "13.00ms",
        "max": "2.00ms",
        "min": "0.00ms",
        "average": "1.00ms",
        "numberOfRecords": 13,
        "from": "2017-09-19T11:24:49.247Z",
        "end": "2017-09-19T11:27:12.246Z"
      }
    }
  }
}`
