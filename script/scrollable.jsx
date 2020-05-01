
import React, { Component } from "react"
import { findDOMNode } from "react-dom"

import { Touchable } from "./various.jsx"

import "./scrollable.scss"


export class Scrollable extends Component {
    state = {
        scrollbarPos: 0,
        scrollbarWidth: 0,
        scrollbarHeight: 0,
    }

    touchStartMouse = 0
    touchStartPos = 0

    componentDidMount() {
        this.onResize()
        window.addEventListener("resize", this.onResize)
        this.contentDom = findDOMNode(this).childNodes[0]
        this.onScroll()
        this.contentDom.addEventListener("scroll", this.onScroll)
    }

    componentWillUnmount() {
        window.removeEventListener("resize", this.onResize)
    }

    componentDidUpdate() {
        this.onScroll()
    }

    onResize = () => {
        this.setState({
            scrollbarWidth: (() => {
                let el = document.createElement("div")
                el.style.width = "100px"
                el.style.height = "100px"
                el.style.overflow = "scroll"
                el.style.position = "absolute"
                el.style.top = "-9999px"
                document.body.appendChild(el)
                let ret = el.offsetWidth - el.clientWidth + 1
                document.body.removeChild(el)
                return ret
            })()
        })
    }

    onScroll = () => {
        let scrollbarPos = this.contentDom.scrollTop / this.contentDom.scrollHeight
        let scrollbarHeight = this.contentDom.offsetHeight / this.contentDom.scrollHeight

        // limit scrollbar top and bottom position for momentum scrolling (ios)
        if((scrollbarPos + scrollbarHeight) > 1) {
            if(scrollbarPos > 0.95) {
                scrollbarPos = 0.95
            }
            scrollbarHeight = 1 - scrollbarPos

        } else if(scrollbarPos < 0) {
            scrollbarHeight = scrollbarHeight + scrollbarPos
            scrollbarPos = 0
            if(scrollbarHeight < 0.05) {
                scrollbarHeight = 0.05
            }
        }

        if(!isNaN(scrollbarPos) && !isNaN(scrollbarHeight) &&
            (scrollbarPos != this.state.scrollbarPos ||
            scrollbarHeight != this.state.scrollbarHeight)) {
            this.setState({
                scrollbarPos,
                scrollbarHeight,
            })
        }
    }

    onScrollbarTouchStart = (evt) => {
        this.touchStartMouse = evt.pageY
        this.touchStartPos = this.contentDom.scrollTop
    }

    onScrollbarTouchMove = (evt) => {
        let newScroll = (evt.pageY - this.touchStartMouse) * this.contentDom.scrollHeight / this.contentDom.offsetHeight + this.touchStartPos
        this.contentDom.scrollTop = newScroll
    }

    render() {
        let clx = [ "scrollable" ]
        if(this.props.className) {
            clx.push(this.props.className)
        }

        return <div className={ clx.join(" ") }>
            <div className="content" style={{ marginRight: (- this.state.scrollbarWidth) + "px" }}>
                { this.props.children }
            </div>
            { (this.state.scrollbarHeight < 1) && <Touchable className="scrollbar" style={{
                top: (this.state.scrollbarPos * 100) + "%",
                height: (this.state.scrollbarHeight * 100) + "%",
            }}
            onStart={ this.onScrollbarTouchStart }
            onMove={ this.onScrollbarTouchMove }
            /> }
        </div>
    }
}
