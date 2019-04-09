
import React, { Component } from "react"

import { DynamicReactRender } from "./various"
import { Scrollable } from "./scrollable"

class BootlegFinfo extends Component {
    render() {
        return <Scrollable>
            <div dangerouslySetInnerHTML={{ __html: this.props.text }} />
        </Scrollable>
    }
}

DynamicReactRender((dom) => <BootlegFinfo text={ dom.innerHTML } />, "bootleg-finfo")
