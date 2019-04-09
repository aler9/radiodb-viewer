
import React, { Component } from "react"
import { findDOMNode } from "react-dom"

import { debounce, urlParamsEncode } from "./various"

export class InfiniteList extends Component {
    state = {
        state: "loading", // loading, waitscroll, fullyloaded, noresults
        items: [],
    }

    curPage = 0

    curRequestId = 0 // id to block concurrent requests

    componentDidMount() {
        this.fetch(++this.curRequestId, 0)
        window.addEventListener("scroll", this.onScrollOrResize)
        window.addEventListener("resize", this.onScrollOrResize)
    }

    componentWillUnmount() {
        window.removeEventListener("scroll", this.onScrollOrResize)
        window.removeEventListener("resize", this.onScrollOrResize)
    }

    componentDidUpdate(prevProps) {
        // params have changed
        if(this.props.params != prevProps.params) {
            this.setState({
                state: "loading",
            })
            this.onItems([], false) // reset items
            this.curPage = 0
            this.fetchAndPushDebounced(++this.curRequestId)
        }
    }

    onScrollOrResize = () => {
        let n = findDOMNode(this)
        if(this.state.state == "waitscroll" &&
            n.offsetParent != null &&
            (n.getBoundingClientRect().bottom - 500) <= window.innerHeight) {

            // go to next page
            this.setState({
                state: "loading",
            })
            this.fetch(++this.curRequestId, ++this.curPage)
        }
    }

    fetchAndPushDebounced = debounce((requestId) => {
        // push params
        history.replaceState(history.state, null, window.location.pathname + "?"
            + urlParamsEncode(this.props.params))

        this.fetch(requestId, 0)
    }, 400)

    fetch = (requestId, fetchPage) => {
        fetch(this.props.url, {
            method: "POST",
            body: JSON.stringify({
                curPage: fetchPage,
                ...this.props.params,
            })
        })
            .then((r) => r.json())
            .then((data) => {
                if(requestId != this.curRequestId) {
                    return
                }

                this.setState({
                    state: (() => {
                        if(!data.items) {
                            return "noresults"
                        }
                        if(data.fullyLoaded) {
                            return "fullyloaded"
                        }
                        return "waitscroll"
                    })(),
                })

                if(fetchPage == 0) {
                    this.props.onChoices(data.choices)
                }
                this.onItems(data.items ? data.items : [], fetchPage > 0)

                // check if next page must be loaded immediately (if screen is too big)
                setTimeout(() => this.onScrollOrResize(), 0)
            })
    }

    onItems = (items, append) => {
        if(append) {
            items = [ ...this.state.items, ...items ]
        }
        this.setState({ items })
    }

    render() {
        let clx = [ "infinitelist" ]
        if(this.state.state == "noresults") {
            clx.push("noresults")
        } else if(this.state.state != "fullyloaded") {
            clx.push("loading")
        }

        return <div className={ clx.join(" ") }><div>
            { this.state.items.map((s) => {
                return <a key={ s.props.key } { ...s.props } dangerouslySetInnerHTML={{ __html: s.cnt }} />
            }) }
        </div></div>
    }
}
