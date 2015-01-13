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
 * Url related Webdriver endpoints
 *
 * @param {Dalek.Internal.Webdriver} Driver Webdriver client instance
 * @return {Dalek.Internal.Webdriver} Driver Webdriver client instance
 */

module.exports = function (Driver) {
  'use strict';

  /**
   * Navigate to a new URL
   *
   * @method url
   * @see https://code.google.com/p/selenium/wiki/JsonWireProtocol#POST_/session/:sessionId/url
   * @param {GET} sessionId  ID of the session to route the command to
   * @param {POST} page The URL to navigate to.
   */

  Driver.addCommand({
    name: 'url',
    url: '/session/:sessionId/url',
    method: 'POST',
    params: ['page'],
    onRequest: function (params) {
      return {url: params.page};
    },
    onResponse: function (request, remote, options, deferred) {
      deferred.resolve(this);
    }
  });

  /**
   * Retrieve the URL of the current page
   *
   * @method getUrl
   * @see https://code.google.com/p/selenium/wiki/JsonWireProtocol#GET_/session/:sessionId/url
   * @param {GET} sessionId  ID of the session to route the command to
   */

  Driver.addCommand({
    name: 'getUrl',
    url: '/session/:sessionId/url',
    method: 'GET'
  });

  /**
   * Navigate forwards in the browser history, if possible.
   *
   * @method forward
   * @see https://code.google.com/p/selenium/wiki/JsonWireProtocol#/session/:sessionId/forward
   * @param {GET} sessionId  ID of the session to route the command to
   */

  Driver.addCommand({
    name: 'forward',
    url: '/session/:sessionId/forward',
    method: 'POST'
  });

  /**
   * Navigate backwards in the browser history, if possible
   *
   * @method back
   * @see https://code.google.com/p/selenium/wiki/JsonWireProtocol#/session/:sessionId/back
   * @param {GET} sessionId  ID of the session to route the command to
   */

  Driver.addCommand({
    name: 'back',
    url: '/session/:sessionId/back',
    method: 'POST'
  });

  /**
   * Refresh the current page
   *
   * @method refresh
   * @see https://code.google.com/p/selenium/wiki/JsonWireProtocol#/session/:sessionId/back
   * @param {GET} sessionId  ID of the session to route the command to
   */

  Driver.addCommand({
    name: 'refresh',
    url: '/session/:sessionId/refresh',
    method: 'GET'
  });
};
