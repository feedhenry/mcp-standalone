'use strict';

/**
 * @ngdoc function
 * @name mobileControlPanelApp.controller:IntegrationsServiceCtrl
 * @description
 * # IntegrationsServiceCtrl
 * Controller of the mobileControlPanelApp
 */
angular.module('mobileControlPanelApp')
  .controller('IntegrationCtrl', ['$scope', 'mcpApi', '$routeParams', function ($scope, mcpApi, $routeParams) {
    $scope.integrations = [];
    mcpApi.mobileService($routeParams.service, "true")
      .then(s => {
        $scope.service = s;
        $scope.integrations = Object.keys(s.integrations);
      })
      .catch(e => {
        console.error("failed to read service ", e);
      });
      mcpApi.mobileApps()
      .then((apps) => {
        $scope.mobileapps =  {};
        for(var i=0; i < apps.length; i++){
          let app = apps[i];
          $scope.mobileapps[app.clientType] = "true";
        }
        $scope.clients = Object.keys($scope.mobileapps);
      })
      .catch(e => {
        console.error(e);
      });  
      $scope.enabled = function(integration, service){
        if(!service){
          return false;
        }
        if(service.integrations[integration]){
          return service.integrations[integration].enabled == true;
        }
        return false;
        
      };
      $scope.enableIntegration = function(service){
        mcpApi.integrateService(service)
        .then((res)=>{
          console.log("Service integrated");
        })
        .catch(e=>{
          console.log("error integrating service ", e);
        })
        return true;
      };

      $scope.clientType = ($scope.clients && $scope.clients.length > 0)? $scope.clients[0] : "cordova";
      
      
      $scope.installationOpt = function(type){
        $scope.clientType = type;
      };

  }]);
