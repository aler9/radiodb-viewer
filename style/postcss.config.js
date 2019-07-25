module.exports = {
    plugins: [
        require('postcss-import')({ path: [__dirname] }),
        require('autoprefixer')(),
        require('cssnano')(),
    ],
}
