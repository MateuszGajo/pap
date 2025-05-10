const fs = require('fs');

// For first phase i wanna test reading everything at once, node has 512mb limits for reading data, so i will split file into twos

const dataCount = 10_000_000;
const fileCount  = 2
const chunkSize = 100_000;
const totalChunks = dataCount / (chunkSize *fileCount);


const oneFile = fs.createWriteStream(`data.json`);
oneFile.write('{"pairs": [');
for (let k =0;k<fileCount;k++){
    const file = fs.createWriteStream(`data${k}.json`);
file.write('{"pairs": [');
for (let chunk = 0; chunk < totalChunks; chunk++) {
    const pairs = [];

    for (let i = 0; i < chunkSize; i++) {
        pairs.push({
            x0: Math.random() * 360 - 180,
            x1: Math.random() * 360 - 180,
            y0: Math.random() * 180 - 90,
            y1: Math.random() * 180 - 90,
        });
    }

    // Serialize chunk (excluding surrounding array brackets)
    const chunkStr = JSON.stringify(pairs).slice(1, -1);

    // Add comma if not the first chunk
    if (chunk > 0) {
        file.write(',');
       
    }
    if (!(k ==0 && chunk ==0)) {
        oneFile.write(',');
    }

    file.write(chunkStr);
    oneFile.write(chunkStr);
}
file.write(']}');
file.end();
}

oneFile.write(']}');
oneFile.end();