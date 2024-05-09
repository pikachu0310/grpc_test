const path = require('path');
const webpack = require('webpack');  // webpackをインポート
const dotenv = require('dotenv');

module.exports = {
  mode: 'development',
  entry: './app.js',  // app.jsの正確なパスを指定してください
  output: {
    path: path.resolve(__dirname, 'dist'),
    filename: 'bundle.js'
  },
  module: {
    rules: [
      {
        test: /\.js$/,
        exclude: /node_modules/,
        use: {
          loader: 'babel-loader',
          options: {
            presets: ['@babel/preset-env']
          }
        }
      }
    ]
  },
  resolve: {
    fallback: {
      "stream": require.resolve("stream-browserify"),
      "buffer": require.resolve("buffer/")
    }
  },
  plugins: [
    new webpack.ProvidePlugin({
      Buffer: ['buffer', 'Buffer'],
      'process.env.GRPC_SERVER_URL': JSON.stringify(process.env.GRPC_SERVER_URL)
    }),
    new webpack.DefinePlugin({
      'process.env.GRPC_SERVER_URL': JSON.stringify(process.env.GRPC_SERVER_URL)
    })
  ]
};
