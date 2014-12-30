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
 * Page related methods
 *
 * @module Driver
 * @class Page
 * @namespace Dalek.DriverNative.Commands
 */

var Page = {

  /**
   * This function is non operational
   *
   * @method noop
   * @param {mixed} message Whatever yu like
   * @param {string} hash Unique hash of that fn call
   * @chainable
   */

  noop: function (message, hash) {
    this.actionQueue.push(this._noopCb.bind(this, message, hash));
    return this;
  },

  /**
   * Sends out an event with the results of the `noop` call
   *
   * @method _noopCb
   * @param {mixed} message Whatever yu like
   * @param {string} hash Unique hash of that fn call
   * @return {object} Promise
   * @private
   */

  _noopCb: function (message, hash) {
    var deferred = Q.defer();
    this.events.emit('driver:message', {key: 'noop', uuid: hash, hash: hash, value: message});
    deferred.resolve();
    return deferred.promise;
  },

  /**
   * Gets the HTML source of a page
   *
   * @method source
   * @param {string} hash Unique hash of that fn call
   * @chainable
   */

  source: function (hash) {
    this.actionQueue.push(this.webdriverClient.source.bind(this.webdriverClient));
    this.actionQueue.push(this._sourceCb.bind(this, hash));
    return this;
  },

  /**
   * Sends out an event with the results of the `source` call
   *
   * @method _sourceCb
   * @param {string} hash Unique hash of that fn call
   * @param {string} source Serialized JSON with the results of the source call
   * @return {object} Promise
   * @private
   */

  _sourceCb: function (hash, source) {
    var deferred = Q.defer();
    this.events.emit('driver:message', {key: 'source', uuid: hash, hash: hash, value: JSON.parse(source).value});
    deferred.resolve();
    return deferred.promise;
  },

  /**
   * Checks the document title of a page
   *
   * @method title
   * @param {string} expected Expected page title
   * @param {string} hash Unique hash of that fn call
   * @chainable
   */

  title: function (expected, hash) {
    this.actionQueue.push(this.webdriverClient.title.bind(this.webdriverClient));
    this.actionQueue.push(this._titleCb.bind(this, expected, hash));
    return this;
  },

  /**
   * Sends out an event with the results of the `title` call
   *
   * @method _titleCb
   * @param {string} expected Expected page title
   * @param {string} hash Unique hash of that fn call
   * @param {string} title Serialized JSON with the results of the title call
   * @return {object} promise Title promise
   * @private
   */

  _titleCb: function (expected, hash, title) {
    var deferred = Q.defer();
    this.events.emit('driver:message', {key: 'title', expected: expected, hash: hash, value: JSON.parse(title).value});
    deferred.resolve();
    return deferred.promise;
  },

  /**
   * Checks the text of an alaert, prompt or confirm dialog
   *
   * @method alertText
   * @param {string} expected Expected alert text
   * @param {string} hash Unique hash of that fn call
   * @chainable
   */

  alertText: function (expected, hash) {
    this.actionQueue.push(this.webdriverClient.alertText.bind(this.webdriverClient));
    this.actionQueue.push(this._alertTextCb.bind(this, expected, hash));
    return this;
  },

  /**
   * Sends out an event with the results of the `alertText` call
   *
   * @method _alertTextCb
   * @param {string} expected Expected alert text
   * @param {string} hash Unique hash of that fn call
   * @param {string} alertText Serialized JSON with the results of the alertText call
   * @return {object} promise alertText promise
   * @private
   */

  _alertTextCb: function (expected, hash, alertText) {
    var deferred = Q.defer();
    this.events.emit('driver:message', {key: 'alertText', expected: expected, hash: hash, value: JSON.parse(alertText).value});
    deferred.resolve();
    return deferred.promise;
  },

  /**
   * Sends text to a javascript prompt dialog box
   *
   * @method promptText
   * @param {object} dimensions New window width & height
   * @param {string} hash Unique hash of that fn call
   * @chainable
   */

  promptText: function (text, hash) {
    this.actionQueue.push(this.webdriverClient.promptText.bind(this.webdriverClient, text));
    this.actionQueue.push(this._promptTextCb.bind(this, text, hash));
    return this;
  },

  /**
   * Sends out an event with the results of the `promptText` call
   *
   * @method _promptTextCb
   * @param {object} dimensions New window width & height
   * @param {string} hash Unique hash of that fn call
   * @param {string} result Serialized JSON with the reuslts of the toFrame call
   * @return {object} promise Exists promise
   * @private
   */

  _promptTextCb: function (text, hash) {
    var deferred = Q.defer();
    this.events.emit('driver:message', {key: 'promptText', text: text, hash: hash, value: true});
    deferred.resolve();
    return deferred.promise;
  },

  /**
   * Accepts (e.g. oressing the OK button) an javascript alert, prompt or confirm dialog
   *
   * @method acceptAlert
   * @param {string} hash Unique hash of that fn call
   * @chainable
   */

  acceptAlert: function (hash) {
    this.actionQueue.push(this.webdriverClient.acceptAlert.bind(this.webdriverClient));
    this.actionQueue.push(this._acceptAlertCb.bind(this, hash));
    return this;
  },

  /**
   * Sends out an event with the results of the `acceptAlert` call
   *
   * @method _acceptAlertCb
   * @param {string} hash Unique hash of that fn call
   * @param {string} result Serialized JSON with the reuslts of the toFrame call
   * @return {object} promise Exists promise
   * @private
   */

  _acceptAlertCb: function (text, hash) {
    var deferred = Q.defer();
    this.events.emit('driver:message', {key: 'acceptAlert', hash: hash, value: true});
    deferred.resolve();
    return deferred.promise;
  },

  /**
   * Accepts (e.g. oressing the OK button) an javascript alert, prompt or confirm dialog
   *
   * @method dismissAlert
   * @param {string} hash Unique hash of that fn call
   * @chainable
   */

  dismissAlert: function (hash) {
    this.actionQueue.push(this.webdriverClient.dismissAlert.bind(this.webdriverClient));
    this.actionQueue.push(this._dismissAlertCb.bind(this, hash));
    return this;
  },

  /**
   * Sends out an event with the results of the `dismissAlert` call
   *
   * @method _dismissAlertCb
   * @param {string} hash Unique hash of that fn call
   * @param {string} result Serialized JSON with the reuslts of the toFrame call
   * @return {object} promise Exists promise
   * @private
   */

  _dismissAlertCb: function (text, hash) {
    var deferred = Q.defer();
    this.events.emit('driver:message', {key: 'dismissAlert', hash: hash, value: true});
    deferred.resolve();
    return deferred.promise;
  },

  /**
   * Wait for a specific amount of time
   *
   * @method wait
   * @param {integer} timeout Time to wait in ms
   * @param {string} hash Unique hash of that fn call
   * @param {string} uuid Unique hash of that fn call
   * @chainable
   */

  wait: function (timeout, hash, uuid) {
    this.actionQueue.push(this.webdriverClient.implicitWait.bind(this.webdriverClient, timeout));
    this.actionQueue.push(this._waitCb.bind(this, timeout, hash, uuid));
    return this;
  },

  /**
   * Sends out an event with the results of the `wait` call
   *
   * @method _waitCb
   * @param {integer} timeout Time to wait in ms
   * @param {string} hash Unique hash of that fn call
   * @param {string} uuid Unique hash of that fn call
   * @return {object} promise WaitForElement promise
   * @private
   */

  _waitCb: function (timeout, hash, uuid) {
    var deferred = Q.defer();
    this.events.emit('driver:message', {key: 'wait', timeout: timeout, uuid: uuid, hash: hash, value: timeout + ' ms'});
    setTimeout(function () {
      deferred.resolve();
    }.bind(this), timeout);
    return deferred.promise;
  }

};

/**
 * Mixes in page methods
 *
 * @param {Dalek.DriverNative} DalekNative Native driver base class
 * @return {Dalek.DriverNative} DalekNative Native driver base class
 */

module.exports = function (DalekNative) {
  // mixin methods
  Object.keys(Page).forEach(function (fn) {
    DalekNative.prototype[fn] = Page[fn];
  });

  return DalekNative;
};
