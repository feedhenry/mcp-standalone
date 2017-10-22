'use strict';

/**
 * @ngdoc component
 * @name mcp.component:mp-default-service-chart
 * @description
 * # mp-default-service-chart
 */
angular.module('mobileControlPanelApp').component('mpDefaultServiceChart', {
  template: `<div class="row">
                <div ng-repeat="chart in $ctrl.chartData">
                  <div class="col-xs-12 col-sm-12 col-md-6 col-lg-4">
                    <div class="mp-sync-trend-chart card-pf">
                      <h3>{{chart[1][0]}}</h3>
                      <mp-line-chart chart-data=chart chart-config=chartConfig></mp-line-chart>
                    </div>
                  </div>
              </div>
            </div>`,
  bindings: {
    chartData: '<'
  },
  controller: [
    '$scope',
    function($scope) {
      $scope.chartConfig = {
        axis: {
          x: {
            type: 'timeseries',
            tick: {
              format: '%H:%M:%S'
            }
          }
        },
        data: {
          x: 'x',
          xFormat: '%Y-%m-%d %H:%M:%S'
        }
      };
    }
  ]
});
