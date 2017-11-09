'use strict';

/**
 * @ngdoc component
 * @name mcp.directive:mp-modal
 * @description
 * # mp-modal
 */
angular.module('mobileControlPanelApp').directive('mpModal', [
  '$timeout',
  function($timeout) {
    return {
      template: `<button class="btn launch" ng-class={{ngClass}}>{{launch}}</button>
              <div class="modal container" tabindex="-1" role="dialog" aria-labelledby="" aria-hidden="true">
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
              </div>`,
      scope: {
        displayControls: '=?',
        modalOpen: '=?',
        launch: '=?',
        modalTitle: '=?',
        cancel: '&?',
        ok: '&?',
        ngClass: '=?',
        modalClass: '=?'
      },
      transclude: true,
      link: function(scope, element, attrs) {
        element.removeClass();

        const modalContainer = $('.modal.container', element);
        modalContainer.addClass(scope.modalClass);
        scope.modal = modalContainer.modal({
          show: false,
          keyboard: true
        });

        scope.modalOpen = scope.modalOpen || false;

        $timeout(() => {
          const launchButton = $('.launch', element);
          const okButton = $('.ok', modalContainer);
          const cancelButton = $('.cancel', modalContainer);
          const closeIcon = $('.close', modalContainer);

          modalContainer.detach();

          launchButton.addClass(attrs.class);
          launchButton.on('click', () => {
            $timeout(() => {
              scope.modalOpen = true;
            });
          });

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
  }
]);
