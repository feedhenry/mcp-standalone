'use strict';

/**
 * @ngdoc function
 * @name mobileControlPanelApp.controller:AppcreateCtrl
 * @description
 * # AppcreateCtrl
 * Controller of the mobileControlPanelApp
 */
angular.module('mobileControlPanelApp')
  .controller('AppCreateCtrl', ['$scope','mcpApi', '$location',function ($scope, mcpApi, $location) {
    $scope.app = {"clientType":""};
    $scope.createApp = function(){
      console.log("called ", $scope.app);
      mcpApi.createMobileApp($scope.app)
      .then((app)=>{
        console.log("app created");
        $location.path("apps");
      })
      .catch(err =>{
        console.error("failed to create app ", err);
      })
    }; 
  }]);
