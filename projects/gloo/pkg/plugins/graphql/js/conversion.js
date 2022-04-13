function msgToBase64String(msg) {
  return uintArrayToBase64(msg.serializeBinary());
}

function base64StringToMsg(string, deserializeFunc) {
  return deserializeFunc(decodeBase64(string));
}

// function that converts uintarray8 to base64 string
function uintArrayToBase64(uintArray) {
  var binary = '';
  for (var i = 0; i < uintArray.length; i++) {
    binary += String.fromCharCode(uintArray[i]);
  }
  return Buffer.from(binary, 'binary').toString('base64');
}

// function that converts base64 string to uintarray8
// function to decode base64 string to byte array
function decodeBase64(base64) {
  var binary_string = Buffer.from(base64, 'base64').toString('binary');
  var len = binary_string.length;
  var bytes = new Uint8Array(len);
  for (var i = 0; i < len; i++) {
    bytes[i] = binary_string.charCodeAt(i);
  }
  return bytes;
}

module.exports = {
  msgToBase64String,
  base64StringToMsg,
};
