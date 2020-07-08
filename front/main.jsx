
import "./reset.scss"
import "./main.scss"

import React, { Component } from "react"
import { render } from "react-dom"
import { BrowserRouter, NavLink, Route } from "react-router-dom"

import { Bootlegs } from "./bootlegs.jsx"
import { Bootleg } from "./bootleg.jsx"
import { Dump } from "./dump.jsx"
import { Front } from "./front.jsx"
import { Random } from "./random.jsx"
import { Shows } from "./shows.jsx"
import { Show } from "./show.jsx"
import { Stats } from "./stats.jsx"

import imgGithub from "./images/github.svg"
import imgLogo from "./images/logo.svg"
import imgMenu from "./images/menu.svg"


class Main extends Component {
    render() {
        return <div id="main">
            <header>
                { (this.props.location.pathname != "/") && <a className="logo" href="/">
                    <img src={ "/static/" + imgLogo } alt="RadioDB" />
                </a> }

                <input type="checkbox" id="menu-toggle" />
                <label htmlFor="menu-toggle">
                    <img src={ "/static/" + imgMenu } alt="" />
                </label>

                <nav><ul>
                    <li><NavLink exact to="/">Home</NavLink></li>
                    <li><NavLink exact to="/shows">shows</NavLink></li>
                    <li><NavLink exact to="/bootlegs">bootlegs</NavLink></li>
                    <li><NavLink exact to="/random">Random</NavLink></li>
                    <li><NavLink exact to="/stats">Stats</NavLink></li>
                    <li><NavLink exact to="/dump">Dump</NavLink></li>
                </ul></nav>
            </header>

            <main className={ (this.props.location.pathname == "/") ? "front" : "" }>
                <Route exact path="/" render={ (props) => <Front { ...props } /> } />
                <Route exact path="/shows" render={ (props) => <Shows { ...props } /> } />
                <Route exact path="/bootlegs" render={ (props) => <Bootlegs { ...props } /> } />
                <Route exact path="/shows/:id" render={ (props) => <Show { ...props } /> } />
                <Route exact path="/bootlegs/:id" render={ (props) => <Bootleg { ...props } /> } />
                <Route exact path="/random" render={ (props) => <Random { ...props } /> } />
                <Route exact path="/stats" render={ (props) => <Stats { ...props } /> } />
                <Route exact path="/dump" render={ (props) => <Dump { ...props } /> } />
            </main>

            <footer>
                <a className="source" href="https://github.com/aler9/radiodb-viewer" target="_blank" rel="noreferrer">
                    <img src={ "/static/" + imgGithub } />
                    source code
                </a>
                <br />
                <br />
                This site is not affiliated with any of the mentioned artists.
                <br />
                Use of low-resolution album covers qualifies as fair use under the copyright law of the United States.
                <br />
                This site does not host any file and is merely a tracker.
                <br />
                Audio and video materials mentioned in the site are exclusively non-commercial, non-official and freely available recordings.
                <br />
                If you want to have an entry removed, write to <a href="mailto:aler9.dev@gmail.com">aler9.dev@gmail.com</a>
            </footer>
        </div>
    }
}

render(<BrowserRouter>
    <Route component={ () => { window.scrollTo(0, 0); return null } } />
    <Route render={ (props) => <Main { ...props } /> } />
</BrowserRouter>, document.getElementById("root"))
