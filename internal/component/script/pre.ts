import {codeToHtml } from "./generated/shiki.bundle";

function findArticleHeaders(targetElement: Element) {
    const headerTags = ["h2", "h3", "h4", "h5", "h6"] as const;
    const headers: Array<[number, HTMLElement]> = [];

    targetElement.childNodes.forEach((child) => {
        if (child.nodeType !== Node.ELEMENT_NODE) return;

        const el = child as Element;

        const tagIndex = headerTags.indexOf(el.tagName.toLowerCase() as typeof headerTags[number]);
        const isHeader = tagIndex !== -1;
        if (isHeader) {
            headers.push([tagIndex + 1, <HTMLElement>el]);
        } else if (el.childNodes.length > 0) {
            headers.push(...findArticleHeaders(el));
        }
    });

    return headers;
}

window._utilMakeOutline = function(
    articleElemId: string, 
    outlineElemId: string
) {
    const article = document.getElementById(articleElemId)
    const outline = document.getElementById(outlineElemId)

    if(!article) {
        console.log(`\`${articleElemId}\` element not found`)
        return
    } else if(!outline) {
        console.log(`\`${outlineElemId}\` element not found`)
        return
    }

    let topMostVisibleEl: HTMLElement | undefined
    const visibleEl: {[key: string]: string} = {}
    const observerOpts = {threshold: 1.0}
    const headers = findArticleHeaders(article)
    const observer = new IntersectionObserver(entries => {
        entries.forEach(entry => {
            const entryEl = <HTMLElement> entry.target
            if(entry.isIntersecting) {
                visibleEl[entryEl.id] = entryEl.dataset["outlineIdx"]!
                return
            }

            // When scrolling up, if there's no intersecting header within
            // viewport, use the previous header of the last top most header
            if(topMostVisibleEl && Object.keys(visibleEl).length == 1) {
                const outlineIdx = Number(topMostVisibleEl.dataset["outlineIdx"])
                const topMostVisibleElPos = topMostVisibleEl.getBoundingClientRect()
                if(outlineIdx != 0 && topMostVisibleElPos.y > 0) {
                    const previousSibling = headers[outlineIdx - 1]![1]
                    visibleEl[previousSibling.id] = previousSibling.dataset["outlineIdx"]!
                }
            }
            delete visibleEl[entryEl.id]
        })

        Object.values(visibleEl).forEach(elIdx => {
            const el = headers[Number(elIdx)]![1]
            if(!topMostVisibleEl) {
                topMostVisibleEl = el
                document
                    .querySelector(`#outline-item${topMostVisibleEl.dataset["outlineIdx"]}`)!
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
                    .querySelector(`#outline-item${topMostVisibleEl.dataset["outlineIdx"]}`)!
                    .classList.remove("specification__outline--active")
                topMostVisibleEl = el
                document
                    .querySelector(`#outline-item${topMostVisibleEl.dataset["outlineIdx"]}`)!
                    .classList.add("specification__outline--active")
            }
        })
    }, observerOpts)

    const root = document.createElement("ul")
    const headerStack = [root]
    let currentLevel = 0 // Bypass H1, which usually be the title
    let nextUlEl: typeof root
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
        hEl.dataset["outlineIdx"] = `${idx}`
        liEl.appendChild(aEl)

        if(level > currentLevel) {
            if(nextUlEl) {
                headerStack.push(nextUlEl)
            }
            nextUlEl = document.createElement("ul")
        } else if(level < currentLevel) {
            const levelDifference = currentLevel - level
            for(let i = 0; i < levelDifference; i++) {
                headerStack.pop()
            }
        } 
        liEl.appendChild(nextUlEl) // Note: this will move the element from old parent node to the new one
        headerStack.at(-1)!.appendChild(liEl)
        currentLevel = level
    })
    outline.replaceChildren(root)
}

// Kudos: https://medium.com/@cerutti.alexander/a-mostly-complete-guide-to-theme-switching-in-css-and-js-c4992d5fd357
const applyTheme = function(theme?: string) {
    const ATTR_THEME = "siteTheme"
    let newTheme = theme ?? "dark"
    if(!theme) {
        switch (document.documentElement.dataset[ATTR_THEME]) {
            case "dark": newTheme = "light"; break;
            case "light": newTheme = "dark"; break;
        }
    }
    localStorage.setItem("colorscheme", newTheme)
    document.documentElement.dataset[ATTR_THEME] = newTheme
}

const colorschemePreference = window.matchMedia?.("(prefers-color-scheme:dark)")
const setTheme = function(ignoreLocalStorage: boolean = false) {
    if(!ignoreLocalStorage) {
        applyTheme(localStorage.getItem("colorscheme") ?? undefined)
    } else if(colorschemePreference) {
        let theme = colorschemePreference.matches? "dark": "light"
        applyTheme(theme)
    }
}

window._utilHighlightCode = async() => {
    const codeblocks = document.getElementsByClassName("codeblock__code")
    for (const cb of codeblocks) {
        const codes = cb.getElementsByTagName("code")
        const canSkip = (
            codes.length < 1 
            || codes[0]!.parentElement!.classList.contains("shiki"))
        if(canSkip) {
            return
        }

        cb.parentElement!.innerHTML = await codeToHtml(
            codes[0]!.textContent, 
            {
                lang: cb.parentElement?.dataset.lang ?? "",
                themes: {
                    light: 'github-light',
                    dark: 'catppuccin-mocha' 
                },
                transformers: [{
                    pre(node) {
                        this.addClassToHast(node, "codeblock__code")
                    }
                }]
        })
        .catch(err => { // fallback to default styling
            console.error(err)
            return cb.parentElement!.innerHTML
        })
    }
}

{(function() {
    setTheme()
    document.addEventListener("DOMContentLoaded", (e) => {
        document.getElementById("site__theme")
            ?.addEventListener("click", (e) => { 
                applyTheme() 
            })
    }, {once: true})
    colorschemePreference.addEventListener("change", e => {
        setTheme(true)
    })
})() }