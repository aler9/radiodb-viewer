
import React, { Component } from "react"

import { DynamicReactRender } from "./various.jsx"
import { Scrollable } from "./scrollable.jsx"

class BootlegFinfo extends Component {
    render() {
        return <Scrollable>
            <div dangerouslySetInnerHTML={{ __html: this.props.text }} />
        </Scrollable>
    }
}

DynamicReactRender((dom) => <BootlegFinfo text={ dom.innerHTML } />, "bootleg-finfo")
