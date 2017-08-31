'use strict';

/**
 * @ngdoc function
 * @name mobileControlPanelApp.controller:MobileappsCtrl
 * @description
 * # MobileappsCtrl
 * Controller of the mobileControlPanelApp
 */
angular.module('mobileControlPanelApp')
  .controller('MobileappsCtrl', ['$scope', 'mcpApi', '$location', function ($scope, mcpApi, $location) {
    $scope.mobileapps = [];
    $scope.services = [];
    mcpApi.mobileApps()
      .then((apps) => {
        console.log("apps", apps);
        $scope.mobileapps = apps;
      })
      .catch(e => {
        console.error(e);
      });

      mcpApi.mobileServices()
      .then(s=>{
        console.log("got services ", s);
        $scope.services = s;
      })
      .catch(e =>{
        console.error("error getting services ", e);
      })

    $scope.openApp = function (id) {
      console.log("open app ", id);
      $location.path("/apps/" + id);
    };
  }]);
