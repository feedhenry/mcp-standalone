'use strict';

/**
 * @ngdoc component
 * @name mcp.directive:mp-sparkline-chart
 * @description
 * # mp-sparkline-chart
 */
angular
  .module('mobileControlPanelApp')
  .directive('mpSparklineChart', function() {
    return {
      template: `<div class="mp-sparkline-container">
                  <div ng-if="noData" class="empty-chart-content">
                    <span class="pficon pficon-info"></span><span>No data available</span>
                  </div>
                </div>`,
      scope: {
        chartData: '<',
        chartConfig: '<?'
      },
      link: function(scope, element, attrs) {
        if (!scope.chartData) {
          scope.noData = true;
          return;
        }

        let defaultConfig = $()
          .c3ChartDefaults()
          .getDefaultSparklineConfig();

        const config = Object.assign(defaultConfig, scope.chartConfig);
        config.data = config.data ? config.data : {};
        config.data.columns = scope.chartData;
        const chart = c3.generate(config);
        element.find('.mp-sparkline-container').append(chart.element);
      }
    };
  });
