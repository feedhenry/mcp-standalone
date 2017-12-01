'use strict';

/**
 * @ngdoc component
 * @name mcp.component:mp-app-download
 * @description
 * # mp-app-download
 */
angular.module('mobileControlPanelApp').component('mpAppDownload', {
  template: `<div class="mp-app-download">
                <div class="url-controls" ng-hide=$ctrl.url>
                  <button ng-disabled="$ctrl.build.status.phase !== 'Complete'" ng-click="$ctrl.generateUrl()" class="btn btn-success btn-xs" type="button">Generate Download URL</button>
                  <div ng-show="$ctrl.build.status.phase === 'Complete'" class="help-block">Download URL will last 30 mins before expiring</div>
                </div>
                <div ng-hide=!$ctrl.url>
                  <label>Download URL: </label><a ng-if=$ctrl.url href="{{$ctrl.url}}">{{$ctrl.url}}</a>
                  <mp-modal modal-class="'mp-app-download-modal'" class="btn-primary btn-xs" launch="'QR Code'" modal-open=$ctrl.modalOpen>
                    <p class="help-block" >Scan the QR code to install this build directly onto a device</p>
                    <mp-qrcode content=$ctrl.url></mp-qrcode>
                  </mp-modal>
                </div>
              </div>`,
  bindings: {
    build: '<'
  },
  controller: [
    '$scope',
    'McpService',
    '$timeout',
    '$window',
    function($scope, McpService, $timeout, $window) {
      let storedValue = $window.localStorage.getItem(this.build.metadata.name);
      storedValue = JSON.parse(storedValue);
      let timeoutPromise = null;
      this.modalOpen = false;

      function timeoutFn() {
        this.url = '';
        this.modalOpen = false;
        $window.localStorage.removeItem(this.build.metadata.name);
      }

      let dateNow = Date.now();
      if (!storedValue || dateNow > storedValue.expires) {
        $window.localStorage.removeItem(this.build.metadata.name);
        this.url = '';
      } else {
        timeoutPromise = $timeout(
          timeoutFn.bind(this),
          storedValue.expires - dateNow
        );
        this.url = storedValue.url;
      }

      this.generateUrl = function() {
        McpService.mobileAppDownloadUrl(this.build.metadata.name).then(res => {
          this.url = res.url;
          $window.localStorage.setItem(
            this.build.metadata.name,
            JSON.stringify(res)
          );
          timeoutPromise = $timeout(
            timeoutFn.bind(this),
            res.expires - Date.now()
          );
        });
      };

      $scope.$on('$destroy', () => {
        timeoutPromise && $timeout.cancel(timeoutPromise);
      });
    }
  ]
});
