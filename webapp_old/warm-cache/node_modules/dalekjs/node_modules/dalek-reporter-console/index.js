/*!
 *
 * Copyright (c) 2013 Sebastian Golasch
 *
 * Permission is hereby granted, free of charge, to any person obtaining a
 * copy of this software and associated documentation files (the "Software"),
 * to deal in the Software without restriction, including without limitation
 * the rights to use, copy, modify, merge, publish, distribute, sublicense,
 * and/or sell copies of the Software, and to permit persons to whom the
 * Software is furnished to do so, subject to the following conditions:
 *
 * The above copyright notice and this permission notice shall be included
 * in all copies or substantial portions of the Software.
 *
 * THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS
 * OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
 * FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL
 * THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
 * LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING
 * FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER
 * DEALINGS IN THE SOFTWARE.
 *
 */

'use strict';

// int. globals
var reporter = null;

/**
 * Daleks basic reporter, all the lovely colors & symbols you see when running dalek.
 * The reporter will be installed by default.
 *
 * If you would like to use the reporter in addition to another one,
 * you can start dalek with a special command line argument
 *
 * ```bash
 * $ dalek your_test.js -r console,junit
 * ```
 *
 * or you can add it to your Dalekfile
 *
 * ```javascript
 * "reporter": ["console", "junit"]
 * ```
 *
 * @class Reporter
 * @constructor
 * @part Console
 * @api
 */

var Reporter = function (opts) {
  var loglevel = opts && opts.logLevel ? parseInt(opts.logLevel, 10) : 1;
  this.level = (loglevel >= -1 && loglevel <= 5) ? loglevel : 1;
  this.events = opts.events;

  // set color & symbols flags
  this.noColor = opts.config.config.noColors;
  this.noSymbols = opts.config.config.noSymbols;

  this.importLogModule();
  this.startListening();
};

/**
 * @module Reporter
 */

module.exports = function (opts) {
  if (reporter === null) {
    reporter = new Reporter(opts);
  }

  return reporter;
};

Reporter.prototype = {

  /**
   * Imports an output module with the correct log state
   *
   * @method importLogModule
   * @param {object} data
   * @chainable
   */

  importLogModule: function () {
    var logModule = require('./lib/levelbase');
    if (this.level !== -1) {
      logModule = require('./lib/loglevel/level' + this.level);
    }

    var methods = Object.keys(logModule.prototype);

    methods.forEach(function (method) {
      this[method] = logModule.prototype[method];
    }.bind(this));
    return this;
  },

  /**
   * Connects to all the event listeners
   *
   * @method startListening
   * @param {object} data
   * @chainable
   */

  startListening: function () {
    // assertion & action status
    this.events.on('report:assertion', this.outputAssertionResult.bind(this));
    this.events.on('report:assertion:status', this.outputAssertionExpectation.bind(this));
    this.events.on('report:action', this.outputAction.bind(this));

    // test status
    this.events.on('report:test:finished', this.outputTestFinished.bind(this));
    this.events.on('report:test:started', this.outputTestStarted.bind(this));

    // runner status
    this.events.on('report:runner:started', this.outputRunnerStarted.bind(this));
    this.events.on('report:runner:finished', this.outputRunnerFinished.bind(this));

    // session & browser status
    this.events.on('report:run:browser', this.outputRunBrowser.bind(this));
    this.events.on('report:driver:status', this.outputOSVersion.bind(this));
    this.events.on('report:driver:session', this.outputBrowserVersion.bind(this));

    // remote connections
    this.events.on('report:remote:ready', this.remoteConnectionReady.bind(this));
    this.events.on('report:remote:established', this.remoteConnectionEstablished.bind(this));
    this.events.on('report:remote:closed', this.remoteConnectionClosed.bind(this));

    // logs
    this.events.on('report:log:system', this.outputLogUser.bind(this, 'system'));
    this.events.on('report:log:driver', this.outputLogUser.bind(this, 'driver'));
    this.events.on('report:log:browser', this.outputLogUser.bind(this, 'browser'));
    this.events.on('report:log:user', this.outputLogUser.bind(this, 'user'));
    this.events.on('report:log:system:webdriver', this.outputLogUser.bind(this, 'webdriver'));

    // errors & warnings
    this.events.on('error', this.outputError.bind(this));
    this.events.on('warning', this.outputWarning.bind(this));

    // reports
    this.events.on('report:written', this.outputReportWritten.bind(this));

    return this;
  }
};
