// Generated on 2017-08-30 using generator-angular 0.16.0
'use strict';

var ADD_CONFIG_LOCAL = false;

// # Globbing
// for performance reasons we're only matching one level down:
// 'test/spec/{,*/}*.js'
// use this if you want to recursively match all subfolders:
// 'test/spec/**/*.js'
module.exports = function(grunt) {
  // Time how long tasks take. Can help when optimizing build times
  require('time-grunt')(grunt);
  // grunt.loadNpmTasks('grunt-contrib-less');
  // Automatically load required Grunt tasks
  require('jit-grunt')(grunt, {
    useminPrepare: 'grunt-usemin'
  });

  // Configurable paths for the application
  var appConfig = {
    app: require('./bower.json').appPath || 'app',
    dist: 'dist'
  };

  // Define the configuration for all the tasks
  grunt.initConfig({
    // Project settings
    yeoman: appConfig,

    less: {
      dist: {
        options: {
          compress: true,
          yuicompress: true,
          optimization: 2
        },
        files: {
          '.tmp/styles/mcp.css': appConfig.app + '/styles/main.less'
        }
      }
    },

    // Watches files for changes and runs tasks based on the changed files
    watch: {
      bower: {
        files: ['bower.json'],
        tasks: ['wiredep']
      },
      js: {
        files: ['<%= yeoman.app %>/scripts/{,*/}*.js'],
        tasks: ['setlocalconfig', 'build']
      },
      jsTest: {
        files: ['test/spec/{,*/}*.js'],
        tasks: ['karma']
      },
      less: {
        files: ['<%= yeoman.app %>/styles/{,*/}*.less'],
        tasks: ['setlocalconfig', 'build']
      },
      gruntfile: {
        files: ['Gruntfile.js']
      }
    },

    // Empties folders to start fresh
    clean: {
      dist: {
        files: [
          {
            dot: true,
            src: [
              '.tmp',
              '<%= yeoman.dist %>/{,*/}*',
              '!<%= yeoman.dist %>/.git{,*/}*'
            ]
          }
        ]
      },
      server: '.tmp'
    },

    // Add vendor prefixed styles
    postcss: {
      options: {
        processors: [
          require('autoprefixer-core')({ browsers: ['last 1 version'] })
        ]
      },
      server: {
        options: {
          map: true
        },
        files: [
          {
            expand: true,
            cwd: '.tmp/styles/',
            src: '{,*/}*.css',
            dest: '.tmp/styles/'
          }
        ]
      },
      dist: {
        files: [
          {
            expand: true,
            cwd: '.tmp/styles/',
            src: '{,*/}*.css',
            dest: '.tmp/styles/'
          }
        ]
      }
    },

    // Automatically inject Bower components into the app
    wiredep: {
      app: {
        src: ['<%= yeoman.app %>/index.html'],
        ignorePath: /\.\.\//
      },
      test: {
        devDependencies: true,
        src: '<%= karma.unit.configFile %>',
        ignorePath: /\.\.\//,
        fileTypes: {
          js: {
            block: /(([\s\t]*)\/{2}\s*?bower:\s*?(\S*))(\n|\r|.)*?(\/{2}\s*endbower)/gi,
            detect: {
              js: /'(.*\.js)'/gi
            },
            replace: {
              js: "'{{filePath}}',"
            }
          }
        }
      }
    },

    // Renames files for browser caching purposes
    filerev: {
      dist: {
        src: [
          '<%= yeoman.dist %>/scripts/{,*/}*.js',
          '<%= yeoman.dist %>/styles/{,*/}*.css',
          '<%= yeoman.dist %>/images/{,*/}*.{png,jpg,jpeg,gif,webp,svg}',
          '<%= yeoman.dist %>/styles/fonts/*'
        ]
      }
    },

    // Reads HTML for usemin blocks to enable smart builds that automatically
    // concat, minify and revision files. Creates configurations in memory so
    // additional tasks can operate on them
    useminPrepare: {
      html: '<%= yeoman.app %>/index.html',
      options: {
        root: '<%= yeoman.app %>',
        dest: '<%= yeoman.dist %>',
        flow: {
          html: {
            steps: {
              js: ['concat', 'uglifyjs'],
              css: ['cssmin']
            },
            post: {}
          }
        }
      }
    },

    // Performs rewrites based on filerev and the useminPrepare configuration
    usemin: {
      html: ['<%= yeoman.dist %>/{,*/}*.html'],
      css: ['.tmp/styles/{,*/}*.css'],
      js: ['<%= yeoman.dist %>/scripts/{,*/}*.js'],
      options: {
        assetsDirs: [
          '<%= yeoman.dist %>',
          '<%= yeoman.dist %>/images',
          '<%= yeoman.dist %>/styles'
        ],
        patterns: {
          js: [
            [
              /(images\/[^''""]*\.(png|jpg|jpeg|gif|webp|svg))/g,
              'Replacing references to images'
            ]
          ]
        }
      }
    },

    imagemin: {
      dist: {
        files: [
          {
            expand: true,
            cwd: '<%= yeoman.app %>/images',
            src: '{,*/}*.{png,jpg,jpeg,gif}',
            dest: '<%= yeoman.dist %>/images'
          }
        ]
      }
    },

    svgmin: {
      dist: {
        files: [
          {
            expand: true,
            cwd: '<%= yeoman.app %>/images',
            src: '{,*/}*.svg',
            dest: '<%= yeoman.dist %>/images'
          }
        ]
      }
    },

    // Copies remaining files to places other tasks can use
    copy: {
      dist: {
        files: [
          {
            expand: true,
            dot: true,
            cwd: '<%= yeoman.app %>',
            dest: '<%= yeoman.dist %>',
            src: [
              '*.{ico,png,txt}',
              '*.html',
              'images/{,*/}*.{webp}',
              'styles/fonts/{,*/}*.*'
            ]
          },
          {
            expand: true,
            cwd: '.tmp/images',
            dest: '<%= yeoman.dist %>/images',
            src: ['generated/*']
          },
          {
            expand: true,
            cwd: 'bower_components/bootstrap/dist',
            src: 'fonts/*',
            dest: '<%= yeoman.dist %>'
          }
        ]
      },
      styles: {
        expand: true,
        cwd: '<%= yeoman.app %>/styles',
        dest: '.tmp/styles/',
        src: '{,*/}*.css'
      }
    },

    // Run some tasks in parallel to speed up the build process
    concurrent: {
      server: ['copy:styles'],
      test: ['copy:styles'],
      dist: ['copy:styles', 'imagemin', 'svgmin']
    },

    // Test settings
    karma: {
      unit: {
        configFile: 'test/karma.conf.js',
        singleRun: true
      }
    }
  });

  grunt.registerTask('prettier', 'Run prettier on js.', function() {
    var done = this.async();
    var args = [
      '--single-quote',
      '--write',
      './Gruntfile.js',
      './' + appConfig.app + '/scripts/**/*.js'
    ];

    grunt.log.writeln('Running task with args: ' + args);
    var prettier = require('child_process').spawn(
      './node_modules/.bin/prettier',
      args
    );
    prettier.stdout.on('data', data => {
      grunt.log.writeln(`stdout: ${data}`);
    });

    prettier.stderr.on('data', data => {
      grunt.log.writeln(`stderr: ${data}`);
    });

    prettier.on('close', code => {
      grunt.log.writeln(`child process exited with code ${code}`);
      done(code === 0);
    });
  });

  grunt.registerTask(
    'setlocalconfig',
    'Sets a flag to enable config.local.js in bundled js',
    function() {
      ADD_CONFIG_LOCAL = true;
    }
  );

  grunt.registerTask(
    'local',
    'Watch local files and serve up development version without uglify',
    function() {
      var useminPrepare = grunt.config.get('useminPrepare');
      var jsSteps = useminPrepare.options.flow.html.steps.js;

      // remove uglifyjs when using `grunt watch`
      jsSteps.splice(jsSteps.indexOf('uglifyjs'), 1);
      grunt.config.set('useminPrepare', useminPrepare);

      var taskList = ['setlocalconfig', 'build', 'watch'];

      grunt.task.run(taskList);
    }
  );

  grunt.registerTask(
    'checklocalconfig',
    'Add config.local.js if enabled',
    function() {
      if (ADD_CONFIG_LOCAL) {
        var concat = grunt.config.get('concat');
        concat.generated.files[1].src.unshift(
          appConfig.app + '/scripts/config.local.js'
        );
        grunt.config.set('concat', concat);
      }
    }
  );

  grunt.registerTask('build', 'Build dist assets', function() {
    var taskList = [
      'clean:dist',
      'wiredep',
      'prettier',
      'useminPrepare',
      'checklocalconfig',
      'concurrent:dist',
      'less:dist',
      'postcss',
      'concat:generated',
      'copy:dist',
      'cssmin',
      'uglify',
      'filerev',
      'usemin'
    ];

    // don't run uglify if the usemin config doesn't say to
    var useminPrepare = grunt.config.get('useminPrepare');
    var jsSteps = useminPrepare.options.flow.html.steps.js;
    if (jsSteps.indexOf('uglifyjs') < 0) {
      taskList.splice(taskList.indexOf('uglify'), 1);
    }

    grunt.task.run(taskList);
  });

  grunt.registerTask('default', ['test', 'build']);
};
