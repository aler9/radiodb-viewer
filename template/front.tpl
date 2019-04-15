<div id="front">
    <div class="content">
        <div class="logo">
            <img src="/static/logo.svg" alt="RadioDB" />
        </div>

        <div id="search-front" class="search"><div>
            <div class="text-input">
                <input type="text" placeholder="Search show date, place or file TTH" />
            </div>
        </div></div>

        <div class="last">
            <div class="shows">
                <h3><a href="/shows">Latest shows</a></h3>
                <div class="list">
                    {{ range .ShowsLast }}
                    <a href="/shows/{{ .Id }}">
                        <span class="tour">
                            <img src="/static/{{ .Tour }}.png" />
                        </span>
                        <span class="title">
                            {{ .Title }}
                            <br />
                            {{ .Country }}
                        </span>
                    </a>
                    {{ end }}
                </div>
            </div>
            <div class="sep"></div>
            <div class="bootlegs">
                <h3><a href="/bootlegs">Latest bootlegs</a></h3>
                <div class="list">
                    {{ range .BootlegsLast }}
                    <a href="/bootlegs/{{ .Id }}">
                        <span class="tour">
                            <img src="/static/{{ .Tour }}.png" />
                        </span>
                        <span class="title">
                            {{ .Name }}
                        </span>
                    </a>
                    {{ end }}
                </div>
            </div>
        </div>

        <div class="info">
            <div class="stats">
                <a href="/stats">
                    Tracking {{ .ShowCount }} shows,
                    {{ .BootlegCount }} bootlegs,
                    {{ .ShareSize }} of data
                </a>
            </div>
            Up to date {{ .Generated }}
        </div>
    </div>
</div>
