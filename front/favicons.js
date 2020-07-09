const favicons = require('favicons');
const path = require('path');
const { promisify } = require('util');
const readFile = promisify(require('fs').readFile);
const writeFile = promisify(require('fs').writeFile);
const mkdir = promisify(require('fs').mkdir);

const APP_NAME = "RadioDB";
const SRC_FILE = "favicon.svg";
const DEST_DIR = "/build/static/fav";
const URL_PREFIX = "/static/fav";
const TEMPLATE_FILE = "/build/frame.html";

(async () => {
    await mkdir(DEST_DIR, { recursive: true });

    let res = await favicons(SRC_FILE, {
        path: URL_PREFIX,
        appName: APP_NAME,
        icons: {
            android: false,
            appleStartup: false,
            coast: false,
            yandex: false,
            firefox: false,
            windows: false,
        },
    });

    let files = res.images.concat(res.files);
    for(let f of files) {
        await writeFile(path.join(DEST_DIR, f.name), f.contents);
    }

    await writeFile(TEMPLATE_FILE, (await readFile(TEMPLATE_FILE, 'utf8'))
        .replace(/<!-- favicons -->/, res.html.join('\n    ')));

})().catch((e) => { console.error(e); process.exit(1); });
