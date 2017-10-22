'use strict';

/**
 * @ngdoc component
 * @name mcp.component:mp-keycloak-chart
 * @description
 * # mp-keycloak-chart
 */
angular.module('mobileControlPanelApp').component('mpKeycloakChart', {
  template: `<div class="col-xs-6 col-md-6 col-lg-6">
              <div class="card-pf card-pf-accented card-pf-aggregate-status">
                <h2 class="card-pf-title">
                  <span class="fa fa-shield"></span><span class="card-pf-aggregate-status-count">{{metrics.logins.total}}</span> Logins
                </h2>
                <div class="card-pf-body">
                  <p class="card-pf-aggregate-status-notifications">
                    <span class="card-pf-aggregate-status-notification"><span class="pficon pficon-ok"></span>{{metrics.logins.success}}</span>
                    <span class="card-pf-aggregate-status-notification"><span class="pficon pficon-error-circle-o"></span>{{metrics.logins.error}}</span>
                  </p>
                </div>
              </div>
            </div>
            <div ng-if-end class="col-xs-6 col-md-6 col-lg-6">
              <div class="card-pf card-pf-accented card-pf-aggregate-status">
                <h2 class="card-pf-title">
                  <span class="fa fa-shield"></span><span class="card-pf-aggregate-status-count">{{metrics.registrations.total}}</span> Registrations
                </h2>
                <div class="card-pf-body">
                  <p class="card-pf-aggregate-status-notifications">
                    <span class="card-pf-aggregate-status-notification"><span class="pficon pficon-ok"></span>{{metrics.registrations.success}}</span>
                    <span class="card-pf-aggregate-status-notification"><span class="pficon pficon-error-circle-o"></span>{{metrics.registrations.error}}</span>
                  </p>
                </div>
              </div>
            </div>`,
  bindings: {
    chartData: '<'
  },
  controller: [
    '$scope',
    function($scope) {
      $scope.metrics = {
        logins: {
          success: 0,
          error: 0
        },
        registrations: {
          success: 0,
          error: 0
        }
      };

      var loginSuccessChart = _.findWhere($scope.$ctrl.chartData, {
        title: 'LOGIN'
      });
      if (loginSuccessChart) {
        $scope.metrics.logins.success =
          loginSuccessChart.data.columns[1][
            loginSuccessChart.data.columns[1].length - 1
          ];
      }
      var loginErrorChart = _.findWhere($scope.$ctrl.chartData, {
        title: 'LOGIN_ERROR'
      });
      if (loginErrorChart) {
        $scope.metrics.logins.error =
          loginErrorChart.data.columns[1][
            loginErrorChart.data.columns[1].length - 1
          ];
      }
      var registrationSuccessChart = _.findWhere($scope.$ctrl.chartData, {
        title: 'REGISTER'
      });
      if (registrationSuccessChart) {
        $scope.metrics.registrations.success =
          registrationSuccessChart.data.columns[1][
            registrationSuccessChart.data.columns[1].length - 1
          ];
      }
      var registrationErrorChart = _.findWhere($scope.$ctrl.chartData, {
        title: 'REGISTER_ERROR'
      });
      if (registrationErrorChart) {
        $scope.metrics.registrations.error =
          registrationErrorChart.data.columns[1][
            registrationErrorChart.data.columns[1].length - 1
          ];
      }

      $scope.metrics.logins.total =
        $scope.metrics.logins.success + $scope.metrics.logins.error;
      $scope.metrics.registrations.total =
        $scope.metrics.registrations.success +
        $scope.metrics.registrations.error;
      $scope.metrics = $scope.metrics;
    }
  ]
});
