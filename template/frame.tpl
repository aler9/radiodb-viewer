<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="utf-8">
    <title>{{ .Title }}</title>
    <meta name="viewport" content="width=device-width,maximum-scale=1.0">

    <!-- favicons -->

    <style>
    /* critical css */
    body {
        background: rgb(35, 35, 35);
        overflow-y: scroll; /* always show scrollbar */
    }
    body > .inner { display: none; }
    </style>
    <link rel="stylesheet" href="/static/style.css">
    <script defer src="/static/script.js"></script>
</head>
<body>
<div class="inner">

<header>
    {{- if ne .CurPath "/" }}
    <a class="logo" href="/">
        <img src="/static/logo.svg" alt="RadioDB" />
    </a>
    {{- end }}

    <input type="checkbox" id="menu-toggle" />
    <label for="menu-toggle">
        <img src="/static/menu.svg" alt="" />
    </label>
    <nav><ul>
        <li><a href="/">Home</a></li>
        <li><a href="/shows" class="{{ if eq .CurPath "/shows" }}current{{ end }}">shows</a></li>
        <li><a href="/bootlegs" class="{{ if eq .CurPath "/bootlegs" }}current{{ end }}">bootlegs</a></li>
        <li><a href="/random" class="{{ if eq .CurPath "/random" }}current{{ end }}">Random</a></li>
        <li><a href="/stats" class="{{ if eq .CurPath "/stats" }}current{{ end }}">Stats</a></li>
        <li><a href="/dump" class="{{ if eq .CurPath "/dump" }}current{{ end }}">Dump</a></li>
    </ul></nav>
</header>

<main{{ if .Class }} class="{{ .Class }}"{{ end }}>
{{ .Content -}}
</main>

<footer>
    <a href="https://github.com/aler9/radiodb-viewer" target="_blank" class="source">
        <img src="/static/github.svg" />
        source code
    </a>
    <br />
    <br />
    This site is not affiliated with any of the mentioned artists.
    <br />
    Use of low-resolution album covers qualifies as fair use under the copyright law of the United States.
    <br />
    This site does not host any file and is merely a tracker.
    <br />
    Audio and video materials mentioned in the site are exclusively non-commercial, non-official and freely available recordings.
    <br />
    If you want to have an entry removed, write to <a href="mailto:aler9.dev@gmail.com">aler9.dev@gmail.com</a>
</footer>

</div>
</body>
</html>
