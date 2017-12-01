'use strict';

/**
 * @ngdoc component
 * @name mcp.component:mp-mobile-client-actions
 * @description
 * # mp-mobile-client-actions
 */
angular.module('mobileControlPanelApp').component('mpMobileClientActions', {
  template: `<div class="mp-config-actions clearfix">
              <div class="pull-right">
                <button ng-disabled="$ctrl.ClientBuildEditorService.state === ClientBuildEditorService.states.EDIT" type="button" class="btn btn-default" ng-click="$ctrl.buildClicked()">Build App</button>
                <mp-action-dropdown actions=$ctrl.dropdownActions selected=$ctrl.selected></mp-action-dropdown>
              </div>
            </div>`,
  bindings: {
    actionSelected: '&'
  },
  controller: [
    'ClientBuildEditorService',
    function(ClientBuildEditorService) {
      this.ClientBuildEditorService = ClientBuildEditorService;
      this.dropdownActions = [
        {
          label: 'Edit',
          value: ClientBuildEditorService.states.EDIT
        }
      ];

      this.selected = value => {
        this.actionSelected()(value);
      };

      this.buildClicked = () => {
        this.actionSelected()('build');
      };
    }
  ]
});
