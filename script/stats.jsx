
import Chart from "chart.js"

window.addEventListener("pjax-load", () => {
    if(!document.getElementById("stats")) {
        return
    }

    let data = JSON.parse(unescape(document.getElementById("stats-peryear").getAttribute("data")))

    Chart.defaults.global.defaultFontColor = "rgb(230, 230, 230)"
    Chart.defaults.global.defaultFontFamily = "'Muli', sans-serif"
    window.myLine = Chart.Scatter(
        document.getElementById("stats-peryear").getContext("2d"), {
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
                tooltips: {
                    callbacks: {
                        label: (i, d) => {
                            if(i.datasetIndex == 2) {
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
})
