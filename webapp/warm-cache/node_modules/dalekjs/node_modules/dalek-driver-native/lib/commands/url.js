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
 * Url related methods
 *
 * @module Driver
 * @class Url
 * @namespace Dalek.DriverNative.Commands
 */

var Url = {

  /**
   * Navigate to a new URL
   *
   * @method open
   * @param {string} url Url to navigate to
   * @param {string} hash Unique hash of that fn call
   * @param {string} uuid Unique hash of that fn call
   * @chainable
   */

  open: function (url, hash, uuid) {
    this.lastCalledUrl = url;
    this.actionQueue.push(this.webdriverClient.url.bind(this.webdriverClient, url));
    this.actionQueue.push(this._openCb.bind(this, url, hash, uuid));
    return this;
  },

  /**
   * Sends out an event with the results of the `open` call
   *
   * @method _openCb
   * @param {string} url Url to navigate to
   * @param {string} hash Unique hash of that fn call
   * @param {string} uuid Unique hash of that fn call
   * @return {object} promise Open promise
   * @private
   */

  _openCb: function (url, hash, uuid) {
    var deferred = Q.defer();
    this.events.emit('driver:message', {key: 'open', value: url, hash: hash, uuid: uuid});
    deferred.resolve();
    return deferred.promise;
  },

  /**
   * Fetches the current url
   *
   * @method url
   * @param {string} expected Expected url
   * @param {string} hash Unique hash of that fn call
   * @chainable
   */

  url: function (expected, hash) {
    this.actionQueue.push(this.webdriverClient.getUrl.bind(this.webdriverClient));
    this.actionQueue.push(this._urlCb.bind(this, expected, hash));
    return this;
  },

  /**
   * Sends out an event with the results of the `url` call
   *
   * @method _urlCb
   * @param {string} expected Expected url
   * @param {string} hash Unique hash of that fn call
   * @param {string} url Serialized JSON result of url call
   * @return {object} promise Url promise
   * @private
   */

  _urlCb: function (expected, hash, url) {
    var deferred = Q.defer();
    this.events.emit('driver:message', {key: 'url', expected: expected, hash: hash, value: JSON.parse(url).value});
    deferred.resolve();
    return deferred.promise;
  },

  /**
   * Navigate backwards in the browser history, if possible.
   *
   * @method back
   * @param {string} hash Unique hash of that fn call
   * @param {string} uuid Unique hash of that fn call
   * @chainable
   */

  back: function (hash, uuid) {
    this._createNonReturnee('back')(hash, uuid);
    return this;
  },

  /**
   * Navigate forwards in the browser history, if possible.
   *
   * @method forward
   * @param {string} hash Unique hash of that fn call
   * @param {string} uuid Unique hash of that fn call
   * @chainable
   */

  forward: function (hash, uuid) {
    this._createNonReturnee('forward')(hash, uuid);
    return this;
  },

  /**
   * Refresh the current page
   *
   * @method refresh
   * @param {string} hash Unique hash of that fn call
   * @param {string} uuid Unique hash of that fn call
   * @chainable
   */

  refresh: function (hash, uuid) {
    this._createNonReturnee('refresh')(hash, uuid);
    return this;
  }
};

/**
 * Mixes in url methods
 *
 * @param {Dalek.DriverNative} DalekNative Native driver base class
 * @return {Dalek.DriverNative} DalekNative Native driver base class
 */

module.exports = function (DalekNative) {
  // mixin methods
  Object.keys(Url).forEach(function (fn) {
    DalekNative.prototype[fn] = Url[fn];
  });

  return DalekNative;
};
