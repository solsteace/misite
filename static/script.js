const ATTR_THEME = "siteTheme"
document.documentElement.dataset[ATTR_THEME] = "light"

// Kudos: https://medium.com/@cerutti.alexander/a-mostly-complete-guide-to-theme-switching-in-css-and-js-c4992d5fd357
const ChangeTheme = function(theme = null) {
    let newTheme = theme
    if(!newTheme) {
        switch (document.documentElement.dataset[ATTR_THEME]) {
            case "dark": newTheme = "light"; break;
            case "light": newTheme = "dark"; break;
        }
    }
    document.documentElement.dataset[ATTR_THEME] = newTheme
}

const colorSchemePreference = window.matchMedia?.("(prefers-color-scheme:dark)")
if(colorSchemePreference) {
    if(colorSchemePreference.matches) {
        ChangeTheme("dark")
    }

    colorSchemePreference.addEventListener("change", e => {
        console.log("bjir")
    })
}

const FindAllArticleHeaders = function(targetElement) {
    const headerTags = ["h1", "h2", "h3", "h4", "h5", "h6"]
    const headers = [];
    targetElement.childNodes.forEach((child) => {
        let isHeader = false
        for(idx = 0; idx < headerTags.length; idx++) {
            if(child.nodeName.toLowerCase() == headerTags[idx]) {
                isHeader = true
                headers.push([idx, child])
                break;
            }
        }

        if(child.childNodes.length > 0 && !isHeader) {
            headers.push(...FindAllArticleHeaders(child))
            return
        }
    })
    return headers
}

const MakeOutline = function(articleElemId, outlineElemId) {
    const article = document.getElementById(articleElemId)
    const outline = document.getElementById(outlineElemId)

    if(!article) {
        console.log(`\`${articleElemId}\` element not found`)
        return
    } else if(!outline) {
        console.log(`\`${outlineElemId}\` element not found`)
        return
    }

    const headers = {}
    FindAllArticleHeaders(article).forEach((header, idx) => {
        const elem = header[1]

        const headerText = elem.innerHTML.toLowerCase().replaceAll(" ", "-")
        if(!headers[headerText]) {
            headers[headerText] = 0
        }
        headers[headerText] += 1
        elem.id = `${headerText}-${headers[headerText]}`

        outlineEntryElem = document.createElement("a")
        outlineEntryElem.innerHTML = elem.innerHTML
        outlineEntryElem.href= `#${elem.id}`
        outline.appendChild(outlineEntryElem)
    })
}
