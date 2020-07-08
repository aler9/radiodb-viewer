
import React, { Component } from "react"

import "./dump.scss"

import imgDump from "./images/dump.svg"


export class Dump extends Component {
    componentDidMount() {
        document.title = "RadioDB - Dump"
    }

    render() {
        return <div id="dump">
            The entire database can be downloaded in JSON format by using the button below.
            <br />
            Keep in mind that the database does not include any media file, but just
            the metadata displayed in the website.
            <br />
            <a href="/dumpget">
                Get dump
                <img src={ "/static/" + imgDump } />
            </a>
        </div>
    }
}
