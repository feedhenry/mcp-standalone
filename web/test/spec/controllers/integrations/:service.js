'use strict';

describe('Controller: IntegrationsServiceCtrl', function () {

  // load the controller's module
  beforeEach(module('mobileControlPanelApp'));

  var IntegrationsServiceCtrl,
    scope;

  // Initialize the controller and a mock scope
  beforeEach(inject(function ($controller, $rootScope) {
    scope = $rootScope.$new();
    IntegrationsServiceCtrl = $controller('IntegrationsServiceCtrl', {
      $scope: scope
      // place here mocked dependencies
    });
  }));

  it('should attach a list of awesomeThings to the scope', function () {
    expect(IntegrationsServiceCtrl.awesomeThings.length).toBe(3);
  });
});
