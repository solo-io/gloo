/* eslint-disable */
const webpack = require('webpack');

module.exports = function configOverload(config) {
    const fallback = {
        ...(config.resolve.fallback || {}),
        ...{
            buffer: require.resolve("buffer"),
            crypto: require.resolve("crypto-browserify"),
            stream: require.resolve("stream-browserify"),
            assert: require.resolve("assert"),
            http: require.resolve("stream-http"),
            https: require.resolve("https-browserify"),
            os: require.resolve("os-browserify"),
            url: require.resolve("url"),
            process: require.resolve("process"),
        }
    };

   config.resolve.fallback = fallback;
   config.resolve.alias["process/browser"] = "process/browser.js";

   config.ignoreWarnings = [
       {
           module: /.*/,
           message: /.*Failed to parse source map.*/
       },
   ]

   config.plugins = (config.plugins || []).concat([
   	new webpack.ProvidePlugin({
    	process: 'process/browser',
        Buffer: ['buffer', 'Buffer']
    }),
   ]);
   return config;
}
