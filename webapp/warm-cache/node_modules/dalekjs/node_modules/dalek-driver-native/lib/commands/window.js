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

/**
 * Window related methods
 *
 * @module Driver
 * @class Window
 * @namespace Dalek.DriverNative.Commands
 */

var Window = {

  /**
   * Switches to another window context
   *
   * @method toFrame
   * @param {string} name Name of the window to switch to
   * @param {string} hash Unique hash of that fn call
   * @chainable
   */

  toWindow: function (name, hash) {
    this.actionQueue.push(this.webdriverClient.windowHandles.bind(this.webdriverClient));
    this.actionQueue.push(function (result) {
      var deferred = Q.defer();
      if (name === null) {
        deferred.resolve(JSON.parse(result).value[0]);
      }
      deferred.resolve(name);
      return deferred.promise;
    });
    this.actionQueue.push(this.webdriverClient.changeWindow.bind(this.webdriverClient));
    this.actionQueue.push(this._windowCb.bind(this, name, hash));
    return this;
  },

  /**
   * Sends out an event with the results of the `toWindow` call
   *
   * @method _windowCb
   * @param {string} name Name of the window to switch to
   * @param {string} hash Unique hash of that fn call
   * @param {string} result Serialized JSON with the reuslts of the toFrame call
   * @return {object} promise Exists promise
   * @private
   */

  _windowCb: function (name, hash) {
    var deferred = Q.defer();
    this.events.emit('driver:message', {key: 'toWindow', name: name, hash: hash, value: true});
    deferred.resolve();
    return deferred.promise;
  },

  /**
   * Resizes the current window
   *
   * @method resize
   * @param {object} dimensions New window width & height
   * @param {string} hash Unique hash of that fn call
   * @chainable
   */

  resize: function (dimensions, hash) {
    this.actionQueue.push(this.webdriverClient.setWindowSize.bind(this.webdriverClient, dimensions.width, dimensions.height));
    this.actionQueue.push(this._resizeCb.bind(this, dimensions, hash));
    return this;
  },

  /**
   * Sends out an event with the results of the `resize` call
   * and stores the current viewport dimensions to the
   * globale Config
   *
   * @method _windowCb
   * @param {object} dimensions New window width & height
   * @param {string} hash Unique hash of that fn call
   * @param {string} result Serialized JSON with the reuslts of the toFrame call
   * @return {object} promise Exists promise
   * @private
   */

  _resizeCb: function (dimensions, hash) {
    var deferred = Q.defer();
    this.config.config.viewport = dimensions;
    this.events.emit('driver:message', {key: 'resize', dimensions: dimensions, hash: hash, value: true});
    deferred.resolve();
    return deferred.promise;
  },

  /**
   * Maximizes the current window
   *
   * @method maximize
   * @param {string} hash Unique hash of that fn call
   * @chainable
   */

  maximize: function (hash) {
    this.actionQueue.push(this.webdriverClient.maximize.bind(this.webdriverClient));
    this.actionQueue.push(this._maximizeCb.bind(this, hash));
    return this;
  },

  /**
   * Sends out an event with the results of the `maximize` call
   *
   * @method _maximizeCb
   * @param {string} hash Unique hash of that fn call
   * @param {string} result Serialized JSON with the reuslts of the toFrame call
   * @return {object} promise Exists promise
   * @private
   */

  _maximizeCb: function (hash) {
    var deferred = Q.defer();
    this.events.emit('driver:message', {key: 'maximize', hash: hash, value: true});
    deferred.resolve();
    return deferred.promise;
  },

  /**
   * Closes the current window
   *
   * @method close
   * @param {string} hash Unique hash of that fn call
   * @chainable
   */

  close: function (hash) {
    this.actionQueue.push(this.webdriverClient.close.bind(this.webdriverClient));
    this.actionQueue.push(this._closeCb.bind(this, hash));
    return this;
  },

  /**
   * Sends out an event with the results of the `close` call
   *
   * @method _closeCb
   * @param {string} hash Unique hash of that fn call
   * @return {object} promise Exists promise
   * @private
   */

  _closeCb: function (hash) {
    var deferred = Q.defer();
    this.events.emit('driver:message', {key: 'close', hash: hash, value: true});
    deferred.resolve();
    return deferred.promise;
  }
};

/**
 * Mixes in element methods
 *
 * @param {Dalek.DriverNative} DalekNative Native driver base class
 * @return {Dalek.DriverNative} DalekNative Native driver base class
 */

module.exports = function (DalekNative) {
  // mixin methods
  Object.keys(Window).forEach(function (fn) {
    DalekNative.prototype[fn] = Window[fn];
  });

  return DalekNative;
};
