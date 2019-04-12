<div id="show">

    <div class="title">
        <div class="parent">
            <a href="/shows">Show</a>
        </div>
        <h1>
            {{ .ArtistLong }}, {{ .Date }}, {{ .City }}, {{ .LabelCountry }}
        </h1>
    </div>

    <section class="info">
        <h2>Data</h2>
        <ul>
            <li>
                <span>Artist</span>
                <span>{{ .ArtistLong }}</span>
            </li>
            <li>
                <span>Date</span>
                <span>{{ .Date }}</span>
            </li>
            <li>
                <span>Location</span>
                <span>{{ .City }}</span>
            </li>
            <li class="country">
                <span>Country</span>
                <span><img src="/static/countryflag/{{ .LabelCountryCode }}.png" /> {{ .LabelCountry }}</span>
            </li>
            <li class="tour">
                <span>Tour</span>
                <span><img src="/static/{{ .Tour }}.png" /> {{ .LabelTour }}</span>
            </li>
            <li>
                <span>Setlist</span>
                <span>
                    {{- range .Urls }}
                    <div>
                        <a href="{{ .Url }}" target="_blank">view on {{ .Name }}</a>
                    </div>
                    {{- end }}
                </span>
            </li>
        </ul>
    </section>

    <section class="bootlegs">
        <h2>Bootlegs</h2>
        <div>
            {{- range .Bootlegs }}
            <a href="/bootlegs/{{ .Id }}">
                <span class="name">
                    {{ .Name }}
                    <div class="released">Released: {{ .FirstSeen }}</div>
                </span>
                <span class="duration">
                    {{ .Duration }}
                </span>
                <span class="res" title="{{ .TypeLong }}">
                    <img class="media-type-icon" src="/static/{{ .Type }}.svg" />
                    {{ .Res }}
                </span>
                <span class="size">{{ .Size }}</span>
            </a>
            {{- end }}
        </div>
    </section>

</div>
