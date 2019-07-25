module.exports = {
    presets: [
        ["@babel/env", {
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
