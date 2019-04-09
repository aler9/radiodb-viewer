module.exports = {
    plugins: [
        require('postcss-import')({ path: [__dirname] }),
        require('autoprefixer')({ browsers: [">= 1%", "iOS >= 9"] }),
        require('cssnano')(),
    ],
}
