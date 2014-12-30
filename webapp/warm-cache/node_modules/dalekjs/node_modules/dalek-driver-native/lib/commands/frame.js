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
 * Frame related methods
 *
 * @module Driver
 * @class Frame
 * @namespace Dalek.DriverNative.Commands
 */

var Frame = {

  /**
   * Switches to frame context
   *
   * @method toFrame
   * @param {string} selector Selector expression to find the element
   * @param {string} hash Unique hash of that fn call
   * @chainable
   */

  toFrame: function (selector, hash) {
    if (selector !== null) {
      this.actionQueue.push(this.webdriverClient.element.bind(this.webdriverClient, selector));
    }
    this.actionQueue.push(this.webdriverClient.frame.bind(this.webdriverClient));
    this.actionQueue.push(this._frameCb.bind(this, selector, hash));
    return this;
  },

  /**
   * Sends out an event with the results of the `toFrame` call
   *
   * @method _frameCb
   * @param {string} selector Selector expression to find the element
   * @param {string} hash Unique hash of that fn call
   * @param {string} result Serialized JSON with the reuslts of the toFrame call
   * @return {object} promise Exists promise
   * @private
   */

  _frameCb: function (selector, hash) {
    var deferred = Q.defer();
    this.events.emit('driver:message', {key: 'toFrame', selector: selector, hash: hash, value: true});
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
  Object.keys(Frame).forEach(function (fn) {
    DalekNative.prototype[fn] = Frame[fn];
  });

  return DalekNative;
};
