module.exports = {
    parser: 'postcss-scss',
    plugins: [
        require("postcss-import")(),
        require("stylelint")(),
        require('postcss-mixins')(),
        require('postcss-simple-vars')(),
        require('postcss-nested')(),
        require('googlefonts-inliner')(),
        require('autoprefixer')(),
        require("postcss-reporter")({ throwError: true }),
    ],
};
