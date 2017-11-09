'use strict';

/**
 * @ngdoc component
 * @name mcp.directive:mp-qrcode
 * @description
 * # mp-qrcode
 */
angular.module('mobileControlPanelApp').directive('mpQrcode', function() {
  return {
    template: `<div class="qr-code-container"></div>`,
    scope: {
      content: '<?'
    },
    link: function($scope, element, attrs) {
      const qrCodeContainer = $('.mp-qrcode-container', element);
      qrCodeContainer.qrcode({
        text: $scope.content,
        size: 250
      });

      $scope.$watch('content', () => {
        qrCodeContainer.find('canvas').remove();
        qrCodeContainer.qrcode({
          text: $scope.content,
          size: 250
        });
      });
    }
  };
});
