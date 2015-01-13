module.exports = function(grunt) {
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
      complexity: ['index.js', 'lib/**/*.js'],
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
          title: 'DalekJS - Documentation - Driver',
          desc: 'DalekJS - Documentation - Driver',
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
      pushTo: 'git@github.com:dalekjs/dalek-driver-native.git',
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
          'report/complexity',
          'report/complexity/files',
          'report/complexity/files/test',
          'report/complexity/files/index_js',
          'report/complexity/files/lib_commands_cookie_js',
          'report/complexity/files/lib_commands_element_js',
          'report/complexity/files/lib_commands_execute_js',
          'report/complexity/files/lib_commands_frame_js',
          'report/complexity/files/lib_commands_page_js',
          'report/complexity/files/lib_commands_screenshot_js',
          'report/complexity/files/lib_commands_window_js',
          'report/complexity/files/lib_commands_url_js'
        ],
        files: [
          'report.history.json',
          'files/report.history.json',
          'files/index_js/report.history.json',
          'files/lib_commands_cookie_js/report.history.json',
          'files/lib_commands_element_js/report.history.json',
          'files/lib_commands_execute_js/report.history.json',
          'files/lib_commands_frame_js/report.history.json',
          'files/lib_commands_page_js/report.history.json',
          'files/lib_commands_screenshot_js/report.history.json',
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

    // archive docs
    archive: {
      options: {
        files: ['drivernative.html']
      }
    },

    // release canary version
    'release-canary': {
      options: {
        files: ['index.js']
      }
    },

    watch: {
      scripts: {
        files: ['lib/**/*.js', 'test/**/*.js'],
        tasks: ['lint', 'mochaTest'],
        options: {
          atBegin: true
        }
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
  grunt.registerTask('docs', ['clean:reportZip', 'clean:report', 'preparePlato', 'concurrent:docs', 'documantix', 'includereplace', 'compress']);
  
  // release tasks
  grunt.registerTask('releasePatch', ['test', 'bump-only:patch', 'contributors', 'changelog', 'bump-commit']);
  grunt.registerTask('releaseMinor', ['test', 'bump-only:minor', 'contributors', 'changelog', 'bump-commit']);
  grunt.registerTask('releaseMajor', ['test', 'bump-only:major', 'contributors', 'changelog', 'bump-commit']);
  
  // clean, test, generate docs (the CI task)
  grunt.registerTask('all', ['clean', 'test', 'docs']);

};
