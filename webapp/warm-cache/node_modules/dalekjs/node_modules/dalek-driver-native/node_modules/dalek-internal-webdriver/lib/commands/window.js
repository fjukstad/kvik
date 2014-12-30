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
 * Window related Webdriver endpoints
 *
 * @param {Dalek.Internal.Webdriver} Driver Webdriver client instance
 * @return {Dalek.Internal.Webdriver} Driver Webdriver client instance
 */

module.exports = function (Driver) {
  'use strict';

  /**
   * Retrieve the current window handle.
   *
   * @method windowHandle
   * @see https://code.google.com/p/selenium/wiki/JsonWireProtocol#GET_/session/:sessionId/window_handle
   * @param {GET} sessionId ID of the session to route the command to
   */

  Driver.addCommand({
    name: 'windowHandle',
    url: '/session/:sessionId/window_handle',
    method: 'GET'
  });

  /**
   * Retrieve the list of all window handles available to the session.
   *
   * @method windowHandles
   * @see https://code.google.com/p/selenium/wiki/JsonWireProtocol#GET_/session/:sessionId/window_handles
   * @param {GET} sessionId ID of the session to route the command to
   */

  Driver.addCommand({
    name: 'windowHandles',
    url: '/session/:sessionId/window_handles',
    method: 'GET'
  });

  /**
   * Change focus to another window.
   * The window to change focus to may be specified by its server assigned window handle, or by the value of its name attribute.
   *
   * @method changeWindow
   * @see https://code.google.com/p/selenium/wiki/JsonWireProtocol#POST_/session/:sessionId/window
   * @param {GET} sessionId ID of the session to route the command to
   * @param {POST} name Name of the window to switch to
   */

  Driver.addCommand({
    name: 'changeWindow',
    url: '/session/:sessionId/window',
    method: 'POST',
    params: ['name']
  });

  /**
   * Get the size of the specified window.
   * If the :windowHandle URL parameter is "current", the size of the currently active window will be returned.
   *
   * @method getWindowSize
   * @see https://code.google.com/p/selenium/wiki/JsonWireProtocol#GET_/session/:sessionId/window/:windowHandle/size
   * @param {GET} sessionId ID of the session to route the command to
   * @param {GET} windowHandle ID of the window to route the command to
   */

  Driver.addCommand({
    name: 'getWindowSize',
    url: '/session/:sessionId/window/current/size',
    method: 'GET'
  });

  /**
   * Change the size of the specified window.
   * If the :windowHandle URL parameter is "current",
   * the currently active window will be resized.
   *
   * @method setWindowSize
   * @see https://code.google.com/p/selenium/wiki/JsonWireProtocol#POST_/session/:sessionId/window/:windowHandle/size
   * @param {POST} sessionId ID of the session to route the command to
   * @param {POST} windowHandle ID of the window to route the command to
   * @param {POST} width The new window width
   * @param {POST} height The new window height
   */

  Driver.addCommand({
    name: 'setWindowSize',
    url: '/session/:sessionId/window/current/size',
    method: 'POST',
    params: ['width', 'height']
  });

  /**
   * Get the position of the specified window.
   * If the :windowHandle URL parameter is "current", the position of the currently active window will be returned.
   *
   * @method getWindowPosition
   * @see https://code.google.com/p/selenium/wiki/JsonWireProtocol#GET_/session/:sessionId/window/:windowHandle/position
   * @param {GET} sessionId ID of the session to route the command to
   * @param {GET} windowHandle ID of the window to route the command to
   */

  Driver.addCommand({
    name: 'getWindowPosition',
    url: '/session/:sessionId/window/:windowHandle/position',
    method: 'GET'
  });

  /**
   * Set the position of the specified window.
   *
   * @method setWindowPosition
   * @see https://code.google.com/p/selenium/wiki/JsonWireProtocol#GET_/session/:sessionId/window/:windowHandle/position
   * @param {POST} sessionId ID of the session to route the command to
   * @param {POST} x The X coordinates for the window, relative to the upper left corner of the screen.
   * @param {POST} y The Y coordinates for the window, relative to the upper left corner of the screen.
   */

  Driver.addCommand({
    name: 'setWindowPosition',
    url: '/session/:sessionId/window/current/position',
    method: 'POST',
    params: ['x', 'y']
  });

  /**
   * Maximize the specified window if not already maximized.
   * If the :windowHandle URL parameter is "current", the currently active window will be maximized.
   *
   * @method windowMaximize
   * @see https://code.google.com/p/selenium/wiki/JsonWireProtocol#GET_/session/:sessionId/window/:windowHandle/position
   * @param {GET} sessionId ID of the session to route the command to
   * @param {GET} windowHandle ID of the window to route the command to
   */

  Driver.addCommand({
    name: 'maximize',
    url: '/session/:sessionId/window/current/maximize',
    method: 'POST'
  });

  /**
   * Closes the current window.
   *
   * @method close
   * @see https://code.google.com/p/selenium/wiki/JsonWireProtocol#DELETE_/session/:sessionId/window
   * @param {GET} sessionId ID of the session to route the command to
   */

  Driver.addCommand({
    name: 'close',
    url: '/session/:sessionId/window',
    method: 'DELETE'
  });

};
