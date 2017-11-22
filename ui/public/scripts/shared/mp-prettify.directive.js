'use strict';

/**
 * @ngdoc component
 * @name mcp.directive:mp-prettify
 * @description
 * # mp-prettify
 */
angular.module('mobileControlPanelApp').directive('mpPrettify', [
  function() {
    return {
      template: `<pre><ng-transclude></ng-transclude></pre>`,
      scope: {
        type: '<?',
        codeClass: '<?'
      },
      transclude: true,
      link: function(scope, element, attrs) {
        const innerHTML = element
          .find('span')
          .html()
          .trim();
        const prettified = prettyPrintOne(innerHTML, scope.type);
        const pre = element.find('pre');
        pre.html(prettified);
        pre.addClass(scope.codeClass);
      }
    };
  }
]);
