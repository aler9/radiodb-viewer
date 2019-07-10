
import { isIE, debounce } from "./various.jsx"

const loadCurPage = (scroll) => {
    window.dispatchEvent(new CustomEvent("pjax-unload"))

    // body: clear content and style
    document.body.innerHTML = ""
    document.body.setAttribute("style", "")

    let xhr = new XMLHttpRequest()
    xhr.onload = () => {
        // in case of 301/302, replace url
        if(xhr.responseURL != document.location.href) {
            history.replaceState(history.state, null, xhr.responseURL)
        }

        // head: replace title, do not touch anything else
        // https://github.com/luruke/barba.js/issues/82
        document.title = xhr.responseXML.title

        // body: replace children
        while(xhr.responseXML.body.hasChildNodes()) {
            document.body.appendChild(xhr.responseXML.body.childNodes[0])
        }

        window.dispatchEvent(new CustomEvent("pjax-load"))

        // for dynamic content: force a minimum height before restoring scroll
        // the check is to avoid breaking fullscreen pages
        if((scroll + window.innerHeight) > document.body.scrollHeight) {
            document.body.style.minHeight = (scroll + window.innerHeight) + "px"
        }

        // restore scroll
        window.scrollTo(0, scroll)
    }
    xhr.open("GET", document.location.href)
    xhr.responseType = "document"
    xhr.send()
}

const saveScroll = () => {
    history.replaceState({
        scroll: window.pageYOffset,
    }, null, window.location.href)
}

const onClick = (evt) => {
    evt.preventDefault()
    saveScroll() // save immediately
    history.pushState({ scroll: 0 }, null, evt.currentTarget.getAttribute("href"))
    loadCurPage(0)
}

const registerLink = (node) => {
    // skip nonexistent/external links
    let href = node.getAttribute("href")
    if(!href || href[0] != "/") {
        return
    }
    // skip special links
    if(node.getAttribute("nopjax") !== null) {
        return
    }
    node.addEventListener("click", onClick)
}

// wait for both js and css loading (main difference with DOMContentLoaded)
window.addEventListener("load", () => {
    // on IE XMLHttpRequest is too buggy and CustomEvent is not implemented
    if(isIE) {
        let evt = document.createEvent("CustomEvent")
        evt.initCustomEvent("pjax-load", false, false, undefined)
        window.dispatchEvent(evt)
        return
    }

    // popstate is called after current history entry has passed
    // so we must continuously save current scroll in state
    // (debounced for performance reasons)
    window.addEventListener("scroll", debounce(saveScroll, 500))

    window.addEventListener("popstate", (evt) => {
        loadCurPage(evt.state.scroll)
    })

    // register current links
    let nodes = document.getElementsByTagName("A")
    for(let i = 0; i < nodes.length; i++) {
        registerLink(nodes[i])
    }

    // register future links
    new MutationObserver((mutationList) => {
        mutationList.forEach((mutation) => {
            for(let i = 0; i < mutation.addedNodes.length; i++) {
                let node = mutation.addedNodes[i]
                if(node.nodeName == "A") {
                    registerLink(node)

                } else if(node.getElementsByTagName != undefined) {
                    let nodes = node.getElementsByTagName("A")
                    for(let j = 0; j < nodes.length; j++) {
                        registerLink(nodes[j])
                    }
                }
            }
        })
    }).observe(document.body, {
        childList: true,
        subtree: true,
    })

    window.dispatchEvent(new CustomEvent("pjax-load"))
})
