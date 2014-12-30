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
 * Cookie related methods
 *
 * @module Driver
 * @class Cookie
 * @namespace Dalek.DriverNative.Commands
 */

var Cookie = {

  /**
   * Sets an cookie
   *
   * @method setCookie
   * @param {string} name Name of the cookie
   * @param {string} contents Contents of the cookie
   * @param {string} hash Unique hash of that fn call
   * @chainable
   */

  setCookie: function (name, contents, hash) {
    this.actionQueue.push(this.webdriverClient.setCookie.bind(this.webdriverClient, {name: name, value: contents}));
    this.actionQueue.push(this._setCookieCb.bind(this, name, contents, hash));
    return this;
  },

  /**
   * Sends out an event with the results of the `setCookie` call
   *
   * @method _setCookieCb
   * @param {string} name Name of the cookie
   * @param {string} contents Contents of the cookie
   * @param {string} hash Unique hash of that fn call
   * @return {object} promise Exists promise
   * @private
   */

  _setCookieCb: function (name, contents, hash) {
    var deferred = Q.defer();
    this.events.emit('driver:message', {key: 'setCookie', value: name, contents: contents, uuid: hash, hash: hash});
    deferred.resolve();
    return deferred.promise;
  },

  /**
   * Retrieves an cookie
   *
   * @method cookie
   * @param {string} name Name of the cookie
   * @param {string} cookie Expected contents of the cookie
   * @param {string} hash Unique hash of that fn call
   * @chainable
   */

  cookie: function (name, cookie, hash) {
    this.actionQueue.push(this.webdriverClient.getCookie.bind(this.webdriverClient, name));
    this.actionQueue.push(this._cookieCb.bind(this, name, hash));
    return this;
  },

  /**
   * Sends out an event with the results of the `setCookie` call
   *
   * @method _cookieCb
   * @param {string} name Name of the cookie
   * @param {string} expected Expected contents of the cookie
   * @param {string} hash Unique hash of that fn call
   * @param {string} result Serialized JSON with the reuslts of the exists call
   * @return {object} promise Exists promise
   * @private
   */

  _cookieCb: function (name, expected, hash, res) {
    var deferred = Q.defer();
    var cookies = JSON.parse(res).value;
    var cookie = false;
    cookies.forEach(function (cookieItem) {
      if (cookieItem.name === name) {
        cookie = cookieItem.value;
      }
    });

    this.events.emit('driver:message', {key: 'cookie', value: cookie, name: name, expected: expected, uuid: hash, hash: hash});
    deferred.resolve();
    return deferred.promise;
  },
};

/**
 * Mixes in cookie methods
 *
 * @param {Dalek.DriverNative} DalekNative Native driver base class
 * @return {Dalek.DriverNative} DalekNative Native driver base class
 */

module.exports = function (DalekNative) {
  // mixin methods
  Object.keys(Cookie).forEach(function (fn) {
    DalekNative.prototype[fn] = Cookie[fn];
  });

  return DalekNative;
};
