'use strict';

/**
 * @ngdoc function
 * @name mobileControlPanelApp.controller:ServicesCtrl
 * @description
 * # ServicesCtrl
 * Controller of the mobileControlPanelApp
 */
angular.module('mobileControlPanelApp')
  .controller('ServicesCtrl', ['$scope','mcpApi',function ($scope, mcpApi) {
    mcpApi.mobileServices()
    .then(s=>{
      console.log("got services ", s);
      $scope.services = s;
    })
    .catch(e =>{
      console.error("error getting services ", e);
    })
  }]);
