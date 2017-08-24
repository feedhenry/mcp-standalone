// Karma configuration
// Generated on 2017-08-20

module.exports = function(config) {
  'use strict';

  config.set({
    // enable / disable watching file and executing tests whenever any file changes
    autoWatch: true,

    // base path, that will be used to resolve files and exclude
    basePath: '../',

    // testing framework to use (jasmine/mocha/qunit/...)
    // as well as any additional frameworks (requirejs/chai/sinon/...)
    frameworks: [
      'jasmine'
    ],

    // list of files / patterns to load in the browser
    files: [
      // bower:js
      'bower_components/jquery/dist/jquery.js',
      'bower_components/angular/angular.js',
      'bower_components/bootstrap-sass-official/assets/javascripts/bootstrap/affix.js',
      'bower_components/bootstrap-sass-official/assets/javascripts/bootstrap/alert.js',
      'bower_components/bootstrap-sass-official/assets/javascripts/bootstrap/button.js',
      'bower_components/bootstrap-sass-official/assets/javascripts/bootstrap/carousel.js',
      'bower_components/bootstrap-sass-official/assets/javascripts/bootstrap/collapse.js',
      'bower_components/bootstrap-sass-official/assets/javascripts/bootstrap/dropdown.js',
      'bower_components/bootstrap-sass-official/assets/javascripts/bootstrap/tab.js',
      'bower_components/bootstrap-sass-official/assets/javascripts/bootstrap/transition.js',
      'bower_components/bootstrap-sass-official/assets/javascripts/bootstrap/scrollspy.js',
      'bower_components/bootstrap-sass-official/assets/javascripts/bootstrap/modal.js',
      'bower_components/bootstrap-sass-official/assets/javascripts/bootstrap/tooltip.js',
      'bower_components/bootstrap-sass-official/assets/javascripts/bootstrap/popover.js',
      'bower_components/angular-route/angular-route.js',
      'bower_components/angular-sanitize/angular-sanitize.js',
      'bower_components/angular-utf8-base64/angular-utf8-base64.js',
      'bower_components/js-logger/src/logger.js',
      'bower_components/hawtio-core/dist/hawtio-core.js',
      'bower_components/hawtio-extension-service/dist/hawtio-extension-service.js',
      'bower_components/hopscotch/dist/js/hopscotch.js',
      'bower_components/sifter/sifter.js',
      'bower_components/microplugin/src/microplugin.js',
      'bower_components/selectize/dist/js/selectize.js',
      'bower_components/lodash/lodash.js',
      'bower_components/kubernetes-label-selector/labelSelector.js',
      'bower_components/kubernetes-label-selector/labelFilter.js',
      'bower_components/bootstrap/dist/js/bootstrap.js',
      'bower_components/bootstrap-datepicker/dist/js/bootstrap-datepicker.min.js',
      'bower_components/bootstrap-select/dist/js/bootstrap-select.js',
      'bower_components/bootstrap-switch/dist/js/bootstrap-switch.js',
      'bower_components/bootstrap-touchspin/src/jquery.bootstrap-touchspin.js',
      'bower_components/d3/d3.js',
      'bower_components/c3/c3.js',
      'bower_components/datatables/media/js/jquery.dataTables.js',
      'bower_components/datatables-colreorder/js/dataTables.colReorder.js',
      'bower_components/datatables-colvis/js/dataTables.colVis.js',
      'bower_components/google-code-prettify/bin/prettify.min.js',
      'bower_components/matchHeight/dist/jquery.matchHeight.js',
      'bower_components/moment/moment.js',
      'bower_components/eonasdan-bootstrap-datetimepicker/build/js/bootstrap-datetimepicker.min.js',
      'bower_components/patternfly-bootstrap-combobox/js/bootstrap-combobox.js',
      'bower_components/patternfly-bootstrap-treeview/dist/bootstrap-treeview.js',
      'bower_components/patternfly/dist/js/patternfly.js',
      'bower_components/uri.js/src/URI.js',
      'bower_components/uri.js/src/URITemplate.js',
      'bower_components/uri.js/src/jquery.URI.js',
      'bower_components/uri.js/src/URI.fragmentURI.js',
      'bower_components/origin-web-common/dist/origin-web-common.js',
      'bower_components/angular-animate/angular-animate.js',
      'bower_components/angular-bootstrap/ui-bootstrap-tpls.js',
      'bower_components/jquery-ui/jquery-ui.js',
      'bower_components/angular-dragdrop/src/angular-dragdrop.js',
      'bower_components/datatables.net/js/jquery.dataTables.js',
      'bower_components/angularjs-datatables/dist/angular-datatables.js',
      'bower_components/angularjs-datatables/dist/plugins/bootstrap/angular-datatables.bootstrap.js',
      'bower_components/angularjs-datatables/dist/plugins/colreorder/angular-datatables.colreorder.js',
      'bower_components/angularjs-datatables/dist/plugins/columnfilter/angular-datatables.columnfilter.js',
      'bower_components/angularjs-datatables/dist/plugins/light-columnfilter/angular-datatables.light-columnfilter.js',
      'bower_components/angularjs-datatables/dist/plugins/colvis/angular-datatables.colvis.js',
      'bower_components/angularjs-datatables/dist/plugins/fixedcolumns/angular-datatables.fixedcolumns.js',
      'bower_components/angularjs-datatables/dist/plugins/fixedheader/angular-datatables.fixedheader.js',
      'bower_components/angularjs-datatables/dist/plugins/scroller/angular-datatables.scroller.js',
      'bower_components/angularjs-datatables/dist/plugins/tabletools/angular-datatables.tabletools.js',
      'bower_components/angularjs-datatables/dist/plugins/buttons/angular-datatables.buttons.js',
      'bower_components/angularjs-datatables/dist/plugins/select/angular-datatables.select.js',
      'bower_components/angular-drag-and-drop-lists/angular-drag-and-drop-lists.js',
      'bower_components/datatables.net-select/js/dataTables.select.js',
      'bower_components/angular-patternfly/dist/angular-patternfly.js',
      'bower_components/urijs/src/URI.js',
      'bower_components/angular-mocks/angular-mocks.js',
      // endbower
      'app/scripts/**/*.js',
      'test/mock/**/*.js',
      'test/spec/**/*.js'
    ],

    // list of files / patterns to exclude
    exclude: [
    ],

    // web server port
    port: 8080,

    // Start these browsers, currently available:
    // - Chrome
    // - ChromeCanary
    // - Firefox
    // - Opera
    // - Safari (only Mac)
    // - PhantomJS
    // - IE (only Windows)
    browsers: [
      'PhantomJS'
    ],

    // Which plugins to enable
    plugins: [
      'karma-phantomjs-launcher',
      'karma-jasmine'
    ],

    // Continuous Integration mode
    // if true, it capture browsers, run tests and exit
    singleRun: false,

    colors: true,

    // level of logging
    // possible values: LOG_DISABLE || LOG_ERROR || LOG_WARN || LOG_INFO || LOG_DEBUG
    logLevel: config.LOG_INFO,

    // Uncomment the following lines if you are using grunt's server to run the tests
    // proxies: {
    //   '/': 'http://localhost:9000/'
    // },
    // URL root prevent conflicts with the site root
    // urlRoot: '_karma_'
  });
};
