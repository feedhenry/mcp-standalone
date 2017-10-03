'use strict';

/**
 * @ngdoc component
 * @name mcp.directive:modal
 * @description
 * # modal
 */
angular.module('mobileControlPanelApp').directive('modal', function($timeout) {
  return {
    template: `<div class="controlPanelAppModal">
                <div class="modal container" tabindex="-1" role="dialog" aria-labelledby="update this" aria-hidden="true">
                  <div class="modal-dialog">
                    <div class="modal-content">
                      <div class="modal-header">
                        <button type="button" class="close icon" aria-hidden="true">
                          <span class="pficon pficon-close"></span>
                        </button>
                        <h4 class="modal-title">{{modalTitle}}</h4>
                      </div>
                      <div class="modal-body">
                        <ng-transclude></ng-transclude>
                      </div>
                      <div ng-if="displayControls === undefined || displayControls === true" class="modal-footer">
                        <button type="button" class="btn btn-default cancel">Cancel</button>
                        <button type="button" class="btn btn-primary ok">Create</button>
                      </div>
                    </div>
                  </div>
                </div>
              </div>`,
    scope: {
      displayControls: '=?',
      modalOpen: '=?',
      launch: '=?',
      modalTitle: '=?',
      cancel: '&?',
      ok: '&?'
    },
    transclude: true,
    link: function(scope, element) {
      scope.modal = $('.modal.container', element).modal({
        show: false,
        keyboard: true
      });
      scope.modalOpen = scope.modalOpen || false;

      $timeout(() => {
        const okButton = $('.ok', element);
        const cancelButton = $('.cancel', element);
        const closeIcon = $('.close', element);

        okButton.on('click', () => {
          $timeout(() => {
            scope.ok && scope.ok()();
            scope.modalOpen = false;
          });
        });

        cancelButton.on('click', () => {
          $timeout(() => {
            scope.cancel && scope.cancel()();
            scope.modalOpen = false;
          });
        });

        closeIcon.on('click', () => {
          $timeout(() => {
            scope.cancel && scope.cancel()();
            scope.modalOpen = false;
          });
        });

        scope.modal.on('hidden.bs.modal', () => {
          if (!scope.modalOpen) {
            return;
          }

          $timeout(() => {
            scope.modalOpen = false;
          });
        });

        scope.$watch('modalOpen', value => {
          if (value) {
            scope.modal.modal('show');
          } else {
            scope.modal.modal('hide');
            $('.modal-backdrop').remove();
          }
        });
      });
    }
  };
});
