
import React, { Component, Fragment } from "react"
import { Link } from "react-router-dom"

import { Scrollable } from "./scrollable.jsx"
import { Fetch } from "./fetch.jsx"
import { imgMediaTypes, imgTours, imgCountries } from "./imagearrays.jsx"

import "./bootleg.scss"


export class Bootleg extends Component {
    state = {}

    componentDidUpdate(prevProps, prevState) {
        if (this.state.loaded && !prevState.loaded) {
            document.title = "RadioDB - " + this.state.title
        }
    }

    render() {
        if (!this.state.loaded) {
            return <Fetch
                path={ "/data/bootleg/" + this.props.match.params.id }
                body={ {} }
                onData={ (data) => this.setState(data) }
            />
        }

        return <div id="bootleg">
            <div className="title">
                <div className="parent">
                    <Link to="/bootlegs">Bootleg</Link>
                </div>
                <h1>
                    { this.state.name }
                </h1>
            </div>

            <section className="info">
                <h2>Data</h2>
                <ul>
                    <li className="show">
                        <span>Show</span>
                        <span>
                            <Link to={ "/shows/" + this.state.showId }>
                                <img src={ "/static/" + imgTours[this.state.tour] } title={ this.state.labelTour } />
                                <img src={ "/static/" + imgCountries[this.state.labelCountryCode] } title={ this.state.labelCountry } />
                                <span>
                                    { this.state.showArtist }
                                    <br />{ this.state.showDate }
                                    <br />
                                    { this.state.showCity }, { this.state.showCountry }
                                </span>
                            </Link>
                        </span>
                    </li>
                    <li>
                        <span>Released</span>
                        <span>{ this.state.firstSeen }</span>
                    </li>
                    <li>
                        <span>Overall size</span>
                        <span>{ this.state.size }</span>
                    </li>
                    <li>
                        <span>Media type</span>
                        <span><img className="media-type-icon" src={ "/static/" + imgMediaTypes[this.state.type] } alt="" /> { this.state.typeLong }</span>
                    </li>
                    { (this.state.type != "misc") && <Fragment><li>
                        <span>Format</span>
                        <span>{ this.state.minfoFormat }</span>
                    </li>
                    <li>
                        <span>Duration</span>
                        <span>{ this.state.duration }</span>
                    </li></Fragment> }
                    { (this.state.type == "video") && <Fragment><li>
                        <span>Video Codec</span>
                        <span>{ this.state.minfoVideoCodec }</span>
                    </li>
                    <li>
                        <span>Video resolution</span>
                        <span>{ this.state.minfoVideoRes }</span>
                    </li></Fragment> }
                    { (this.state.type == "video" || this.state.type == "audio") && <Fragment><li>
                        <span>Audio Codec</span>
                        <span>{ this.state.minfoAudioCodec }</span>
                    </li>
                    <li>
                        <span>Audio resolution</span>
                        <span>{ this.state.minfoAudioRes }</span>
                    </li></Fragment> }
                </ul>
            </section>

            { (this.state.finfo != "") && <section className="finfo">
                <h2>Info</h2>
                <Scrollable>
                    { this.state.finfo }
                </Scrollable>
            </section> }

            <section className="files">
                <h2>Files</h2>
                <div className="note">
                    Download requires a Direct Connect client installed
                    (Windows: <a href="http://dcplusplus.sourceforge.net/" target="_blank" rel="noreferrer">DC++</a>,
                    Mac/Linux: <a href="https://sourceforge.net/projects/eiskaltdcpp/files/macOS/" target="_blank" rel="noreferrer">EiskaltDC++</a>)
                    and connected to the RadioHub (<a href="http://radiohub.wikidot.com/" target="_blank" rel="noreferrer">instructions</a>).
                </div>

                <div className="list">
                    { this.state.files.map((file) => {
                        return <a key={ file.TTH } href={ file.magnet }>
                            <span className="name">{ file.name }</span>
                            <span className="duration">{ file.duration }</span>
                            <span className="size">{ file.size }</span>
                            <span className="tth">{ file.TTH }</span>
                        </a>
                    }) }
                </div>
            </section>
        </div>
    }
}
