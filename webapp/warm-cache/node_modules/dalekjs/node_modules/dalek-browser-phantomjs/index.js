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

// ext. libs
var Q = require('q');
var fs = require('fs');
var phantomjs = require('phantomjs');
var portscanner = require('portscanner');
var spawn = require('child_process').spawn;

/**
 * This module is a browser plugin for [DalekJS](//github.com/dalekjs/dalek).
 * It provides a browser launcher as well the PhantomJS browser itself.
 *
 * The browser plugin comes bundled with the DalekJS base framework.
 *
 * You can use the browser plugin beside others (it is the default)
 * by adding a config option to the your Dalekfile:
 *
 * ```javascript
 * "browser": ["phantomjs", "chrome"]
 * ```
 *
 * Or you can tell Dalek that it should test in this & another browser via the command line:
 *
 * ```bash
 * $ dalek mytest.js -b phantomjs,chrome
 * ```
 *
 * The Webdriver Server tries to open Port 9001 by default,
 * if this port is blocked, it tries to use a port between 9002 & 9091
 * You can specifiy a different port from within your [Dalekfile](/pages/config.html) like so:
 *
 * ```javascript
 * "browsers": {
 *   "phantomjs": {
 *     "port": 5555 
 *   }
 * }
 * ```
 *
 * It is also possible to specify a range of ports:
 *
 * ```javascript
 * "browsers": {
 *   "phantomjs": {
 *     "portRange": [6100, 6120] 
 *   }
 * }
 * ```
 *
 * If you would like to use a different Phantom version than the one that comes bundled with the driver,
 * your are able to specify its location in your [Dalekfile](/pages/config.html):
 *
 * ```javascript
 * "browsers": {
 *   "phantomjs": {
 *     "binary": "~/bin/phantomjs" 
 *   }
 * }
 * ```
 *
 * If you would like to preserve the ability to use the bundled version,
 * you can also add an additional browser launcher in your [Dalekfile](/pages/config.html).
 *
 * ```javascript
 * "browsers": {
 *   "phantomjs:1.9.1": {
 *     "binary": "~/bin/phantomjs" 
 *   }
 * }
 * ```
 *
 * And then launch it like this:
 * 
 * ```bash
 * $ dalek mytest.js -b phantomjs:1.9.1
 * ```
 * 
 * @module DalekJS
 * @class PhantomJSDriver
 * @namespace Browser
 * @part PhantomJS
 * @api
 */

var PhantomJSDriver = {

  /**
   * Verbose version of the browser name
   *
   * @property
   * @type string
   * @default PhantomJS
   */

  longName: 'PhantomJS',

  /**
   * Default port of the PhantomJSDriver
   * The port may change, cause the port conflict resultion
   * tool might pick another one, if the default one is blocked
   *
   * @property
   * @type integer
   * @default 9001
   */

  port: 9001,

  /**
   * Default maximum port of the Ghostdriver Server
   * The port is the highest port in the range that can be allocated
   * by the Ghostdriver Server
   *
   * @property maxPort
   * @type integer
   * @default 9091
   */

  maxPort: 9091,

  /**
   * Default host of the PhantomJSDriver
   * The host may be overriden with
   * a user configured value
   *
   * @property
   * @type string
   * @default localhost
   */

  host: 'localhost',

  /**
   * Root path of the PhantomJSDriver
   *
   * @property
   * @type string
   * @default /wd/hub
   */

  path: '/wd/hub',

  /**
   * Default desired capabilities that should be
   * transferred when the browser session gets requested
   *
   * @property desiredCapabilities
   * @type object
   */

  desiredCapabilities: {
    version: phantomjs.version,
    browserName: 'phantomjs'
  },

  /**
   * Driver defaults, what should the driver be able to access.
   *
   * @property driverDefaults
   * @type object
   */

  driverDefaults: {
    viewport: true,
    status: true,
    sessionInfo: true
  },

  /**
   * Child process instance of the PhantomJS browser
   *
   * @property
   * @type null|Object
   */

  spawned: null,

  /**
   * Resolves the driver port
   *
   * @method getPort
   * @return integer
   */

  getPort: function () {
    return this.port;
  },

  /**
   * Resolves the maximum range for the driver port
   *
   * @method getMaxPort
   * @return {integer} port Max WebDriver server port range
   */

  getMaxPort: function () {
    return this.maxPort;
  },

  /**
   * Returns the driver host
   *
   * @method getHost
   * @type string
   */

  getHost: function () {
    return this.host;
  },

  /**
   * Launches PhantomJS, negoatiates a port & checks for a user set binary
   *
   * @method launch
   * @param {object} configuration Browser configuration
   * @param {EventEmitter2} events EventEmitter (Reporter Emitter instance)
   * @param {Dalek.Internal.Config} config Dalek configuration class
   * @return {object} promise Browser promise
   */

  launch: function (configuration, events, config) {
    var deferred = Q.defer();

    // store injected configuration/log event handlers
    this.reporterEvents = events;
    this.configuration = configuration;
    this.config = config;

    // check for a user set port
    var browsers = this.config.get('browsers');
    if (browsers && Array.isArray(browsers)) {
      browsers.forEach(this._checkUserDefinedPorts.bind(this));
    }

    // check if the current port is in use, if so, scan for free ports
    portscanner.findAPortNotInUse(this.getPort(), this.getMaxPort(), this.getHost(), this._checkPorts.bind(this, deferred));
    return deferred.promise;
  },

  /**
   * Kills the PhantomJSDriver processe
   *
   * @method kill
   * @chainable
   */

  kill: function () {
    this.spawned.kill('SIGTERM');
    return this;
  },

  /**
   * Checks if the def. port is blocked & if we need to switch to another port
   * Kicks off the process manager (for closing the opened browsers after the run has been finished)
   * Also starts the chromedriver instance 
   *
   * @method _checkPorts
   * @param {object} deferred Promise
   * @param {null|object} error Error object
   * @param {integer} port Found open port
   * @private
   * @chainable
   */

  _checkPorts: function (deferred, error, port) {
    // check if the port was blocked & if we need to switch to another port
    if (this.port !== port) {
      this.reporterEvents.emit('report:log:system', 'dalek-browser-phantomjs: Switching to port: ' + port);
      this.port = port;
    }

    // check the binary
    var binary = this._checkUserDefinedBinary(this.configuration, phantomjs.path);

    // launch the browser process
    this.spawned = spawn(binary, ['--webdriver', this.getPort(), '--ignore-ssl-errors=true']);
    this.spawned.stdout.on('data', this._launch.bind(this, deferred));
    return this;
  },

  /**
   * Checks the data stream from the launched phantom process
   *
   * @method _launch
   * @param {object} deferred Promise
   * @param {buffer} data Console output from Ghostdriver
   * @chainable
   * @private
   */

  _launch: function (deferred, data) {
    var stream = data + '';
    
    // check if ghostdriver could be launched    
    if (stream.search('GhostDriver - Main - running') !== -1) {
      deferred.resolve();
    } else if (stream.search('Could not start Ghost Driver') !== -1) {
      this.reporterEvents.emit('error', 'dalek-browser-phantomjs: Could not start Ghost Driver');
      deferred.reject('Could not start Ghost Driver');
      process.exit(127);
    }

    return this;
  },

  /**
   * Process user defined ports
   *
   * @method _checkUserDefinedPorts
   * @param {object} browser Browser configuration
   * @chainable
   * @private
   */

  _checkUserDefinedPorts: function (browser) {
    // check for a single defined port
    if (browser.phantomjs && browser.phantomjs.port) {
      this.port = parseInt(browser.phantomjs.port, 10);
      this.maxPort = this.port + 90;
      this.reporterEvents.emit('report:log:system', 'dalek-browser-phantomjs: Switching to user defined port: ' + this.port);
    }

    // check for a port range
    if (browser.phantomjs && browser.phantomjs.portRange && browser.phantomjs.portRange.length === 2) {
      this.port = parseInt(browser.phantomjs.portRange[0], 10);
      this.maxPort = parseInt(browser.phantomjs.portRange[1], 10);
      this.reporterEvents.emit('report:log:system', 'dalek-browser-phantomjs: Switching to user defined port(s): ' + this.port + ' -> ' + this.maxPort);
    }

    return this;
  },

  /**
   * Checks if the binary exists,
   * when set manually by the user
   *
   * @method _checkUserDefinedBinary
   * @param {string} binary Path to the browser binary
   * @return {bool|string} Binary path if binary exists, else false
   * @private
   */

  _checkUserDefinedBinary: function (configuration, defaultBinary) {
    var binary = defaultBinary;

    // check if we have a user defined binary
    if (configuration && configuration.binary) {
      binary = configuration.binary;
    }

    // check if we need to replace the users home directory
    if (process.platform === 'darwin' && binary.trim()[0] === '~') {
      binary = binary.replace('~', process.env.HOME);
    }

    // check if the binary exists
    if (!fs.existsSync(binary)) {
      this.reporterEvents.emit('error', 'dalek-driver-phantomjs: Binary not found: ' + binary);
      process.exit(127);
      return false;
    }

    return binary;
  }

};

module.exports = PhantomJSDriver;
