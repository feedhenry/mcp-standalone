'use strict';

describe('Controller: ErrorCtrl', function () {

  // load the controller's module
  beforeEach(module('mobileControlPanelApp'));

  var ErrorCtrl,
    scope;

  // Initialize the controller and a mock scope
  beforeEach(inject(function ($controller, $rootScope) {
    scope = $rootScope.$new();
    ErrorCtrl = $controller('ErrorCtrl', {
      $scope: scope
      // place here mocked dependencies
    });
  }));

  it('should attach a list of awesomeThings to the scope', function () {
    expect(ErrorCtrl.awesomeThings.length).toBe(3);
  });
});
