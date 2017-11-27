'use strict';

/**
 * @ngdoc service
 * @name mobileControlPanelApp.ClientBuildEditorService
 * @description
 * # ClientBuildEditorService
 * ClientBuildEditorService
 */
angular.module('mobileControlPanelApp').service('ClientBuildEditorService', [
  function() {
    const states = {
      VIEW: 'view',
      EDIT: 'edit',
      CREATE: 'create'
    };

    this.states = states;
    this.state = '';
  }
]);
