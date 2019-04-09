
import React, { Component } from "react"
import { render, unmountComponentAtNode, findDOMNode } from "react-dom"

export const isIE = (navigator.userAgent.indexOf("MSIE") !== -1 || navigator.appVersion.indexOf("Trident/") > -1)
const supportsTouch = (window.ontouchstart !== undefined) || navigator.msMaxTouchPoints

export const urlParamsEncode = (params) => {
    return Object.keys(params).map((key) => {
        let val = params[key]
        if(val instanceof Array) {
            val = val.join(",")
        }
        return key + "=" + encodeURIComponent(val)
    }).join("&")
}

export const urlParamsDecode = (types) => {
    if(window.location.search.length == 0) {
        return {}
    }

    return window.location.search.substr(1).split("&").reduce((ret, e) => {
        let kv = e.split("=")
        let [ key, val ] = [ kv[0], decodeURIComponent(kv[1]) ]
        if(key in types) {
            switch(types[key]) {
            case "string": ret[key] = val; break
            case "number": ret[key] = parseInt(val); break
            case "string_array":
                ret[key] = val.split(",").reduce((ret, e) => {
                    if(e != "") {
                        ret.push(e)
                    }
                    return ret
                }, [])
                break
            case "number_array":
                ret[key] = val.split(",").reduce((ret, e) => {
                    if(e != "") {
                        ret.push(parseInt(e))
                    }
                    return ret
                }, [])
                break
            default: throw Error("unknown type")
            }
        }
        return ret
    }, {})
}

export const debounce = (func, delay) => {
    let timer = null
    return function() { // arrow functions does not expose arguments
        clearTimeout(timer)
        let savedCtx = this
        let savedArgs = arguments
        timer = setTimeout(() => func.apply(savedCtx, savedArgs), delay)
    }
}

export const DynamicReactRender = (cb, elId) => {
    let lastMounted = null
    window.addEventListener("pjax-load", () => {
        let dom = document.getElementById(elId)
        if(dom) {
            lastMounted = dom
            render(cb(dom), dom)
        }
    })
    window.addEventListener("pjax-unload", () => {
        if(lastMounted) {
            unmountComponentAtNode(lastMounted)
            lastMounted = null
        }
    })
}

class TouchableMouse extends Component {
    static defaultProps = {
        tag: "div",
        onStart: null,
        onMove: null,
        onEnd: null,
    }

    onMouseDown = (e) => {
        e.preventDefault()
        e.stopPropagation()
        if(this.props.onStart) {
            this.props.onStart({
                pageX: e.pageX,
                pageY: e.pageY,
            })
        }
        if(this.props.onMove || this.props.onEnd) {
            window.addEventListener("mousemove", this.onMouseMove)
            window.addEventListener("mouseup", this.onMouseUp)
        }
    }

    onMouseMove = (e) => {
        e.stopPropagation()
        if(this.props.onMove) {
            this.props.onMove({
                pageX: e.pageX,
                pageY: e.pageY,
            })
        }
    }

    onMouseUp = (e) => {
        e.stopPropagation()
        if(this.props.onEnd) {
            this.props.onEnd({ pageX: e.pageX, pageY: e.pageY })
        }
        window.removeEventListener("mousemove", this.onMouseMove)
        window.removeEventListener("mouseup", this.onMouseUp)
    }

    render() {
        let Tag = this.props.tag
        let props = { ...this.props }
        delete props["tag"]
        delete props["onStart"]
        delete props["onMove"]
        delete props["onEnd"]
        return <Tag { ...props } onMouseDown={ this.onMouseDown } />
    }
}

class TouchableTouch extends Component {
    static defaultProps = {
        tag: "div",
        onStart: null,
        onMove: null,
        onEnd: null,
    }

    componentDidMount() {
        // workaround: preventDefault() can not be called from TouchStart() unless binding is manual
        // https://github.com/facebook/react/issues/9809
        let dom = findDOMNode(this)
        dom.addEventListener("touchstart", this.onTouchStart)
        dom.addEventListener("touchmove", this.onTouchMove)
        dom.addEventListener("touchend", this.onTouchEnd)
    }

    onTouchStart = (e) => {
        e.stopPropagation()
        e.preventDefault() // prevent scroll
        if(this.props.onStart) {
            this.props.onStart({ pageX: e.changedTouches[0].pageX, pageY: e.changedTouches[0].pageY })
        }
    }

    onTouchMove = (e) => {
        e.stopPropagation()
        if(this.props.onMove) {
            this.props.onMove({ pageX: e.changedTouches[0].pageX, pageY: e.changedTouches[0].pageY })
        }
    }

    onTouchEnd = (e) => {
        e.stopPropagation()
        if(this.props.onEnd) {
            this.props.onEnd({ pageX: e.changedTouches[0].pageX, pageY: e.changedTouches[0].pageY })
        }
    }

    render() {
        const Tag = this.props.tag
        let props = { ...this.props }
        delete props["tag"]
        delete props["onStart"]
        delete props["onMove"]
        delete props["onEnd"]
        return <Tag { ...props } />
    }
}

export const Touchable = supportsTouch ? TouchableTouch : TouchableMouse
