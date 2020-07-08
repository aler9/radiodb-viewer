
import React, { Component } from "react"
import { Link } from "react-router-dom"

import { Fetch } from "./fetch.jsx"
import { imgMediaTypes, imgTours, imgCountries } from "./imagearrays.jsx"

import "./show.scss"


export class Show extends Component {
    state = {}

    componentDidUpdate(prevProps, prevState) {
        if (this.state.loaded && !prevState.loaded) {
            document.title = "RadioDB - " + this.state.title
        }
    }

    render() {
        if (!this.state.loaded) {
            return <Fetch
                path={ "/data/show/" + this.props.match.params.id }
                body={ {} }
                onData={ (data) => this.setState(data) }
            />
        }

        return <div id="show">
            <div className="title">
                <div className="parent">
                    <Link to="/shows">Show</Link>
                </div>
                <h1>
                    { this.state.artistLong }, { this.state.date }, { this.state.city }, { this.state.labelCountry }
                </h1>
            </div>

            <section className="info">
                <h2>Data</h2>
                <ul>
                    <li>
                        <span>Artist</span>
                        <span>{ this.state.artistLong }</span>
                    </li>
                    <li>
                        <span>Date</span>
                        <span>{ this.state.date }</span>
                    </li>
                    <li className="tour">
                        <span>Tour</span>
                        <span><img src={ "/static/" + imgTours[this.state.tour] } /> { this.state.labelTour }</span>
                    </li>
                    <li className="country">
                        <span>Country</span>
                        <span><img src={ "/static/" + imgCountries[this.state.labelCountryCode] } /> { this.state.labelCountry }</span>
                    </li>
                    <li>
                        <span>Location</span>
                        <span>{ this.state.city }</span>
                    </li>
                    <li>
                        <span>Setlist</span>
                        <span>
                            { this.state.urls.map((url) => {
                                return <div key={ url.url }>
                                    <a href={ url.url } target="_blank" rel="noreferrer">view on { url.name }</a>
                                </div>
                            }) }
                        </span>
                    </li>
                </ul>
            </section>

            <section className="bootlegs">
                <h2>Bootlegs</h2>
                <div>
                    { this.state.bootlegs.map((bootleg) => {
                        return <Link key={ bootleg.id } to={ "/bootlegs/" + bootleg.id }>
                            <span className="name">
                                { bootleg.name }
                                <div className="released">Released: { bootleg.firstSeen }</div>
                            </span>
                            <span className="duration">
                                { bootleg.duration }
                            </span>
                            <span className="res" title={ bootleg.typeLong }>
                                <img className="media-type-icon" src={ "/static/" + imgMediaTypes[bootleg.type] } />
                                { bootleg.res }
                            </span>
                            <span className="size">{ bootleg.size }</span>
                        </Link>
                    }) }
                </div>
            </section>
        </div>
    }
}
