<div id="bootleg">

    <div class="title">
        <div class="parent">
            <a href="/bootlegs">Bootleg</a>
        </div>
        <h1>
            {{ .Name }}
        </h1>
    </div>

    <section class="info">
        <h2>Data</h2>
        <ul>
            <li class="show">
                <span>Show</span>
                <span>
                    <a href="{{ .ShowUrl }}">
                        <img src="/static/{{ .Tour }}.png" title="{{ .TourLabel }}" />
                        <img src="/static/countryflag/{{ .CountryCodeShort }}.png" title="{{ .CountryLabel }}" />
                        <span>
                            {{ .ShowArtist }}
                            <br />{{ .ShowDate }}
                            <br />
                            {{ .ShowCity }}, {{ .ShowCountry }}
                        </span>
                    </a>
                </span>
            </li>
            <li>
                <span>Released</span>
                <span>{{ .FirstSeen }}</span>
            </li>
            <li>
                <span>Overall size</span>
                <span>{{ .Size }}</span>
            </li>
            <li>
                <span>Media type</span>
                <span><img class="media-type-icon" src="/static/{{ .Type }}.svg" alt="" /> {{ .TypeLong }}</span>
            </li>
            {{- if ne .Type "misc" }}
            <li>
                <span>Format</span>
                <span>{{ .MinfoFormat }}</span>
            </li>
            <li>
                <span>Duration</span>
                <span>{{ .Duration }}</span>
            </li>
            {{- end }}
            {{- if eq .Type "video" }}
            <li>
                <span>Video Codec</span>
                <span>{{ .MinfoVideoCodec }}</span>
            </li>
            <li>
                <span>Video resolution</span>
                <span>{{ .MinfoVideoRes }}</span>
            </li>
            {{- end }}
            {{- if eq .Type "video" "audio" }}
            <li>
                <span>Audio Codec</span>
                <span>{{ .MinfoAudioCodec }}</span>
            </li>
            <li>
                <span>Audio resolution</span>
                <span>{{ .MinfoAudioRes }}</span>
            </li>
            {{- end }}
        </ul>
    </section>

    {{- if .Finfo }}
    <section class="finfo">
        <h2>Info</h2>
        <div id="bootleg-finfo">
            {{ .Finfo }}
        </div>
    </section>
    {{- end }}

    <section class="files">
        <h2>Files</h2>
        <div class="note">
            Download requires a Direct Connect client installed
            (Windows: <a href="http://dcplusplus.sourceforge.net/" target="_blank">DC++</a>,
            macOS/Linux: <a href="https://sourceforge.net/projects/eiskaltdcpp/files/macOS/" target="_blank">EiskaltDC++</a>)
            and connected to the RadioHub (<a href="http://radiohub.wikidot.com/" target="_blank">instructions</a>).
        </div>

        <div class="list">
            {{- range .Files }}
            <a href="{{ .Magnet }}">
                <span class="name">{{ .Name }}</span>
                <span class="duration">{{ .Duration }}</span>
                <span class="size">{{ .Size }}</span>
                <span class="tth">{{ .TTH }}</span>
            </a>
            {{- end }}
        </div>
    </section>
</div>
