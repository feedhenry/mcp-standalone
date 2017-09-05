'use strict';

describe('Controller: ServiceCtrl', function () {

  // load the controller's module
  beforeEach(module('mobileControlPanelApp'));

  var ServiceCtrl,
    scope;

  // Initialize the controller and a mock scope
  beforeEach(inject(function ($controller, $rootScope) {
    scope = $rootScope.$new();
    ServiceCtrl = $controller('ServiceCtrl', {
      $scope: scope
      // place here mocked dependencies
    });
  }));

  it('should attach a list of awesomeThings to the scope', function () {
    expect(ServiceCtrl.awesomeThings.length).toBe(3);
  });
});
