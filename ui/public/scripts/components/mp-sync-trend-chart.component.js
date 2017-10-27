'use strict';

/**
 * @ngdoc component
 * @name mcp.component:mp-sync-trend-chart
 * @description
 * # mp-sync-trend-chart
 */
angular.module('mobileControlPanelApp').component('mpSyncTrendChart', {
  template: `<div class="col-xs-12 col-sm-4 col-md-4">
              <div class="mp-sync-trend-chart card-pf">
                <h2 class="card-pf-heading">{{$ctrl.config.title}}</h2>
                <div>
                   <span class="count">{{$ctrl.config.total}}</span>
                   <span class="unit">{{$ctrl.config.unit}}</span>
                </div>
                <mp-sparkline-chart chart-data="$ctrl.chartData" chart-config=$ctrl.config.chart></mp-sparkline-chart>
              </div>
            </div>`,
  bindings: {
    chartData: '<',
    config: '<?'
  }
});
