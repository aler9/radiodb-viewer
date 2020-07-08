
import React, { Component } from "react"
import { Redirect } from "react-router-dom"

import { Fetch } from "./fetch.jsx"


export class Random extends Component {
    state = {}

    render() {
        if (!this.state.loaded) {
            return <Fetch
                path={ "/data/random" }
                body={ {} }
                onData={ (data) => this.setState(data) }
            />
        }

        return <Redirect to={ this.state.path } />
    }
}
