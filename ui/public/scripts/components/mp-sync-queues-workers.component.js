'use strict';

/**
 * @ngdoc component
 * @name mcp.component:mp-sync-queues-workers
 * @description
 * # mp-sync-queues-workers
 */
angular.module('mobileControlPanelApp').component('mpSyncQueuesWorkers', {
  template: `<div class="card-pf mp-sync-queues-workers">
              <div class="card-pf-heading row">
                <h2 class="col-sm-6 col-xs-6">Queues &amp; Workers</h2>
                <span class="col-sm-6 col-xs-6">Last 30 Readings</span>
              </div>
              <div class="row">
                <mp-sync-trend-chart ng-repeat="trend in trendData" chart-data=trend.chart config=trend.config></mp-sync-trend-chart>
              </div>
            </div>`,
  bindings: {
    chartData: '<'
  },
  controller: [
    '$scope',
    function($scope) {
      $scope.trendData = [];

      const metrics = {
        pending_worker_queue_count: {
          title: 'Pending Queue',
          unit: 'Items Processed ',
          supplementId: 'pending_worker_queue_count_total',
          color: '#1f77b4'
        },
        ack_worker_queue_count: {
          title: 'Ack Queue',
          unit: 'Items Processed',
          supplementId: 'ack_worker_queue_count_total',
          color: '#2ca02c'
        },
        sync_worker_queue_count: {
          title: 'Sync Queue',
          unit: 'Items Processed',
          supplementId: 'sync_worker_queue_count_total',
          color: '#ff7f0e'
        },
        pending_worker_process_time_ms: {
          title: 'Pending Worker',
          unit: 'ms avg process time',
          supplementId: 'pending_worker_process_time_ms_avg',
          color: '#1f77b4'
        },
        ack_worker_process_time_ms_avg: {
          title: 'Ack Worker',
          unit: 'ms avg process time',
          supplementId: 'ack_worker_process_time_ms_avg_avg',
          color: '#2ca02c'
        },
        sync_worker_process_time_ms: {
          title: 'Sync Worker',
          unit: 'ms avg process time',
          supplementId: 'sync_worker_process_time_ms_avg',
          color: '#ff7f0e'
        }
      };

      const chartDataReduce = function(chartData, acc, key) {
        let metricConfig = metrics[key];
        let chart = chartData.filter(chart => chart[1][0] === key).pop();
        let supplementChart = chartData
          .filter(chart => chart[1][0] === metricConfig.supplementId)
          .pop();

        if (!chart) {
          acc.push({ config: { title: metricConfig.title } });
          return acc;
        }

        let config = {
          chart: {
            data: {
              type: 'area'
            },
            point: {
              r: 0
            },
            color: {
              pattern: [metricConfig.color]
            },
            tooltip: {
              contents: function(d) {
                return (
                  '<p class="mp-sync-trend-tooltip">' + d[0].value + '</p>'
                );
              }
            }
          }
        };

        config.title = metricConfig.title;
        config.unit = metricConfig.unit;
        if (supplementChart) {
          config.total = supplementChart[1][supplementChart[1].length - 1];
        }

        acc.push({ config: config, chart: [chart[1]] });
        return acc;
      };

      $scope.trendData = Object.keys(metrics).reduce(
        chartDataReduce.bind(null, $scope.$ctrl.chartData),
        []
      );
    }
  ]
});
