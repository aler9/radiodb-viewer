module.exports = {
    plugins: [
        require("postcss-import")(),
        require("stylelint")(),
        require('postcss-mixins')(),
        require('postcss-simple-vars')(),
        require('postcss-nested')(),
        require('googlefonts-inliner')({ localPath: '/build/static/googlefonts' }),
        require('autoprefixer')(),
        require("postcss-reporter")({ throwError: true }),
    ],
};
