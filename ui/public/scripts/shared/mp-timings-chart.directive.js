'use strict';

/**
 * @ngdoc component
 * @name mcp.directive:mp-timings-chart
 * @description
 * # mp-timings-chart
 */
angular.module('mobileControlPanelApp').directive('mpTimingsChart', function() {
  return {
    template: `<div class="mp-timings-container col-xs-12"></div>`,
    scope: {
      chartData: '<',
      chartConfig: '<?'
    },
    link: function(scope, element, attrs) {
      let defaultConfig = $()
        .c3ChartDefaults()
        .getDefaultLineConfig();

      defaultConfig.axis = {
        x: {
          type: 'timeseries',
          tick: {
            format: '%Y-%m-%d %H:%M:%S'
          }
        }
      };
      defaultConfig.data = {
        x: 'x',
        xFormat: '%Y-%m-%d %H:%M:%S',
        columns: scope.chartData,
        type: 'spline'
      };
      defaultConfig.point = {
        r: 2
      };
      let config = Object.assign(defaultConfig, scope.chartConfig);
      const chart = c3.generate(config);
      element.find('.mp-timings-container').append(chart.element);
    }
  };
});
