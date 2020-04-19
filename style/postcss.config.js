
module.exports = {
    plugins: [
        require('postcss-import')(),
        require('googlefonts-inliner')({ localPath: '/build/static/googlefonts' }),
        require('autoprefixer')(),
        require('cssnano')(),
    ],
};
