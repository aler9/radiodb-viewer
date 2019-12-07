
module.exports = {
    "extends": "stylelint-config-recommended",
    "rules": {
        "no-descending-specificity": null,
        "at-rule-no-unknown": [true, {
            "ignoreAtRules": ["mixin", "include"]
        }],
        "indentation": 4,
        "block-opening-brace-newline-after": "always",
        "block-opening-brace-space-before": "always",
        "declaration-colon-space-after": "always",
        "property-blacklist": [ "float" ],
    },
};
