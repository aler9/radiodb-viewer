module.exports = {
    presets: [
        ["@babel/env", {
            "targets": {
                "browsers": [">= 1%", "iOS >= 9"],
            },
            "debug": true,
            "useBuiltIns": "usage",
            "corejs": 3,
        }],
        ["@babel/preset-react"],
    ],
    "plugins": [
        ["@babel/plugin-proposal-class-properties"],
    ]
}
