
import React, { Component } from "react"
import ChartJs from "chart.js"

import { Fetch } from "./fetch.jsx"

import "./stats.scss"


class Chart extends Component {
    componentDidMount() {
        const el = document.getElementById("stats-peryear")
        const data = this.props.data

        ChartJs.defaults.global.defaultFontColor = "rgb(230, 230, 230)"
        ChartJs.defaults.global.defaultFontFamily = "'Muli', sans-serif"
        window.myLine = ChartJs.Scatter(
            el.getContext("2d"), {
                data: {
                    datasets: [
                        {
                            showLine: true,
                            yAxisID: "A",
                            label: "Shows",
                            fill: false,
                            borderColor: "rgba(220, 220, 220, 0.7)",
                            data: Object.keys(data[0]).map((key) => {
                                return { x: key, y: data[0][key] }
                            })
                        },
                        {
                            showLine: true,
                            yAxisID: "A",
                            label: "Bootlegs",
                            fill: false,
                            borderColor: "rgba(0, 150, 0, 0.7)",
                            data: Object.keys(data[1]).map((key) => {
                                return { x: key, y: data[1][key] }
                            })
                        },
                        {
                            showLine: true,
                            yAxisID: "B",
                            label: "Bootleg size",
                            fill: false,
                            borderColor: "rgba(220, 220, 0, 0.7)",
                            data: Object.keys(data[2]).map((key) => {
                                return { x: key, y: data[2][key] }
                            }),
                        },
                    ],
                },
                options: {
                    animation: false,
                    elements: {
                        line: {
                            tension: 0,
                        },
                    },
                    tooltips: {
                        callbacks: {
                            label: (i, d) => {
                                if (i.datasetIndex == 2) {
                                    return i.xLabel + ", " + Math.round(i.yLabel / 1024 / 1024 / 1024) + " GB"
                                }
                                return i.xLabel + ", " + i.yLabel
                            },
                        }
                    },
                    aspectRatio: 3,
                    scales: {
                        yAxes: [
                            {
                                id: "A",
                                type: "linear",
                                position: "left",
                            },
                            {
                                id: "B",
                                type: "linear",
                                position: "right",
                                ticks: {
                                    callback: (v) => Math.round(v / 1024 / 1024 / 1024) + " GB",
                                },
                            }
                        ]
                    }
                }
            }
        )
    }

    render() {
        return <canvas id="stats-peryear"></canvas>
    }
}

export class Stats extends Component {
    state = {}

    componentDidMount() {
        document.title = "RadioDB - Statistics"
    }

    render() {
        if (!this.state.loaded) {
            return <Fetch
                path="/data/stats"
                body={ {} }
                onData={ (data) => this.setState(data) }
            />
        }

        return <div id="stats">
            <section>
                <h2>By year</h2>
                <Chart data={ this.state.perYear } />
            </section>

            <section>
                <h2>General</h2>
                <ul>
                    <li>
                        <span>Date of last update</span>
                        <span>{ this.state.generated }</span>
                    </li>
                </ul>
            </section>

            <section>
                <h2>Files</h2>
                <ul>
                    <li>
                        <span>File count</span>
                        <span>{ this.state.stats.FileCount }</span>
                    </li>
                    <li>
                        <span>Unique files</span>
                        <span>{ this.state.stats.FileUniqueCount }</span>
                    </li>
                    <li>
                        <span>Overall size</span>
                        <span>{ this.state.shareSize }</span>
                    </li>
                    <li>
                        <span>Unique size</span>
                        <span>{ this.state.shareUniqueSize }</span>
                    </li>
                </ul>
            </section>

            <section>
                <h2>Bootlegs</h2>
                <ul>
                    <li>
                        <span>Bootleg count</span>
                        <span>{ this.state.stats.BootlegCount }</span>
                    </li>
                    <li>
                        <span>Audio bootlegs</span>
                        <span>{ this.state.stats.AudioCount }</span>
                    </li>
                    <li>
                        <span>Video bootlegs</span>
                        <span>{ this.state.stats.VideoCount }</span>
                    </li>
                    <li>
                        <span>Misc bootlegs</span>
                        <span>{ this.state.stats.MiscCount }</span>
                    </li>
                    <li>
                        <span>Date of latest bootleg</span>
                        <span>{ this.state.dateLastBootleg }</span>
                    </li>
                </ul>
            </section>

            <section>
                <h2>Shows</h2>
                <ul>
                    <li>
                        <span>Show count</span>
                        <span>{ this.state.stats.ShowCount }</span>
                    </li>
                    <li>
                        <span>Earliest show</span>
                        <span>{ this.state.dateEarliest }</span>
                    </li>
                    <li>
                        <span>Latest show</span>
                        <span>{ this.state.dateLatest }</span>
                    </li>
                </ul>
            </section>
        </div>
    }
}
