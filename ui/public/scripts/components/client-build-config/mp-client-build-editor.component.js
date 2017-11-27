'use strict';

/**
 * @ngdoc component
 * @name mcp.component:mp-client-build-editor
 * @description
 * # mp-client-build-editor
 */
angular.module('mobileControlPanelApp').component('mpClientBuildEditor', {
  template: `<div>
              <mp-create-app-config ng-if="$ctrl.ClientBuildEditorService.state === $ctrl.ClientBuildEditorService.states.CREATE" created=$ctrl.buildConfigCreated></mp-create-app-config>
              <mp-view-app-config ng-if="$ctrl.ClientBuildEditorService.state === $ctrl.ClientBuildEditorService.states.VIEW" config=$ctrl.buildConfig></mp-view-app-config>
              <mp-edit-app-config ng-if="$ctrl.ClientBuildEditorService.state === $ctrl.ClientBuildEditorService.states.EDIT" config=$ctrl.buildConfig updated=$ctrl.buildConfigUpdated cancelled=$ctrl.cancelEdit></mp-edit-app-config>
            </div>`,
  bindings: {
    buildConfig: '<',
    configCreated: '&',
    configUpdated: '&'
  },
  controller: [
    'ClientBuildEditorService',
    function(ClientBuildEditorService) {
      const ctrl = this;
      ctrl.ClientBuildEditorService = ClientBuildEditorService;

      ctrl.$onChanges = function(changes) {
        const buildConfig =
          changes.buildConfig && changes.buildConfig.currentValue;

        ClientBuildEditorService.state = buildConfig
          ? ClientBuildEditorService.states.VIEW
          : ClientBuildEditorService.states.CREATE;
      };

      ctrl.buildConfigCreated = function(appConfig) {
        ctrl.configCreated()(appConfig);
        ClientBuildEditorService.state = ClientBuildEditorService.states.VIEW;
      };

      ctrl.buildConfigUpdated = function(appConfig) {
        ctrl.configUpdated()(appConfig);
        ClientBuildEditorService.state = ClientBuildEditorService.states.VIEW;
      };

      ctrl.cancelEdit = function() {
        ClientBuildEditorService.state = ClientBuildEditorService.states.VIEW;
      };
    }
  ]
});
