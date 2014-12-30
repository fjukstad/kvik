module.exports = function (grunt) {
  'use strict';

  // check task runtime
  require('time-grunt')(grunt);

  // load generic configs
  var configs = require('dalek-build-tools');

  // project config
  grunt.initConfig({

    // load module meta data
    pkg: grunt.file.readJSON('package.json'),

    // define a src set of files for other tasks
    src: {
      lint: ['Gruntfile.js', 'index.js', 'lib/**/*.js', 'test/*.js'],
      complexity: ['index.js', 'lib/*.js', 'lib/comands/*.js'],
      test: ['test/*.js'],
      src: ['index.js', 'lib/**/*.js']
    },

    // clean automatically generated helper files & docs
    clean: configs.clean,

    // speed up build by defining concurrent tasks
    concurrent: configs.concurrent,

    // linting
    jshint: configs.jshint,

    // testing
    mochaTest: configs.mocha,

    // code metrics
    complexity: configs.complexity,
    plato: configs.plato(grunt.file.readJSON('.jshintrc')),

    // api docs
    yuidoc: configs.yuidocs(),

    // user docs
    documantix: {
      options: {
        header: 'dalekjs/dalekjs.com/master/assets/header.html',
        footer: 'dalekjs/dalekjs.com/master/assets/footer.html',
        target: 'report/docs',
        vars: {
          title: 'DalekJS - Documentation - Webdriver',
          desc: 'DalekJS - Documentation - Webdriver',
          docs: true
        }
      },
      src: ['index.js']
    },

    // add current timestamp to the html document
    includereplace: {
      dist: {
        options: {
          globals: {
            timestamp: '<%= grunt.template.today("dddd, mmmm dS, yyyy, h:MM:ss TT") %>'
          },
        },
        src: 'report/docs/*.html',
        dest: './'
      }
    },

    // up version, tag & commit
    bump: configs.bump({
      pushTo: 'git@github.com:dalekjs/dalek-internal-webdriver.git',
      files: ['package.json', 'CONTRIBUTORS.md', 'CHANGELOG.md']
    }),

    // generate contributors file
    contributors: configs.contributors,

    // compress artifacts
    compress: configs.compress,

    // prepare files for grunt-plato to
    // avoid error messages (weird issue...)
    preparePlato: {
      options: {
        folders: [
          'coverage',
          'report',
          'report/coverage',
          'report/complexity',
          'report/complexity/files',
          'report/complexity/files/index_js',
          'report/complexity/files/lib_driver_js',
          'report/complexity/files/lib_webdriver_js',
          'report/complexity/files/lib_commands_cookie_js',
          'report/complexity/files/lib_commands_element_js',
          'report/complexity/files/lib_commands_execute_js',
          'report/complexity/files/lib_commands_frame_js',
          'report/complexity/files/lib_commands_ime_js',
          'report/complexity/files/lib_commands_interaction_js',
          'report/complexity/files/lib_commands_page_js',
          'report/complexity/files/lib_commands_screenshot_js',
          'report/complexity/files/lib_commands_session_js',
          'report/complexity/files/lib_commands_storage_js',
          'report/complexity/files/lib_commands_timeout_js',
          'report/complexity/files/lib_commands_url_js',
          'report/complexity/files/lib_commands_window_js'
        ],
        files: [
          'report.history.json',
          'files/index_js/report.history.json',
          'files/lib_driver_js/report.history.json',
          'files/lib_webdriver_js/report.history.json',
          'files/lib_commands_cookie_js/report.history.json',
          'files/lib_commands_element_js/report.history.json',
          'files/lib_commands_execute_js/report.history.json',
          'files/lib_commands_frame_js/report.history.json',
          'files/lib_commands_ime_js/report.history.json',
          'files/lib_commands_interaction_js/report.history.json',
          'files/lib_commands_page_js/report.history.json',
          'files/lib_commands_screenshot_js/report.history.json',
          'files/lib_commands_session_js/report.history.json',
          'files/lib_commands_storage_js/report.history.json',
          'files/lib_commands_timeout_js/report.history.json',
          'files/lib_commands_url_js/report.history.json',
          'files/lib_commands_window_js/report.history.json'
        ]
      }
    },

    // prepare files & folders for coverage
    prepareCoverage: {
      options: {
        folders: ['coverage', 'report', 'report/coverage'],
        pattern: '[require("fs").realpathSync(__dirname + "/../index.js"), require("fs").realpathSync(__dirname + "/../lib/")]'
      }
    },

    // list requires that need to be changed
    // for generating a canary build
    'release-canary': {
      options: {
        files: ['index.js']
      }
    },

    // archives the docs if a new version appears
    archive: {
      options: {
        file: 'webdriver.html'
      }
    }

  });

  // load 3rd party tasks
  require('load-grunt-tasks')(grunt);
  grunt.loadTasks('./node_modules/dalek-build-tools/tasks');

  // define runner tasks
  grunt.registerTask('lint', 'jshint');
  
  // split test & docs for speed
  grunt.registerTask('test', ['clean:coverage', 'prepareCoverage', 'concurrent:test', 'generateCoverageBadge']);
  grunt.registerTask('docs', ['clean:reportZip', 'clean:report', 'preparePlato', 'documantix', 'includereplace', 'concurrent:docs', 'compress']);
  
  // release tasks
  grunt.registerTask('releasePatch', ['test', 'bump-only:patch', 'contributors', 'changelog', 'bump-commit']);
  grunt.registerTask('releaseMinor', ['test', 'bump-only:minor', 'contributors', 'changelog', 'bump-commit']);
  grunt.registerTask('releaseMajor', ['test', 'bump-only:major', 'contributors', 'changelog', 'bump-commit']);
  
  // clean, test, generate docs (the CI task)
  grunt.registerTask('all', ['clean', 'test', 'docs']);

};