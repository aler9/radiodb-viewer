
import React, { Component } from "react"

import { DynamicReactRender, debounce } from "./various.jsx"
import { TextInput } from "./inputs.jsx"

export class Search extends Component {
    state = {
        query: "",
        state: "noquery", // noquery, loading, loaded
        res: [],
    }

    curRequestId = 0 // id to block concurrent requests

    doSearch = debounce((requestId) => {
        if(this.state.state == "noquery") {
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
                if(requestId != this.curRequestId) {
                    return
                }

                this.setState({
                    state: "loaded",
                    res: res.res,
                })
            })
    }, 400)

    onQuery = (query) => {
        this.setState({
            query,
            state: (query.length > 0) ? "loading" : "noquery",
            res: [],
        })
        // always go on to stop any previous request
        this.doSearch(++this.curRequestId)
    }

    render() {
        return <div>
            <TextInput placeholder={ this.props.placeholder } value={ this.state.query } onChange={ this.onQuery } />
            { (this.state.state == "loading" || this.state.state == "loaded") && <ul>
                { (() => {
                    if(this.state.state == "loading") {
                        return <li className="msg" key="loading">loading...</li>
                    }
                    if(!this.state.res) {
                        return <li className="msg" key="nores">no results</li>
                    }
                    return this.state.res.map((r) => {
                        return <li key={ r.url }>
                            <a href={ r.url }>{ r.title }</a>
                        </li>
                    })
                })() }
            </ul> }
        </div>
    }
}

DynamicReactRender((dom) => <Search
    placeholder={ dom.querySelector("input").getAttribute("placeholder") } />,
"search-front")
