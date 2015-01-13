'use strict';

// ext. libs
var Q = require('q');
var fs = require('fs');
var PNG = require('node-pngjs').PNG;

/**
 * Execute related methods
 *
 * @module Driver
 * @class Image
 * @namespace Dalek.DriverNative.Image
 */

var Image = {
  imagecompareResults : {
    equal         : 'equal',
    sizeDifferent : 'images have different sizes',
    imageDiffent  : 'images are different'
  },

  imagecompare : function (opts, expected, makediff, hash) {
    this.actionQueue.push(this._imagecompareCb.bind(this, opts, expected, makediff, hash));
    return this;
  },

  imagecut : function (opts, hash) {
    var position = {};

    this.actionQueue.push(this.webdriverClient.element.bind(this.webdriverClient, opts.selector));
    this.actionQueue.push(this.webdriverClient.location.bind(this.webdriverClient, opts.selector));

    /* get left top positions of element*/
    this.actionQueue.push(function (result) {
      var value = JSON.parse(result).value;
      position.x = Math.ceil(value.x);
      position.y = Math.ceil(value.y);
    });

    this.actionQueue.push(this.webdriverClient.size.bind(this.webdriverClient, opts.selector));
    this.actionQueue.push(this._imagecutCb.bind(this, opts, hash, opts.selector, position));

    return this;
  },

  _imagecutCb: function (opts, hash, selector, position, result) {
    var rect;
    var size = JSON.parse(result);
    var deferred = Q.defer();

    rect = this._convertToRect(position, size.value);
    var cb = function () {
      this.events.emit('driver:message', {key: 'imagecut', value: selector, uuid: hash, hash: hash});
      deferred.resolve();
    };

    this._doImagecut(opts.realpath, opts.selector, rect, cb.bind(this));

    return deferred.promise;
  },

  _doCalcSize: function (dstPos, dstSize, srcSize) {
    if (dstPos + dstSize <= srcSize) {
      return dstSize;
    } else if (dstSize <= srcSize) {
      return srcSize - dstPos;
    }
    return 0;
  },

  _doImagecut: function (path, element, rect, cb) {
    this._parsePNG(path).then(function (image) {
      var width = this._doCalcSize (rect.x, rect.width, image.png.width);
      var height = this._doCalcSize (rect.y, rect.height, image.png.height);

      var pngDst = new PNG({
        width :width,
        height:height
      });

      if (width === 0 || height === 0) {
        cb();
        return 'Error in image cutting. Size of element is bigger of size image.';
      }

      image.png.bitblt(pngDst, rect.x, rect.y, width, height, 0, 0);
      this._savePNG(pngDst, image.path).then(cb);

    }.bind(this));
  },

  _imagecompareCb: function (opts, expected, makediff, hash) {
    var deferred = Q.defer();
    var current = opts.realpath;

    var cb = function (result) {
      this.events.emit('driver:message', {key: 'imagecompare', value: result, uuid: hash, hash: hash});
      deferred.resolve();
    };

    this._doImagecompare(current, expected, makediff, cb.bind(this));

    return deferred.promise;
  },
  _doImagecompare: function (actualPath, etalonPath, makediff, cb) {
    var promises = [this._parsePNG(actualPath), this._parsePNG(etalonPath)];

    Q.all(promises)
    .spread(function (parsedImageA, parsedImageB) {
      var result = this._checkIfImagesTheSameAndMakeDiff(parsedImageA, parsedImageB, makediff);
      if (result === this.imagecompareResults.imageDiffent && makediff) {
        this._doSaveImageDiff(parsedImageB, result)
        .then(function (filePath) {
          cb(result + '. See differences by path ' + filePath);
        },
        cb);
      } else {
        cb(result);
      }
    }.bind(this),
    cb);
  },
  _doSaveImageDiff : function (image) {
    var imageDiffPath = image.path.replace('.png', '.diff.png');

    return this._savePNG(image.png, imageDiffPath);
  },
  _checkIfImagesTheSameAndMakeDiff : function (a, b, makediff) {

    if (a.data.length !== b.data.length) {
      return this.imagecompareResults.sizeDifferent;
    }

    var result = this.imagecompareResults.equal;

    for (var i = 0, len = a.data.length; i < len; i += 4) {
      if (a.data[i]     !== b.data[i] ||
          a.data[i + 1] !== b.data[i + 1] ||
          a.data[i + 2] !== b.data[i + 2]) {
        result = this.imagecompareResults.imageDiffent;
        if (makediff) {
          this._setErrorPixel(b.data, i);
        } else {
          return result;
        }
      }
    }

    return result;
  },
  _parsePNG: function (path) {
    var deferred = Q.defer();

    fs.createReadStream(path)
      .pipe(new PNG({
        filterType: 1
      }))
      .on('parsed', function() {

        var parsedImage = {
          path   : path,
          data   : this.data,
          png    : this
        };

        deferred.resolve(parsedImage);
      })
      .on('error', function (event) {
        deferred.reject(event);
      });

    return deferred.promise;
  },

  _savePNG : function (png, filePath) {
    var deferred = Q.defer();

    png
    .pack()
    .pipe(fs.createWriteStream(filePath, {flags: 'w'}))
    .on('error', function (errors) {
      deferred.reject(errors);
    })
    .on('close', function () {
      deferred.resolve(filePath);
    });

    return deferred.promise;
  },
  _setErrorPixel : function (dst, index) {
    var errorPixelColor = { // Color for Error Pixels. Between 0 and 255.
      red: 255,
      green: 0,
      blue: 255,
      alpha: 255
    };

    dst[index]     = errorPixelColor.red;
    dst[index + 1] = errorPixelColor.green;
    dst[index + 2] = errorPixelColor.blue;
    dst[index + 3] = errorPixelColor.alpha;
  },
  _convertToRect: function (position, size) {
    return {
      x:position.x,
      y:position.y,
      width : size.width,
      height: size.height
    };
  }
};

module.exports = function (DalekNative) {
  // mixin methods

  Object.keys(Image).forEach(function (fn) {
    DalekNative.prototype[fn] = Image[fn];
  });

  return DalekNative;
};

