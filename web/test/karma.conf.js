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
      'app/bower_components/jquery/dist/jquery.js',
      'app/bower_components/angular/angular.js',
      'app/bower_components/bootstrap-sass-official/assets/javascripts/bootstrap/affix.js',
      'app/bower_components/bootstrap-sass-official/assets/javascripts/bootstrap/alert.js',
      'app/bower_components/bootstrap-sass-official/assets/javascripts/bootstrap/button.js',
      'app/bower_components/bootstrap-sass-official/assets/javascripts/bootstrap/carousel.js',
      'app/bower_components/bootstrap-sass-official/assets/javascripts/bootstrap/collapse.js',
      'app/bower_components/bootstrap-sass-official/assets/javascripts/bootstrap/dropdown.js',
      'app/bower_components/bootstrap-sass-official/assets/javascripts/bootstrap/tab.js',
      'app/bower_components/bootstrap-sass-official/assets/javascripts/bootstrap/transition.js',
      'app/bower_components/bootstrap-sass-official/assets/javascripts/bootstrap/scrollspy.js',
      'app/bower_components/bootstrap-sass-official/assets/javascripts/bootstrap/modal.js',
      'app/bower_components/bootstrap-sass-official/assets/javascripts/bootstrap/tooltip.js',
      'app/bower_components/bootstrap-sass-official/assets/javascripts/bootstrap/popover.js',
      'app/bower_components/angular-route/angular-route.js',
      'app/bower_components/angular-sanitize/angular-sanitize.js',
      'app/bower_components/angular-utf8-base64/angular-utf8-base64.js',
      'app/bower_components/js-logger/src/logger.js',
      'app/bower_components/hawtio-core/dist/hawtio-core.js',
      'app/bower_components/hawtio-extension-service/dist/hawtio-extension-service.js',
      'app/bower_components/hopscotch/dist/js/hopscotch.js',
      'app/bower_components/lodash/lodash.js',
      'app/bower_components/bootstrap/dist/js/bootstrap.js',
      'app/bower_components/bootstrap-datepicker/dist/js/bootstrap-datepicker.min.js',
      'app/bower_components/bootstrap-select/dist/js/bootstrap-select.js',
      'app/bower_components/bootstrap-switch/dist/js/bootstrap-switch.js',
      'app/bower_components/bootstrap-touchspin/src/jquery.bootstrap-touchspin.js',
      'app/bower_components/d3/d3.js',
      'app/bower_components/c3/c3.js',
      'app/bower_components/datatables/media/js/jquery.dataTables.js',
      'app/bower_components/datatables-colreorder/js/dataTables.colReorder.js',
      'app/bower_components/datatables-colvis/js/dataTables.colVis.js',
      'app/bower_components/google-code-prettify/bin/prettify.min.js',
      'app/bower_components/matchHeight/dist/jquery.matchHeight.js',
      'app/bower_components/moment/moment.js',
      'app/bower_components/eonasdan-bootstrap-datetimepicker/build/js/bootstrap-datetimepicker.min.js',
      'app/bower_components/patternfly-bootstrap-combobox/js/bootstrap-combobox.js',
      'app/bower_components/patternfly-bootstrap-treeview/dist/bootstrap-treeview.js',
      'app/bower_components/patternfly/dist/js/patternfly.js',
      'app/bower_components/uri.js/src/URI.js',
      'app/bower_components/uri.js/src/URITemplate.js',
      'app/bower_components/uri.js/src/jquery.URI.js',
      'app/bower_components/uri.js/src/URI.fragmentURI.js',
      'app/bower_components/origin-web-common/dist/origin-web-common.js',
      'app/bower_components/angular-animate/angular-animate.js',
      'app/bower_components/angular-bootstrap/ui-bootstrap-tpls.js',
      'app/bower_components/jquery-ui/jquery-ui.js',
      'app/bower_components/angular-dragdrop/src/angular-dragdrop.js',
      'app/bower_components/datatables.net/js/jquery.dataTables.js',
      'app/bower_components/angularjs-datatables/dist/angular-datatables.js',
      'app/bower_components/angularjs-datatables/dist/plugins/bootstrap/angular-datatables.bootstrap.js',
      'app/bower_components/angularjs-datatables/dist/plugins/colreorder/angular-datatables.colreorder.js',
      'app/bower_components/angularjs-datatables/dist/plugins/columnfilter/angular-datatables.columnfilter.js',
      'app/bower_components/angularjs-datatables/dist/plugins/light-columnfilter/angular-datatables.light-columnfilter.js',
      'app/bower_components/angularjs-datatables/dist/plugins/colvis/angular-datatables.colvis.js',
      'app/bower_components/angularjs-datatables/dist/plugins/fixedcolumns/angular-datatables.fixedcolumns.js',
      'app/bower_components/angularjs-datatables/dist/plugins/fixedheader/angular-datatables.fixedheader.js',
      'app/bower_components/angularjs-datatables/dist/plugins/scroller/angular-datatables.scroller.js',
      'app/bower_components/angularjs-datatables/dist/plugins/tabletools/angular-datatables.tabletools.js',
      'app/bower_components/angularjs-datatables/dist/plugins/buttons/angular-datatables.buttons.js',
      'app/bower_components/angularjs-datatables/dist/plugins/select/angular-datatables.select.js',
      'app/bower_components/angular-drag-and-drop-lists/angular-drag-and-drop-lists.js',
      'app/bower_components/datatables.net-select/js/dataTables.select.js',
      'app/bower_components/angular-patternfly/dist/angular-patternfly.js',
      'app/bower_components/urijs/src/URI.js',
      'app/bower_components/angular-mocks/angular-mocks.js',
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
