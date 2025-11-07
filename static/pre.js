const findArticleHeaders = function(targetElement) {
    const headerTags = ["h2", "h3", "h4", "h5", "h6"]
    const headers = [];
    targetElement.childNodes.forEach((child) => {
        let isHeader = false
        for(idx = 0; idx < headerTags.length; idx++) {
            if(child.nodeName.toLowerCase() == headerTags[idx]) {
                isHeader = true
                headers.push([idx + 1, child])
                break;
            }
        }

        if(child.childNodes.length > 0 && !isHeader) {
            headers.push(...findArticleHeaders(child))
            return
        }
    })
    return headers
}

const makeArticleOutline = function(articleElemId, outlineElemId) {
    const article = document.getElementById(articleElemId)
    const outline = document.getElementById(outlineElemId)

    if(!article) {
        console.log(`\`${articleElemId}\` element not found`)
        return
    } else if(!outline) {
        console.log(`\`${outlineElemId}\` element not found`)
        return
    }

    let topMostVisibleEl
    const visibleEl = {}
    const observerOpts = {threshold: 1.0}
    const headers = findArticleHeaders(article)
    const observer = new IntersectionObserver(entries => {
        entries.forEach(entry => {
            const entryEl = entry.target
            if(entry.isIntersecting) {
                visibleEl[entryEl.id] = entryEl.dataset["outlineIdx"]
                return
            }

            // When scrolling up, if there's no intersecting header within
            // viewport, use the previous header of the last top most header
            if(topMostVisibleEl && Object.keys(visibleEl).length == 1) {
                const outlineIdx = topMostVisibleEl.dataset["outlineIdx"] 
                const topMostVisibleElPos = topMostVisibleEl.getBoundingClientRect()
                if(outlineIdx != "0" && topMostVisibleElPos.y > 0) {
                    const [_, previousSibling] = headers[outlineIdx - 1]
                    visibleEl[previousSibling.id] = previousSibling.dataset["outlineIdx"]
                }
            }
            delete visibleEl[entryEl.id]
        })

        Object.values(visibleEl).forEach(elIdx => {
            const [_, el] = headers[elIdx]
            if(!topMostVisibleEl) {
                topMostVisibleEl = el
                document
                    .querySelector(`#outline-item${topMostVisibleEl.dataset["outlineIdx"]}`)
                    .classList.add("specification__outline--active")
                return
            }

            const elPos = el.getBoundingClientRect()
            const topMostVisibleElPos = topMostVisibleEl.getBoundingClientRect()
            const shouldUpdate = 
                !visibleEl[topMostVisibleEl.id] // When top most not within viewport
                || elPos.y < topMostVisibleElPos.y // if there's another element, choose the upper one
                || topMostVisibleElPos.y < 0 // when top most left the view port from top direction
            if(shouldUpdate) {
                document
                    .querySelector(`#outline-item${topMostVisibleEl.dataset["outlineIdx"]}`)
                    .classList.remove("specification__outline--active")
                topMostVisibleEl = el
                document
                    .querySelector(`#outline-item${topMostVisibleEl.dataset["outlineIdx"]}`)
                    .classList.add("specification__outline--active")
            }
        })
    }, observerOpts)

    let currentLevel = 0 // Bypass H1, which usually be the title
    let nextUlEl
    const root = document.createElement("ul")
    const headerStack = [root]
    headers.forEach((header, idx) => {
        const [level, hEl] = header
        observer.observe(hEl)

        const headingText = hEl.innerHTML.toLowerCase().replaceAll(" ", "-")
        const liEl = document.createElement("li")
        const aEl = document.createElement("a")
        aEl.innerHTML = hEl.innerHTML
        aEl.id = `outline-item${idx}`
        aEl.href= `#${headingText}-${idx}`
        hEl.id = `${headingText}-${idx}`
        hEl.dataset["outlineIdx"] = idx
        liEl.appendChild(aEl)

        if(level > currentLevel) {
            if(nextUlEl) {
                headerStack.push(nextUlEl)
            }
            nextUlEl = document.createElement("ul")
        } else if(level < currentLevel) {
            levelDifference = currentLevel - level
            for(let i = 0; i < levelDifference; i++) {
                headerStack.pop()
            }
        } 
        liEl.appendChild(nextUlEl) // Note: this will move the element from old parent node to the new one
        headerStack.at(-1).appendChild(liEl)
        currentLevel = level
    })
    outline.replaceChildren(root)
}

// Kudos: https://medium.com/@cerutti.alexander/a-mostly-complete-guide-to-theme-switching-in-css-and-js-c4992d5fd357
const applyTheme = function(theme = null) {
    const ATTR_THEME = "siteTheme"
    let newTheme = theme
    if(!theme) {
        switch (document.documentElement.dataset[ATTR_THEME]) {
            case "dark": newTheme = "light"; break;
            case "light": newTheme = "dark"; break;
            default: newTheme = "dark";
        }
    }
    localStorage.setItem("colorscheme", newTheme)
    document.documentElement.dataset[ATTR_THEME] = newTheme
}

const colorschemePreference = window.matchMedia?.("(prefers-color-scheme:dark)")
const setTheme = function(ignoreLocalStorage = false) {
    if(!ignoreLocalStorage) {
        applyTheme(localStorage.getItem("colorscheme"))
    } else if(colorschemePreference) {
        let theme = colorschemePreference.matches? "dark": "light"
        applyTheme(theme)
    }
}

setTheme()
colorschemePreference.addEventListener("change", e => {
    setTheme(true)
})