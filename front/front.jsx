
import React, { Component } from "react"
import { Link } from "react-router-dom"

import { Fetch } from "./fetch.jsx"
import { debounce } from "./various.jsx"
import { TextInput } from "./inputs.jsx"
import { imgTours } from "./imagearrays.jsx"

import "./front.scss"

import imgLogo from "./images/logo.svg"


class Search extends Component {
    state = {
        query: "",
        state: "NOQUERY",
        res: [],
    }

    curRequestId = 0 // id to block concurrent requests

    doSearch = debounce((requestId) => {
        if (this.state.state == "NOQUERY") {
            return
        }

        fetch("/data/search", {
            method: "POST",
            body: JSON.stringify({
                query: this.state.query,
            })
        })
            .then((r) => r.json())
            .then((res) => {
                if (requestId != this.curRequestId) {
                    return
                }

                this.setState({
                    state: "LOADED",
                    res: res.res,
                })
            })
    }, 400)

    onQuery = (query) => {
        this.setState({
            query,
            state: (query.length > 0) ? "LOADING" : "NOQUERY",
            res: [],
        })
        // always go on to stop any previous request
        this.doSearch(++this.curRequestId)
    }

    render() {
        return <div className="search"><div>
            <TextInput
                placeholder="Search show date, place or file TTH"
                value={ this.state.query }
                onChange={ this.onQuery }
            />

            { (this.state.state == "LOADING" || this.state.state == "LOADED") && <ul>
                { (() => {
                    if (this.state.state == "LOADING") {
                        return <li className="msg">loading...</li>
                    }

                    if (!this.state.res) {
                        return <li className="msg">no results</li>
                    }

                    return this.state.res.map((r) => {
                        return <li key={ r.url }>
                            <Link to={ r.url }>{ r.title }</Link>
                        </li>
                    })
                })() }
            </ul> }
        </div></div>
    }
}

export class Front extends Component {
    state = {}

    componentDidMount() {
        document.title = "RadioDB | The complete Radiohead bootlegs database"
    }

    render() {
        if (!this.state.loaded) {
            return <Fetch
                path="/data/front"
                body={ {} }
                onData={ (data) => this.setState(data) }
            />
        }

        return <div id="front">
            <div className="content">
                <div className="logo">
                    <img src={ "/static/" + imgLogo } alt="RadioDB" />
                </div>

                <Search />

                <div className="last">
                    <div className="shows">
                        <h3><Link to="/shows">Latest shows</Link></h3>
                        <div className="list">
                            { this.state.showsLast.map((show) => {
                                return <Link key={ show.id } to={ "/shows/" + show.id }>
                                    <span className="tour">
                                        <img src={ "/static/" + imgTours[show.tour] } />
                                    </span>
                                    <span className="title">
                                        { show.title }
                                        <br />
                                        { show.country }
                                    </span>
                                </Link>
                            }) }
                        </div>
                    </div>
                    <div className="sep"></div>
                    <div className="bootlegs">
                        <h3><Link to="/bootlegs">Latest bootlegs</Link></h3>
                        <div className="list">
                            { this.state.bootlegsLast.map((bootleg) => {
                                return <Link key={ bootleg.id } to={ "/bootlegs/" + bootleg.id }>
                                    <span className="tour">
                                        <img src={ "/static/" + imgTours[bootleg.tour] } />
                                    </span>
                                    <span className="title">
                                        { bootleg.name }
                                    </span>
                                </Link>
                            }) }
                        </div>
                    </div>
                </div>

                <div className="info">
                    <div className="stats">
                        <Link to="/stats">
                            { "Tracking " + this.state.showCount + " shows," +
                            " " + this.state.bootlegCount + " bootlegs," +
                            " " + this.state.shareSize + " of data" }
                        </Link>
                    </div>
                    Up to date { this.state.generated }
                </div>
            </div>
        </div>
    }
}
