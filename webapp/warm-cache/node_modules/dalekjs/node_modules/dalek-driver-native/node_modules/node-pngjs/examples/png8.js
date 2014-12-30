#!/usr/bin/env node

var fs = require('fs'),
    PNG = require('../lib/png').PNG;


var pngSrc = new PNG({
		filterType: -1
	});

pngSrc.on('parsed', function() {
	var pngDst = new PNG({
			width: pngSrc.width,
			height:pngSrc.height,
			bitDepth: 8,
			colorType: 3,
			filterType: -1
		});
	pngDst.data=new Buffer(pngDst.width*pngDst.height*4);
    for (var y = 0; y < pngSrc.height; y++) {
        for (var x = 0; x < pngSrc.width; x++) {
			for(var k = 0; k < 4; k++)
			{
				pngDst.data[(pngSrc.width * y + x)*4+k] = pngSrc.data[(pngSrc.width * y + x)*4+k];
			}
        }
    }
	pngDst.pack().pipe(fs.createWriteStream(process.argv[3] || 'out.png',{flags: 'w'}));
});

fs.createReadStream(process.argv[2]).pipe(pngSrc);
