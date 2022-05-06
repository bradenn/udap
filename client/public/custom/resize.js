// Copyright (c) 2022 Braden Nicholson

let glob = require("glob")
const sharp = require('sharp');
// options is optional
glob("**/*@4x.png", {}, function (er, files) {
    files.forEach(f => {
        sharp(f).resize({height: 640}).toFile(f.replace("@4x", "@2x"))
            .then(function (newFileInfo) {
                // newFileInfo holds the output file properties
                console.log("Success")
            })
            .catch(function (err) {
                console.log("Error occured" + err);
            });
    })

})