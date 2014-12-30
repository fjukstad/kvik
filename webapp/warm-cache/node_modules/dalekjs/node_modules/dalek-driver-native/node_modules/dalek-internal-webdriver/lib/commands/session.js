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
 * Browser session related Webdriver endpoints
 *
 * @param {Dalek.Internal.Webdriver} Driver Webdriver client instance
 * @return {Dalek.Internal.Webdriver} Driver Webdriver client instance
 */

module.exports = function (Driver) {
  'use strict';

  /**
   * Create a new session. The server should attempt to create a session that most closely matches
   * the desired and required capabilities. Required capabilities have higher priority than
   * desired capabilities and must be set for the session to be created.
   *
   * @method createSession
   * @see https://code.google.com/p/selenium/wiki/JsonWireProtocol#POST_/session
   * @param {POST} desiredCapabilities An object describing the session's desired capabilities.
   * @param {POST} requiredCapabilities An object describing the session's required capabilities
   */

  Driver.addCommand({
    name: 'createSession',
    url: '/session',
    method: 'POST',
    params: ['desiredCapabilities'],
    onRequest: function (params) {
      return params.desiredCapabilities;
    },
    onResponse: function (request, remote, options, deferred, rawData) {
      if (!request.headers.location) {
        this.options.sessionId = JSON.parse(rawData).sessionId;
      } else {
        this.options.sessionId = request.headers.location.replace('http://' + options.hostname + ':' + options.port + options.path + '/', '').replace('/session/', '').replace('/wd/hub', '');
        
        // fix for sauce (no port in url)
        if (this.options.sessionId.search('http://') !== -1) {
          this.options.sessionId = request.headers.location.replace('http://' + options.hostname + options.path + '/', '').replace('/session/', '').replace('/wd/hub', '');
        }
      }

      deferred.resolve(this);
    }
  });

  /**
   * Returns a list of the currently active sessions
   *
   * @method sessions
   * @see https://code.google.com/p/selenium/wiki/JsonWireProtocol#GET_/sessions
   */

  Driver.addCommand({
    name: 'sessions',
    url: '/sessions',
    method: 'GET'
  });

  /**
   * Retrieve the capabilities of the specified session.
   *
   * @method session
   * @see https://code.google.com/p/selenium/wiki/JsonWireProtocol#GET_/session/:sessionId
   * @param {GET} sessionId ID of the session to route the command to
   */

  Driver.addCommand({
    name: 'sessionInfo',
    url: '/session/:sessionId',
    method: 'GET'
  });

  /**
   * Delete the session.
   *
   * @method session
   * @see https://code.google.com/p/selenium/wiki/JsonWireProtocol#GET_/session/:sessionId
   * @param {GET} sessionId ID of the session to route the command to
   */

  Driver.addCommand({
    name: 'deleteSession',
    url: '/session/:sessionId',
    method: 'DELETE'
  });

};
