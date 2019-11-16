
const { promisify } = require("util"),
    readFile = promisify(require("fs").readFile),
    writeFile = promisify(require("fs").writeFile),
    mkdir = promisify(require("fs").mkdir),
    { dirname, join } = require("path"),
    fetch = require("node-fetch"),
    postcss = require("postcss"),
    cssnano = require("cssnano");

const replaceAsync = async (str, regex, cb) => {
    const promises = [];
    str.replace(regex, (...m) => {
        promises.push(cb(...m));
    });
    const data = await Promise.all(promises);
    return str.replace(regex, () => data.shift());
};

(async () => {
    let css = await readFile(process.argv[2], "utf-8");

    const fontsDir = join(dirname(process.argv[2]), "fonts");
    await mkdir(fontsDir);

    // download and replace css
    css = await replaceAsync(css, /@import url\(["'](https:\/\/fonts\.googleapis\.com.+?)["']\);/g, async (...m) => {
        let fcss = await fetch(m[1], {
            headers: {
                "User-Agent": "Mozilla/5.0 (iPad; CPU OS 10_3_3 like Mac OS X) AppleWebKit/603.1.30 (KHTML, like Gecko) CriOS/63.0.3239.73 Mobile/14G60 Safari/602.1",
            },
        });
        fcss = await fcss.text();

        // download and replace fonts
        fcss = await replaceAsync(fcss, /url\((https:\/\/.+\/(.+?))\)/g, async (...m) => {
            let font = await fetch(m[1]);
            font = await font.arrayBuffer();
            await writeFile(join(fontsDir, m[2]), new Buffer(font));
            return "url(fonts/" + m[2] + ")";
        });
        return fcss;
    });

    // reapply compression
    css = (await postcss([ cssnano ]).process(css, { from: undefined })).css;

    await writeFile(process.argv[2], css, "utf-8");

})().catch((e) => { console.error(e); process.exit(1); });
