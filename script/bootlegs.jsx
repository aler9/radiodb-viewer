
import React, { Component } from "react"

import { DynamicReactRender, urlParamsDecode } from "./various"
import { InfiniteList } from "./infinitelist"
import { MultiChoice, Select, TextInput } from "./inputs"

class Bootlegs extends Component {
    static paramTypes = {
        sort: "string",
        text: "string",
        media: "string_array",
        audioRes: "string_array",
        videoRes: "string_array",
    }

    static initParams ={
        sort: "date",
        text: "",
        media: [],
        audioRes: [],
        videoRes: [],
    }

    state = {
        choices: {
            sort: {
                date: "Release date",
                sdate_desc: "Show date, descending",
                sdate_asc: "Show date, ascending",
                size_desc: "Size, descending",
                size_asc: "Size, ascending",
            },
            media: {},
            audioRes: {},
            videoRes: {},
        },
        params: { ...this.constructor.initParams, ...urlParamsDecode(this.constructor.paramTypes) },
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
        return <div>
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
                    <div className="button-clear" onClick={ this.onReset }>Reset filters</div>
                </span>
                <span className="text">
                    <h3>Search</h3>
                    <TextInput
                        placeholder={ "Enter text or TTH" }
                        value={ this.state.params.text }
                        onChange={ (text) => this.onParams({ text }) }
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
                <span>
                    <h3>Audio resolution</h3>
                    <MultiChoice
                        options={ this.state.choices.audioRes }
                        value={ this.state.params.audioRes.reduce((ret, e) => { ret[e] = 1; return ret }, {}) }
                        onChange={ (audioRes) => this.onParams({ audioRes: Object.keys(audioRes) }) }
                    />
                </span>
                <span>
                    <h3>Video resolution</h3>
                    <MultiChoice
                        options={ this.state.choices.videoRes }
                        value={ this.state.params.videoRes.reduce((ret, e) => { ret[e] = 1; return ret }, {}) }
                        onChange={ (videoRes) => this.onParams({ videoRes: Object.keys(videoRes) }) }
                    />
                </span>
            </div>
            <InfiniteList
                url={ "/data/bootlegs" }
                params={ this.state.params }
                onChoices={ this.onChoices } />
        </div>
    }
}

DynamicReactRender(() => <Bootlegs />, "bootlegs")
