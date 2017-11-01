'use strict';

/**
 * @ngdoc component
 * @name mcp.directive:mp-switch
 * @description
 * # mp-switch
 */
angular.module('mobileControlPanelApp').directive('mpSwitch', function() {
  return {
    template: `<input class="bootstrap-switch" type="checkbox">`,
    scope: {
      disabled: '=?',
      checked: '=?',
      switched: '&?'
    },
    link: function(scope, element) {
      const onText = 'Enabled';
      const offText = 'Disabled';

      element.find('input').bootstrapSwitch({
        onText: onText,
        offText: offText,
        disabled: scope.disabled,
        state: scope.checked,
        onSwitchChange: function(event, value) {
          scope.switched && scope.switched()(event, value);
        }
      });
    }
  };
});
