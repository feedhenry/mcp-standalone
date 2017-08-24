'use strict';

/**
 * @ngdoc function
 * @name mobileControlPanelApp.controller:MobileappsCtrl
 * @description
 * # MobileappsCtrl
 * Controller of the mobileControlPanelApp
 */
angular.module('mobileControlPanelApp')
  .controller('MobileappsCtrl', ['$scope', function($scope) {
  // TODO: read mobile apps from the mobile server
  $scope.mobileapps = [{
    "name": "Mock Cordova App",
    "clientType": "cordova"
  }, {
    "name": "Mock Android App",
    "clientType": "android"
  }];
}]);
