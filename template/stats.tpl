<div id="stats">

    <section>
        <h2>By year</h2>
        <canvas id="stats-peryear" data="{{ .PerYear }}"></canvas>
    </section>

    <section>
        <h2>General</h2>
        <ul>
            <li>
                <span>Date of last update</span>
                <span>{{ .Generated }}</span>
            </li>
        </ul>
    </section>

    <section>
        <h2>Files</h2>
        <ul>
            <li>
                <span>File count</span>
                <span>{{ .Stats.FileCount }}</span>
            </li>
            <li>
                <span>Unique files</span>
                <span>{{ .Stats.FileUniqueCount }}</span>
            </li>
            <li>
                <span>Overall size</span>
                <span>{{ .ShareSize }}</span>
            </li>
            <li>
                <span>Unique size</span>
                <span>{{ .ShareUniqueSize }}</span>
            </li>
        </ul>
    </section>

    <section>
        <h2>Bootlegs</h2>
        <ul>
            <li>
                <span>Bootleg count</span>
                <span>{{ .Stats.BootlegCount }}</span>
            </li>
            <li>
                <span>Audio bootlegs</span>
                <span>{{ .Stats.AudioCount }}</span>
            </li>
            <li>
                <span>Video bootlegs</span>
                <span>{{ .Stats.VideoCount }}</span>
            </li>
            <li>
                <span>Misc bootlegs</span>
                <span>{{ .Stats.MiscCount }}</span>
            </li>
            <li>
                <span>Date of last bootleg</span>
                <span>{{ .DateLastBootleg }}</span>
            </li>
        </ul>
    </section>

    <section>
        <h2>Shows</h2>
        <ul>
            <li>
                <span>Show count</span>
                <span>{{ .Stats.ShowCount }}</span>
            </li>
            <li>
                <span>Earliest show</span>
                <span>{{ .DateEarliest }}</span>
            </li>
            <li>
                <span>Latest show</span>
                <span>{{ .DateLatest }}</span>
            </li>
        </ul>
    </section>

</div>
