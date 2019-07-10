const webpack = require('webpack'),
    path = require('path'),
    TerserPlugin = require('terser-webpack-plugin')

module.exports = {
    mode: (process.env.BUILD_MODE == 'prod') ? 'production' : 'development',
    entry: './script/main.jsx',
    output: {
        path: '/build/static',
        filename: 'script.js',
    },
    module: { rules: [
        {
            test: /\.js$/,
            exclude: /node_modules/,
            loader: 'do-not-use',
        },
        {
            test: /\.jsx$/,
            exclude: /node_modules/,
            use: [
                { loader: 'babel-loader' },
                { loader: 'eslint-loader' }
            ]
        },
    ]},
    optimization: {
        minimizer: [new TerserPlugin({
            terserOptions: {
                output: { comments: false }
            },
        })],
    },
}
