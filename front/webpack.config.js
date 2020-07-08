
const TerserPlugin = require('terser-webpack-plugin');
const MiniCssExtractPlugin = require('mini-css-extract-plugin');
const OptimizeCSSAssetsPlugin = require('optimize-css-assets-webpack-plugin');

module.exports = {
    mode: process.env.BUILD_MODE,
    entry: './main.jsx',
    output: {
        path: '/build/static',
        filename: 'script.[hash].js',
    },
    plugins: [new MiniCssExtractPlugin({
        filename: 'style.[hash].css',
    })],
    module: { rules: [
        {
            test: /\.jsx$/,
            exclude: /node_modules/,
            use: [
                'babel-loader',
                'eslint-loader',
            ]
        },
        {
            test: /\.scss$/,
            exclude: /node_modules/,
            use: [
                MiniCssExtractPlugin.loader,
                { loader: 'css-loader', options: { url: true, import: false } },
                'postcss-loader',
            ]
        },
        {
            test: /\.svg$/,
            exclude: /node_modules/,
            use: [
                "file-loader",
                "svgo-loader",
            ]
        },
        {
            test: /\.(jpg|woff2|png|gif)$/,
            exclude: /node_modules/,
            use: [
                "file-loader",
            ]
        },
        // node_modules
        {
            test: /\.css$/,
            include: /node_modules/,
            use: [
                MiniCssExtractPlugin.loader,
                { loader: 'css-loader', options: { url: false, import: false } },
            ]
        },
    ]},
    optimization: {
        minimizer: [
            new TerserPlugin({
                terserOptions: {
                    output: { comments: false }
                },
            }),
            new OptimizeCSSAssetsPlugin(),
        ],
    },
    stats: { children: false },
};
