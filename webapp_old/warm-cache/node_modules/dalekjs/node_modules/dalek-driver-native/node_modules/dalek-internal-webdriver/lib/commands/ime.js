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
 * IME related Webdriver endpoints
 *
 * @param {Dalek.Internal.Webdriver} Driver Webdriver client instance
 * @return {Dalek.Internal.Webdriver} Driver Webdriver client instance
 */

module.exports = function (Driver) {
  'use strict';

  /**
   * Query the server's current status.
   * The server should respond with a general "HTTP 200 OK" response if it is alive and accepting commands.
   * The response body should be a JSON object describing the state of the server.
   * All server implementations should return two basic objects describing the server's current platform
   * and when the server was built. All fields are optional; if omitted, the client should assume the value is uknown.
   * Furthermore, server implementations may include additional fields not listed here.
   *
   * @method status
   * @see https://code.google.com/p/selenium/wiki/JsonWireProtocol#GET_/status
   */

  Driver.addCommand({
    name: 'status',
    url: '/status',
    method: 'GET'
  });

  /**
   * List all available engines on the machine.
   * To use an engine, it has to be present in this list.
   *
   * @method availableEngines
   * @see https://code.google.com/p/selenium/wiki/JsonWireProtocol#GET_/session/:sessionId/ime/available_engines
   * @param {GET} sessionId ID of the session to route the command to
   */

  Driver.addCommand({
    name: 'availableEngines',
    url: '/session/:sessionId/ime/available_engines',
    method: 'GET'
  });

  /**
   * Get the name of the active IME engine.
   * The name string is platform specific.
   *
   * @method activeEngine
   * @see https://code.google.com/p/selenium/wiki/JsonWireProtocol#GET_/session/:sessionId/ime/active_engine
   * @param {GET} sessionId ID of the session to route the command to
   */

  Driver.addCommand({
    name: 'activeEngine',
    url: '/session/:sessionId/ime/active_engine',
    method: 'GET'
  });

  /**
   * Indicates whether IME input is active at the moment
   *
   * @method activatedEngine
   * @see https://code.google.com/p/selenium/wiki/JsonWireProtocol#GET_/session/:sessionId/ime/activated
   * @param {GET} sessionId ID of the session to route the command to
   */

  Driver.addCommand({
    name: 'activatedEngine',
    url: '/session/:sessionId/ime/activated',
    method: 'GET'
  });

  /**
   * De-activates the currently-active IME engine
   *
   * @method deactivateEngine
   * @see https://code.google.com/p/selenium/wiki/JsonWireProtocol#POST_/session/:sessionId/ime/deactivate
   * @param {GET} sessionId ID of the session to route the command to
   */

  Driver.addCommand({
    name: 'deactivateEngine',
    url: '/session/:sessionId/ime/deactivate',
    method: 'POST'
  });

  /**
   * Make an engines that is available (appears on the list returned by getAvailableEngines) active.
   * After this call, the engine will be added to the list of engines loaded in the IME daemon
   * and the input sent using sendKeys will be converted by the active engine.
   * Note that this is a platform-independent method of activating IME
   * (the platform-specific way being using keyboard shortcuts)
   *
   * @method activateEngine
   * @see https://code.google.com/p/selenium/wiki/JsonWireProtocol#POST_/session/:sessionId/ime/activate
   * @param {GET} sessionId ID of the session to route the command to
   * @param {POST} engine Name of the engine to activate
   */

  Driver.addCommand({
    name: 'activateEngine',
    url: '/session/:sessionId/ime/activate',
    method: 'POST',
    params: ['engine']
  });

};
