'use strict'

exports.handler = (event, context, callback) => {
    const response = {
        statusCode: 200,
        body: 'AWS Lambda, brought to you by Gloo.'
    };
    callback(null, response);
};
