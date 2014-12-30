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
 * Page information related Webdriver endpoints
 *
 * @param {Dalek.Internal.Webdriver} Driver Webdriver client instance
 * @return {Dalek.Internal.Webdriver} Driver Webdriver client instance
 */

module.exports = function (Driver) {
  'use strict';

  /**
   * Get the current page source
   *
   * @method source
   * @see https://code.google.com/p/selenium/wiki/JsonWireProtocol#GET_/session/:sessionId/title
   * @param {GET} sessionId ID of the session to route the command to
   */

  Driver.addCommand({
    name: 'source',
    url: '/session/:sessionId/source',
    method: 'GET'
  });

  /**
   * Get the current page title
   *
   * @method title
   * @see https://code.google.com/p/selenium/wiki/JsonWireProtocol#GET_/session/:sessionId/title
   * @param {GET} sessionId ID of the session to route the command to
   */

  Driver.addCommand({
    name: 'title',
    url: '/session/:sessionId/title',
    method: 'GET'
  });

  /**
   * Checks the device orientation
   *
   * @method orientation
   * @see https://code.google.com/p/selenium/wiki/JsonWireProtocol#GET_/session/:sessionId/orientation
   * @param {GET} sessionId ID of the session to route the command to
   */

  Driver.addCommand({
    url: '/session/:sessionId/orientation',
    method: 'GET'
  });

  /**
   * Checks a prompt text
   *
   * @method alertText
   * @see https://code.google.com/p/selenium/wiki/JsonWireProtocol#GET_/session/:sessionId/alert_text
   * @param {GET} sessionId ID of the session to route the command to
   */

  Driver.addCommand({
    name: 'alertText',
    url: '/session/:sessionId/alert_text',
    method: 'GET'
  });

  /**
   * Sets a prompt text
   *
   * @method promptText
   * @see https://code.google.com/p/selenium/wiki/JsonWireProtocol#POST_/session/:sessionId/alert_text
   * @param {GET} sessionId ID of the session to route the command to
   * @param {POST} text Text to set
   */

  Driver.addCommand({
    name: 'promptText',
    url: '/session/:sessionId/alert_text',
    method: 'POST',
    params: ['text']
  });

  /**
   * Accept an dialog box
   *
   * @method acceptAlert
   * @see https://code.google.com/p/selenium/wiki/JsonWireProtocol#POST_/session/:sessionId/accept_alert
   * @param {GET} sessionId ID of the session to route the command to
   */

  Driver.addCommand({
    name: 'acceptAlert',
    url: '/session/:sessionId/accept_alert',
    method: 'POST'
  });

  /**
   * Cancel an dialog box
   *
   * @method dismissAlert
   * @see https://code.google.com/p/selenium/wiki/JsonWireProtocol#POST_/session/:sessionId/dismiss_alert
   * @param {GET} sessionId ID of the session to route the command to
   */

  Driver.addCommand({
    name: 'dismissAlert',
    url: '/session/:sessionId/dismiss_alert',
    method: 'POST'
  });

  /**
   * Get the current geo location
   *
   * @method geoLocation
   * @see https://code.google.com/p/selenium/wiki/JsonWireProtocol#GET_/session/:sessionId/location
   * @param {GET} sessionId ID of the session to route the command to
   */

  Driver.addCommand({
    name: 'geoLocation',
    url: '/session/:sessionId/location',
    method: 'GET'
  });

  /**
   * Set the geo location
   *
   * @method setGeoLocation
   * @see https://code.google.com/p/selenium/wiki/JsonWireProtocol#POST_/session/:sessionId/location
   * @param {GET} sessionId ID of the session to route the command to
   * @param {POST} latitude The new location
   * @param {POST} longitude The new location
   * @param {POST} altitude The new location
   */

  Driver.addCommand({
    name: 'setGeoLocation',
    url: '/session/:sessionId/location',
    method: 'POST',
    params: ['latitude', 'longitude', 'altitude']
  });

  /**
   * Get the log for a given log type.
   * Log buffer is reset after each request.
   *
   * @method log
   * @see https://code.google.com/p/selenium/wiki/JsonWireProtocol#POST_/session/:sessionId/log
   * @param {GET} sessionId ID of the session to route the command to
   */

  Driver.addCommand({
    name: 'log',
    url: '/session/:sessionId/log',
    method: 'GET'
  });

  /**
   * Get available log types
   *
   * @method logTypes
   * @see https://code.google.com/p/selenium/wiki/JsonWireProtocol#GET_/session/:sessionId/log/types
   * @param {GET} sessionId ID of the session to route the command to
   */

  Driver.addCommand({
    name: 'logTypes',
    url: '/session/:sessionId/log/types',
    method: 'GET'
  });

};
