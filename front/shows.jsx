
import React, { Component } from "react"
import { Link } from "react-router-dom"

import { urlParamsDecode } from "./various.jsx"
import { InfiniteList } from "./infinitelist.jsx"
import { Select, MultiChoice, RangeSelect, TextInput } from "./inputs.jsx"
import { imgTours, imgCountries } from "./imagearrays.jsx"

import "./shows.scss"

import imgAudio from "./images/audio.svg"
import imgVideo from "./images/video.svg"
import imgMisc from "./images/misc.svg"


export class Shows extends Component {
    static paramTypes = {
        sort: "string",
        text: "string",
        artist: "string_array",
        tour: "string_array",
        year: "number_array",
        country: "string_array",
        media: "string_array",
    }

    static initParams = {
        sort: "date_desc",
        text: "",
        artist: [],
        tour: [],
        year: [],
        country: [],
        media: [],
    }

    state = {
        choices: {
            sort: {
                date_desc: "Date, descending",
                date_asc: "Date, ascending",
            },
            artist: {},
            tour: {},
            year: [],
            country: {},
            media: {},
        },
        params: { ...this.constructor.initParams, ...urlParamsDecode(this.constructor.paramTypes) },
    }

    componentDidMount() {
        document.title = "RadioDB - Shows"
    }

    onReset = () => {
        let params = { ...this.constructor.initParams }
        delete params["sort"]
        this.onParams(params)
    }

    onParams = (params) => {
        this.setState({ params: { ...this.state.params, ...params } })
    }

    onChoices = (choices) => {
        this.setState({ choices: { ...this.state.choices, ...choices } })
    }

    render() {
        return <div id="shows"><div>
            <div className="controls">
                <span className="sort">
                    <h3>Sort by</h3>
                    <Select
                        options={ this.state.choices.sort }
                        value={ this.state.params.sort }
                        onChange={ (sort) => this.onParams({ sort }) }
                    />
                </span>
                <span>
                    <div className="button-clear" onClick={ this.onReset }>
                        Reset filters
                    </div>
                </span>
                <span className="text">
                    <h3>Search</h3>
                    <TextInput
                        placeholder={ "Enter text" }
                        value={ this.state.params.text }
                        onChange={ (text) => this.onParams({ text }) }
                    />
                </span>
                <span>
                    <h3>Artist</h3>
                    <MultiChoice
                        options={ this.state.choices.artist }
                        value={ this.state.params.artist.reduce((ret, e) => { ret[e] = 1; return ret }, {}) }
                        onChange={ (artist) => this.onParams({ artist: Object.keys(artist) }) }
                    />
                </span>
                <span>
                    <h3>Year</h3>
                    <RangeSelect
                        range={ this.state.choices.year }
                        value={ this.state.params.year }
                        onChange={ (year) => this.onParams({ year }) }
                    />
                </span>
                <span>
                    <h3>Tour</h3>
                    <MultiChoice
                        options={ this.state.choices.tour }
                        value={ this.state.params.tour.reduce((ret, e) => { ret[e] = 1; return ret }, {}) }
                        onChange={ (tour) => this.onParams({ tour: Object.keys(tour) }) }
                    />
                </span>
                <span>
                    <h3>Country</h3>
                    <MultiChoice
                        options={ this.state.choices.country }
                        value={ this.state.params.country.reduce((ret, e) => { ret[e] = 1; return ret }, {}) }
                        onChange={ (country) => this.onParams({ country: Object.keys(country) }) }
                    />
                </span>
                <span>
                    <h3>Media type</h3>
                    <MultiChoice
                        options={ this.state.choices.media }
                        value={ this.state.params.media.reduce((ret, e) => { ret[e] = 1; return ret }, {}) }
                        onChange={ (media) => this.onParams({ media: Object.keys(media) }) }
                    />
                </span>
            </div>
            <InfiniteList
                url={ "/data/shows" }
                params={ this.state.params }
                onChoices={ this.onChoices }
                onItem={ (item) => <Link key={ item.id } className="entry" to={ "/shows/" + item.id } title={ item.title }>
                    <span className="tour">
                        <img src={ "/static/" + imgTours[item.tour] } title={ item.labelTour } />
                    </span>

                    <span className="country">
                        <img src={ "/static/" + imgCountries[item.labelCountryCode] } title={ item.country } />
                    </span>

                    <span className="title">
                        <div>
                            { item.artistLong }
                        </div>
                        <h3>
                            { item.date }, { item.city }, { item.country }
                        </h3>
                    </span>

                    <span className="count">
                        { (item.audioCount > 0) && <span>
                            <img src={ "/static/" + imgAudio } />
                            { item.audioCount } { (item.audioCount != 1) ? "audios" : "audio" }
                        </span> }
                        { (item.videoCount > 0) && <span>
                            <img src={ "/static/" + imgVideo } />
                            { item.videoCount } { (item.videoCount != 1) ? "videos" : "video" }
                        </span> }
                        { (item.miscCount > 0) && <span>
                            <img src={ "/static/" + imgMisc } />
                            { item.miscCount } { (item.miscCount != 1) ? "miscs" : "misc" }
                        </span> }
                    </span>
                </Link> }
            />
        </div></div>
    }
}
