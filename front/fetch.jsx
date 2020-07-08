
import React, { Component } from "react"


const LOADING = 0
const ERROR = 1

export class Fetch extends Component {
    state = {
        state: LOADING
    }

    curReqId = 0

    componentDidMount() {
        this.load()
    }

    componentWillUnmount() {
        this.curReqId = -1
    }

    load = () => {
        const reqId = ++this.curReqId

        fetch(this.props.path, {
            method: "POST",
            body: JSON.stringify(this.props.body),
        })
            .then((res) => {
                if (reqId != this.curReqId) {
                    return
                }
                if (res.status != 200) {
                    return res.json().then((res) => { throw res.error })
                }
                return res.json()
            })
            .then((res) => this.props.onData({ ...res, loaded: true }))
            .catch(() => this.setState({ state: ERROR }))
    }

    render() {
        switch (this.state.state) {
        case LOADING:
            return null

        case ERROR:
            return <div className="error-net">network error</div>
        }
    }
}
