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

/**
 * Storage related Webdriver endpoints
 *
 * @param {Dalek.Internal.Webdriver} Driver Webdriver client instance
 * @return {Dalek.Internal.Webdriver} Driver Webdriver client instance
 */

module.exports = function (Driver) {
  'use strict';

  /**
   * Get all keys of the browsers local storage
   *
   * @method getLocalStorage
   * @see https://code.google.com/p/selenium/wiki/JsonWireProtocol#GET_/session/:sessionId/local_storage
   * @param {GET} sessionId ID of the session to route the command to
   */

  Driver.addCommand({
    name: 'getLocalStorage',
    url: '/session/:sessionId/local_storage',
    method: 'GET'
  });

  /**
   * Set the storage item for the given key
   *
   * @method setLocalStorage
   * @see https://code.google.com/p/selenium/wiki/JsonWireProtocol#POST_/session/:sessionId/local_storage
   * @param {GET} sessionId ID of the session to route the command to
   * @param {POST} key The key to set
   * @param {POST} value The value to set
   */

  Driver.addCommand({
    name: 'setLocalStorage',
    url: '/session/:sessionId/local_storage/key/:key',
    method: 'POST',
    params: ['key', 'value']
  });

  /**
   * Get the number of items in the storage
   *
   * @method getLocalStorageSize
   * @see https://code.google.com/p/selenium/wiki/JsonWireProtocol#GET_/session/:sessionId/local_storage/size
   * @param {GET} sessionId ID of the session to route the command to
   */

  Driver.addCommand({
    name: 'getLocalStorageSize',
    url: '/session/:sessionId/local_storage/size',
    method: 'GET'
  });

  /**
   * Get all keys of the browsers session storage
   *
   * @method getSessionStorage
   * @see https://code.google.com/p/selenium/wiki/JsonWireProtocol#GET_/session/:sessionId/session_storage
   * @param {GET} sessionId ID of the session to route the command to
   */

  Driver.addCommand({
    name: 'getSessionStorage',
    url: '/session/:sessionId/session_storage',
    method: 'GET'
  });

  /**
   * Get the storage item for the given key
   *
   * @method getSessionStorageEntry
   * @see https://code.google.com/p/selenium/wiki/JsonWireProtocol#GET_/session/:sessionId/session_storage/key/:key
   * @param {GET} sessionId ID of the session to route the command to
   * @param {GET} key  The key to get
   */

  Driver.addCommand({
    name: 'getSessionStorageEntry',
    url: '/session/:sessionId/session_storage/key/:key',
    method: 'GET'
  });

  /**
   * Get the number of items in the storage
   *
   * @method getSessionStorageSize
   * @see https://code.google.com/p/selenium/wiki/JsonWireProtocol#GET_/session/:sessionId/session_storage/size
   * @param {GET} sessionId ID of the session to route the command to
   */

  Driver.addCommand({
    name: 'getSessionStorageSize',
    url: '/session/:sessionId/session_storage/size',
    method: 'GET'
  });

  /**
   * Get the status of the html5 application cache.
   *
   * @method applicationCacheStatus
   * @see https://code.google.com/p/selenium/wiki/JsonWireProtocol#GET_/session/:sessionId/application_cache/status
   * @param {GET} sessionId ID of the session to route the command to
   */

  Driver.addCommand({
    name: 'applicationCacheStatus',
    url: '/session/:sessionId/application_cache/status',
    method: 'GET'
  });

  return Driver;
};
