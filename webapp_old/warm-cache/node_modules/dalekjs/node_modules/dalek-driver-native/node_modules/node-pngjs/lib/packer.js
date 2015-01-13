// Copyright (c) 2012 Kuba Niegowski
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
// THE SOFTWARE.

'use strict';


var util = require('util'),
    Stream = require('stream'),
    zlib = require('zlib'),
    Filter = require('./filter'),
    CrcStream = require('./crc'),
    constants = require('./constants');


var Packer = module.exports = function(options) {
    Stream.call(this);

    this._options = options;
	options.bitDepth = options.bitDepth || 8;
	options.colorType = (typeof options.colorType=="number")? options.colorType : 6;
	options.deflateChunkSize = options.deflateChunkSize || 32 * 1024;
    options.deflateLevel = options.deflateLevel || 9;
    options.deflateStrategy = options.deflateStrategy || 3;

    this.readable = true;
};
util.inherits(Packer, Stream);


Packer.prototype.pack = function(data, width, height) {

    // Signature
    this.emit('data', new Buffer(constants.PNG_SIGNATURE));
    this.emit('data', this._packIHDR(width, height, this._options.bitDepth, this._options.colorType));

	if(this._options.colorType==3){
		this.emit('data', this._packPLTE(data));
	}

    // filter pixel data
    var filter = new Filter(width, height, colorTypeToBppMap[this._options.colorType], data, this._options);
    this._data = filter.filter();

    // compress it
    var deflate = zlib.createDeflate({
            chunkSize: this._options.deflateChunkSize,
            level: this._options.deflateLevel,
            strategy: this._options.deflateStrategy
        });
    deflate.on('error', this.emit.bind(this, 'error'));

    deflate.on('data', function(data) {
        this.emit('data', this._packIDAT(data));
    }.bind(this));

    deflate.on('end', function() {
        this.emit('data', this._packIEND());
        this.emit('end');
    }.bind(this));

    deflate.end(this._data);
};

Packer.prototype._packChunk = function(type, data) {

    var len = (data ? data.length : 0),
        buf = new Buffer(len + 12);

    buf.writeUInt32BE(len, 0);
    buf.writeUInt32BE(type, 4);

    if (data) data.copy(buf, 8);

    buf.writeInt32BE(CrcStream.crc32(buf.slice(4, buf.length - 4)), buf.length - 4);
    return buf;
};

Packer.prototype._packIHDR = function(width, height, bitDepth, colorType) {
    var buf = new Buffer(13);
    buf.writeUInt32BE(width, 0);
    buf.writeUInt32BE(height, 4);
    buf[8] = bitDepth;
    buf[9] = colorType; // colorType
    buf[10] = 0; // compression
    buf[11] = 0; // filter
    buf[12] = 0; // interlace

    return this._packChunk(constants.TYPE_IHDR, buf);
};

Packer.prototype._packIDAT = function(data) {
    return this._packChunk(constants.TYPE_IDAT, data);
};

Packer.prototype._packIEND = function() {
    return this._packChunk(constants.TYPE_IEND, null);
};
var colorTypeToBppMap = {
	0: 1,
	2: 3,
	3: 1,
	4: 2,
	6: 4
};

Packer.prototype._packPLTE = function(data) {
	var k,
		newData=new Buffer(data.length),
		buf = new Buffer(256*3),
		count,palette;
	function gc(color,k)
	{
		return parseInt(color-(color%Math.pow(2,k/8))+Math.pow(2,k/8-1));
	}
	//颜色量化算法
	for(k=0;;k++){
		palette={};
		count= 0;
		for(var i=data.length/4-1;i>=0;i--){//find all colors
			if(data[i*4+3]<0x80){//opacity<50% treat as  transparent
				if(palette._===undefined){
					palette._=(count++);
					if(count>256){break;}
				}
				newData[i*4]=palette._;
			}
			else{
				var colorR =gc(data[i*4+0],k),
					colorG =gc(data[i*4+1],k),
					colorB =gc(data[i*4+2],k),
					color=(colorR<<16)+(colorG<<8)+colorB;
				if(palette[color]===undefined){
					palette[color]=(count++);
					if(count>256){break;}
				}
				newData[i*4]=palette[color];
				buf[palette[color]*3+0]=colorR;
				buf[palette[color]*3+1]=colorG;
				buf[palette[color]*3+2]=colorB;
			}
		}
		if(count<=256){break;}
	}
	if(k>0 && count<256)
	{//如果颜色值不够256色,则将颜色表按照色差大小顺序还原
		var list= [];
		for (var i = data.length / 4 - 1; i >= 0; i--) {//find all colors
			if (newData[i * 4] === palette._) {continue;}
			var diffR = Math.abs(buf[newData[i*4]*3+0] - data[i * 4 + 0]),
				diffG = Math.abs(buf[newData[i*4]*3+1] - data[i * 4 + 1]),
				diffB = Math.abs(buf[newData[i*4]*3+2] - data[i * 4 + 2]),
				diff = diffR + diffG + diffB;
			if(!diff){continue;}
			list.push([diff,i]);
		}
		list.sort(function(a,b){return a[0]-b[0];});
		for(var i=list.length-1;i>=0;i--)
		{
			if(count>=256){break;}
			var colorR =data[list[i][1]*4+0],
				colorG =data[list[i][1]*4+1],
				colorB =data[list[i][1]*4+2],
				color=(colorR<<16)+(colorG<<8)+colorB;
			palette[color]=(count++);
			newData[list[i][1]*4]=palette[color];
			buf[palette[color]*3+0]=colorR;
			buf[palette[color]*3+1]=colorG;
			buf[palette[color]*3+2]=colorB;
		}
	}
	newData.copy(data);
	return this._packChunk(constants.TYPE_PLTE, buf.slice(0,count*3));
};