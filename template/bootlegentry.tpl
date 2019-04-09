
<span class="tour">
    <img src="/static/{{ .Tour }}.png" title="{{ .TourLabel }}" />
</span>

<span class="name">
    <span class="show">
        {{ .Show }}
    </span>
    <h2>
        {{ .Name }}
    </h2>
    <span class="released">
        Released: {{ .FirstSeen }}
    </span>
</span>

<span class="duration">
    {{ .Duration }}
</span>

<span class="res" title="{{ .TypeLong }}">
    <img class="media-type-icon" src="/static/{{ .Type }}.svg" /> {{ .Res }}
</span>

<span class="size">
    {{ .Size }}
</span>
