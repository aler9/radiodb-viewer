
<span class="tour">
    <img src="/static/{{ .Tour }}.png" title="{{ .TourLabel }}" />
</span>

<span class="country">
    <img src="/static/countryflag/{{ .CountryCodeShort }}.png" title="{{ .Country }}" />
</span>

<span class="title">
    <div>
        {{ .ArtistLong }}
    </div>
    <h3>
        {{ .Date }}, {{ .City }}, {{ .Country }}
    </h3>
</span>

<span class="count">
    {{ if gt .AudioCount 0 }}
    <span>
        <img src="/static/audio.svg" />
        {{ .AudioCount }} audio
    </span>
    {{ end }}
    {{ if gt .VideoCount 0 }}
    <span>
        <img src="/static/video.svg" />
        {{ .VideoCount }} video{{ if gt .VideoCount 1 }}s{{ end }}
    </span>
    {{ end }}
    {{ if gt .MiscCount 0 }}
    <span>
        <img src="/static/misc.svg" />
        {{ .MiscCount }} misc{{ if gt .MiscCount 1 }}s{{ end }}
    </span>
    {{ end }}
</span>
