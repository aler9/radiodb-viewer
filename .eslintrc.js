module.exports = {
    "env": {
        "browser": true,
    },
    "parser": "babel-eslint",
    "parserOptions": {
        "ecmaFeatures": {
            "legacyDecorators": true
        }
    },
    "extends": [
        "eslint:recommended",
        "plugin:react/recommended"
    ],
    "settings": {
        "react": {
            "version": "16.6",
        }
    },

    "rules": {
        "react/prop-types": "off",
        "react/no-find-dom-node": "off",
        "no-console": "off",

        "linebreak-style": ["error", "unix"],
        "indent": ["error", 4],
        "brace-style": ["error", "1tbs", { "allowSingleLine": true }],
        "semi": ["error", "never"],
        "quotes": ["error", "double"],
        "comma-spacing": "error",
        "space-before-blocks": ["error"],
        "no-var": "error",
        "curly": ["error", "multi-line"],
        "quote-props": ["error", "as-needed"],
        "space-in-parens": ["error", "never"],
        "no-unused-vars": ["error", { "args": "none" }],
        "space-infix-ops": "error",
        "no-trailing-spaces": ["error", { "ignoreComments": true }],
        "object-curly-spacing": ["error", "always"],
        "array-bracket-spacing": ["error", "always"],
        "space-before-function-paren": ["error", {
            "anonymous": "never",
            "named": "never",
            "asyncArrow": "always"
        }],
        "key-spacing": "error",
        "keyword-spacing": ["error", { "overrides": {
            "if": { "after": false },
            "for": { "after": false },
            "while": { "after": false },
            "switch": { "after": false },
        }}],
    }
}
