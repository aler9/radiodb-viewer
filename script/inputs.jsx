
import React, { Component } from "react"
import { findDOMNode } from "react-dom"

import { Scrollable } from "./scrollable.jsx"
import { Touchable } from "./various.jsx"

export class Select extends Component {
    static defaultProps = {
        disabled: false,
    }

    state = {
        opened: false,
    }

    onClick = () => {
        if(this.props.disabled) {
            return
        }

        this.setState({
            opened: !this.state.opened,
        })
    }

    render() {
        return <div className={ [ "select", this.state.opened ? "opened" : "" ].join(" ") } onClick={ this.onClick }>
            { this.props.options[this.props.value] &&
                this.props.options[this.props.value] }
            { this.state.opened && <div className="options">
                { Object.keys(this.props.options).map((k) => {
                    return <div key={ k } onClick={ () => this.props.onChange(k) }>{ this.props.options[k] }</div>
                }) }
            </div> }
        </div>
    }
}

export class TextInput extends Component {
    onClear = () => {
        this.props.onChange("")
    }

    render() {
        return <div className={ [ "text-input", (this.props.value) ? "clearable" : "" ].join(" ") }>
            <input type="text" placeholder={ this.props.placeholder }
                value={ this.props.value } onChange={ (e) => this.props.onChange(e.target.value) } />
            { (this.props.value) && <div className="clear" onClick={ this.onClear }></div> }
        </div>
    }
}

export class RangeSelect extends Component {
    static defaultProps = {
        disabled: false,
    }

    state = {
        validRange: [ 0, 1 ],
        validValue: [ 0, 1 ],
    }

    startPosMouse = [ 0, 0 ]
    startPosKnob = [ 0, 0 ]

    componentDidMount() {
        this.updateState()
    }

    componentDidUpdate(prevProps) {
        this.updateState(prevProps)
    }

    updateState = (prevProps) => {
        if(!prevProps ||
            prevProps.range != this.props.range ||
            prevProps.value != this.props.value) {
            if(this.props.range.length == 2) {
                this.setState({ validRange: this.props.range })
                if(this.props.value.length != 2) {
                    this.setState({ validValue: this.props.range })
                }
            }
            if(this.props.value.length == 2) {
                this.setState({ validValue: this.props.value })
            }
        }
    }

    onStart = (e, i) => {
        if(this.props.disabled) {
            return
        }

        let domWidth = findDOMNode(this).offsetWidth
        this.startPosMouse[i] = e.pageX
        this.startPosKnob[i] = (this.state.validValue[i] - this.state.validRange[0]) / (this.state.validRange[1] - this.state.validRange[0]) * domWidth
    }

    onMove = (e, i) => {
        if(this.props.disabled) {
            return
        }

        // value to pos
        let pos = e.pageX - this.startPosMouse[i] + this.startPosKnob[i]

        // pos to value
        let domWidth = findDOMNode(this).offsetWidth
        let newval = [ ...this.state.validValue ]
        newval[i] = Math.round(pos / domWidth * (this.state.validRange[1] - this.state.validRange[0]) + this.state.validRange[0])

        // limits
        for(let n = 0; n < 2; n++) {
            if(newval[n] < this.state.validRange[0]) {
                newval[n] = this.state.validRange[0]
            } else if(newval[n] > this.state.validRange[1]) {
                newval[n] = this.state.validRange[1]
            }
        }
        if(i == 0 && newval[0] >= newval[1]) {
            newval[0] = newval[1] - 1
        }
        if(i == 1 && newval[1] <= newval[0]) {
            newval[1] = newval[0] + 1
        }

        if(newval[0] != this.state.validValue[0] || newval[1] != this.state.validValue[1]) {
            this.props.onChange(newval)
        }
    }

    render() {
        return <div className="range-select">
            <div className="min">{ this.state.validValue[0] }</div>
            <div className="max">{ this.state.validValue[1] }</div>
            <div className="slider">
                <div className="range" style={ {
                    left: ((this.state.validValue[0] - this.state.validRange[0]) / (this.state.validRange[1] - this.state.validRange[0]) * 100) + "%",
                    right: ((1 - (this.state.validValue[1] - this.state.validRange[0]) / (this.state.validRange[1] - this.state.validRange[0])) * 100) + "%",
                } }>
                    <Touchable className="knob-left"
                        onStart={ (e) => this.onStart(e, 0) }
                        onMove={ (e) => this.onMove(e, 0) } />
                    <Touchable className="knob-right"
                        onStart={ (e) => this.onStart(e, 1) }
                        onMove={ (e) => this.onMove(e, 1) } />
                </div>
            </div>
        </div>
    }
}

export class MultiChoice extends Component {
    static defaultProps = {
        sort: "asc",
        disabled: false,
    }

    onChange = (k) => {
        if(this.props.disabled) {
            return
        }

        let checked = this.props.value[k] ? false : true
        if(checked) {
            this.props.onChange({
                ...this.props.value,
                [k]: 1,
            })
        } else {
            let value = { ...this.props.value }
            delete value[k]
            this.props.onChange(value)
        }
    }

    render() {
        return <Scrollable className="multichoices">
            <ul>
                { Object.keys(this.props.options).map((k) => {
                    return <li key={ k }
                        className={ this.props.value[k] ? "checked" : "" }
                        onClick={ () => this.onChange(k) }>
                        <div className="checkbox" />
                        { this.props.options[k] }
                    </li>
                }) }
            </ul>
        </Scrollable>
    }
}
