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
 * Interaction related Webdriver endpoints
 *
 * @param {Dalek.Internal.Webdriver} Driver Webdriver client instance
 * @return {Dalek.Internal.Webdriver} Driver Webdriver client instance
 */

module.exports = function (Driver) {
  'use strict';

  /**
   * Move the mouse by an offset of the specificed element.
   * If no element is specified, the move is relative to the current mouse cursor.
   * If an element is provided but no offset, the mouse will be moved to the center of the element.
   * If the element is not visible, it will be scrolled into view.
   *
   * @method moveto
   * @see https://code.google.com/p/selenium/wiki/JsonWireProtocol#POST_/session/:sessionId/moveto
   * @param {GET} sessionId ID of the session to route the command to
   * @param {POST} element Opaque ID assigned to the element to move to, as described in the WebElement JSON Object. If not specified or is null, the offset is relative to current position of the mouse.
   * @param {POST} xoffset X offset to move to, relative to the top-left corner of the element. If not specified, the mouse will move to the middle of the element.
   * @param {POST} yoffset Y offset to move to, relative to the top-left corner of the element. If not specified, the mouse will move to the middle of the element.
   */

  Driver.addCommand({
    name: 'moveto',
    url: '/session/:sessionId/moveto',
    method: 'POST',
    params: ['element', 'xoffset', 'yoffset'],
    onRequest: function (params) {
      return params.element;
    }
  });

  /**
   * Click any mouse button (at the coordinates set by the last moveto command).
   * Note that calling this command after calling buttondown and before calling button up
   * (or any out-of-order interactions sequence) will yield undefined behaviour).
   *
   * @method buttonClick
   * @see https://code.google.com/p/selenium/wiki/JsonWireProtocol#POST_/session/:sessionId/click
   * @param {GET} sessionId ID of the session to route the command to
   * @param {POST} button Which button, enum: {LEFT = 0, MIDDLE = 1 , RIGHT = 2}. Defaults to the left mouse button if not specified.
   */

  // FIXME: rename button* commands to pointerClick, pointerDown, etc.
  Driver.addCommand({
    name: 'buttonClick',
    url: '/session/:sessionId/click',
    method: 'POST',
    params: ['button']
  });

  /**
   * Click and hold the left mouse button (at the coordinates set by the last moveto command).
   * Note that the next mouse-related command that should follow is buttonup.
   * Any other mouse command (such as click or another call to buttondown) will yield undefined behaviour.
   *
   * @method buttondown
   * @see https://code.google.com/p/selenium/wiki/JsonWireProtocol#POST_/session/:sessionId/buttondown
   * @param {GET} sessionId ID of the session to route the command to
   * @param {POST} button Which button, enum: {LEFT = 0, MIDDLE = 1 , RIGHT = 2}. Defaults to the left mouse button if not specified.
   */

  Driver.addCommand({
    name: 'buttondown',
    url: '/session/:sessionId/buttondown',
    method: 'POST',
    params: ['button']
  });

  /**
   * Releases the mouse button previously held (where the mouse is currently at).
   * Must be called once for every buttondown command issued.
   * See the note in click and buttondown about implications of out-of-order commands.
   *
   * @method buttonup
   * @see https://code.google.com/p/selenium/wiki/JsonWireProtocol#POST_/session/:sessionId/buttonup
   * @param {GET} sessionId ID of the session to route the command to
   * @param {POST} button Which button, enum: {LEFT = 0, MIDDLE = 1 , RIGHT = 2}. Defaults to the left mouse button if not specified.
   */

  Driver.addCommand({
    name: 'buttonup',
    url: '/session/:sessionId/buttonup',
    method: 'POST',
    params: ['button']
  });

  /**
   * Double-clicks at the current mouse coordinates (set by moveto).
   *
   * @method doubleclickPage
   * @see https://code.google.com/p/selenium/wiki/JsonWireProtocol#POST_/session/:sessionId/doubleclick
   * @param {GET} sessionId ID of the session to route the command to
   */

  Driver.addCommand({
    name: 'buttonDoubleclick',
    url: '/session/:sessionId/doubleclick',
    method: 'POST'
  });

  /**
   * Single tap on the touch enabled device.
   *
   * @method tap
   * @see https://code.google.com/p/selenium/wiki/JsonWireProtocol#POST_/session/:sessionId/touch/click
   * @param {GET} sessionId ID of the session to route the command to
   */

  // FIXME: rename touch* and tap commands to touchDown, touchUp, touchFlich, etc.
  Driver.addCommand({
    name: 'tap',
    url: '/session/:sessionId/touch/click',
    method: 'POST'
  });

  /**
   * Finger down on the screen.
   *
   * @method touchdown
   * @see https://code.google.com/p/selenium/wiki/JsonWireProtocol#POST_/session/:sessionId/touch/down
   * @param {GET} sessionId ID of the session to route the command to
   * @param {POST} x X coordinate on the screen.
   * @param {POST} y Y coordinate on the screen.
   */

  Driver.addCommand({
    name: 'touchdown',
    url: '/session/:sessionId/touch/down',
    method: 'POST',
    params: ['x', 'y']
  });

  /**
   * Finger up on the screen.
   *
   * @method touchup
   * @see https://code.google.com/p/selenium/wiki/JsonWireProtocol#POST_/session/:sessionId/touch/up
   * @param {GET} sessionId ID of the session to route the command to
   * @param {POST} x X coordinate on the screen.
   * @param {POST} y Y coordinate on the screen.
   */

  Driver.addCommand({
    name: 'touchup',
    url: '/session/:sessionId/touch/up',
    method: 'POST',
    params: ['x', 'y']
  });

  /**
   * Finger move on the screen.
   *
   * @method touchmove
   * @see https://code.google.com/p/selenium/wiki/JsonWireProtocol#POST_/session/:sessionId/touch/move
   * @param {GET} sessionId ID of the session to route the command to
   * @param {POST} x X coordinate on the screen.
   * @param {POST} y Y coordinate on the screen.
   */

  Driver.addCommand({
    name: 'touchmove',
    url: 'session/:sessionId/touch/move',
    method: 'POST',
    params: ['x', 'y']
  });

  /**
   * Scroll on the touch screen using finger based motion events.
   * Use this command if you don't care where the scroll starts on the screen.
   *
   * @method touchscroll
   * @see https://code.google.com/p/selenium/wiki/JsonWireProtocol#POST_/session/:sessionId/touch/scroll
   * @param {GET} sessionId ID of the session to route the command to
   * @param {POST} x X coordinate on the screen.
   * @param {POST} y Y coordinate on the screen.
   */

  Driver.addCommand({
    name: 'touchscroll',
    url: 'session/:sessionId/touch/scroll',
    method: 'POST',
    params: ['x', 'y']
  });

  /**
   * Double tap on the touch screen using finger motion events.
   *
   * @method doubletap
   * @see https://code.google.com/p/selenium/wiki/JsonWireProtocol#POST_/session/:sessionId/touch/doubleclick
   * @param {GET} sessionId ID of the session to route the command to
   */

  Driver.addCommand({
    name: 'doubleTap',
    url: 'session/:sessionId/touch/doubleclick',
    method: 'POST'
  });

  /**
   * Long press on the touch screen using finger motion events.
   *
   * @method longpress
   * @see https://code.google.com/p/selenium/wiki/JsonWireProtocol#POST_/session/:sessionId/touch/longclick
   * @param {GET} sessionId ID of the session to route the command to
   */

  Driver.addCommand({
    name: 'longpress',
    url: 'session/:sessionId/touch/longclick',
    method: 'POST'
  });

  /**
   * Flick on the touch screen using finger motion events.
   * This flickcommand starts at a particulat screen location.
   *
   * @method flick
   * @see https://code.google.com/p/selenium/wiki/JsonWireProtocol#POST_/session/:sessionId/touch/flick
   * @param {GET} sessionId ID of the session to route the command to
   * @param {POST} element ID of the element where the flick starts
   * @param {POST} xoffset The x offset in pixels to flick by
   * @param {POST} yoffset The y offset in pixels to flick by
   * @param {POST} speed The speed in pixels per seconds
   */

  Driver.addCommand({
    name: 'flick',
    url: 'session/:sessionId/touch/flick',
    method: 'POST',
    params: ['id', 'x', 'y', 'speed']
  });


  Driver.addCommand({
    url: '/session/:sessionId/keys',
    method: 'GET'
  });

};
