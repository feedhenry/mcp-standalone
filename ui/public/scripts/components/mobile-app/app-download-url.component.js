'use strict';

/**
 * @ngdoc component
 * @name mcp.component:app-download-url
 * @description
 * # app-download-url
 */
angular.module('mobileControlPanelApp').component('appDownloadUrl', {
  template: `<div class="app-download-url">
                <div class="url-controls" ng-hide=url>
                  <button ng-disabled="$ctrl.build.status.phase !== 'Complete'" ng-click="generateUrl()" class="btn btn-success btn-xs" type="button">Generate Download URL</button>
                  <div ng-show="$ctrl.build.status.phase === 'Complete'" class="help-block">Download URL will last 10 mins before expiring</div>
                </div>
                <div ng-hide=!url>
                  <label>Download URL: </label><a ng-if=url href="{{url}}">{{url}}</a>
                  <modal modal-class="'app-download'" class="btn-primary btn-xs" display-controls=false launch="'QR Code'" modal-open=modalOpen>
                    <qr-code content=url></qr-code>
                  </modal>
                </div>
              </div>`,
  bindings: {
    build: '<'
  },
  controller: [
    '$scope',
    'mcpApi',
    '$timeout',
    '$window',
    function($scope, mcpApi, $timeout, $window) {
      let value = $window.localStorage.getItem(
        $scope.$ctrl.build.metadata.name
      );
      value = JSON.parse(value);
      if (!value || Date.now() > value.expires) {
        $window.localStorage.removeItem($scope.$ctrl.build.metadata.name);
        $scope.url = '';
      } else {
        $scope.url = value.url;
      }

      $scope.modalOpen = false;

      let timeoutPromise = null;
      $scope.generateUrl = function() {
        mcpApi
          .mobileAppDownloadUrl($scope.$ctrl.build.metadata.name)
          .then(res => {
            $scope.url = res.url;
            $scope.$apply();
            $window.localStorage.setItem(
              $scope.$ctrl.build.metadata.name,
              JSON.stringify(res)
            );

            const timeoutPromise = $timeout(() => {
              $scope.url = '';
              $scope.modalOpen = false;
              $window.localStorage.removeItem($scope.$ctrl.build.metadata.name);
            }, res.expires - Date.now());
          });
      };

      $scope.$on('destroy', () => {
        timeoutPromise && $timeout.cancel(timeoutPromise);
      });
    }
  ]
});
