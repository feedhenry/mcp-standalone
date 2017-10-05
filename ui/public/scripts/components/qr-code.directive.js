'use strict';

/**
 * @ngdoc component
 * @name mcp.directive:qr-code
 * @description
 * # qr-code
 */
angular.module('mobileControlPanelApp').directive('qrCode', function($timeout) {
  return {
    template: `<div class="qr-code-container"></div>`,
    scope: {
      content: '<?'
    },
    link: function($scope, element, attrs) {
      const qrCodeContainer = $('.qr-code-container', element);
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
