'use strict';

/**
 * @ngdoc function
 * @name mobileControlPanelApp.controller:MobileappCtrl
 * @description
 * # MobileappCtrl
 * Controller of the mobileControlPanelApp
 */
angular.module('mobileControlPanelApp')
  .controller('MobileappCtrl', function () {
    this.mobileapp = {
      name: "mock app",
      clientType: "cordova"
    };
  });
