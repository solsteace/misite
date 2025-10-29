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

    colorSchemePreference.addEventListener("change", e => { })
}

const FindAllArticleHeaders = function(targetElement) {
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

    let currentLevel = 0 // Bypass H1, which usually be the title
    let nextUlEl
    const root = document.createElement("ul")
    const headerStack = [root]
    const heading = {}
    FindAllArticleHeaders(article).forEach((header, _) => {
        // Generate id
        const [level, hEl] = header
        const headingText = hEl.innerHTML.toLowerCase().replaceAll(" ", "-")
        if(!heading[headingText]) {
            heading[headingText] = 0
        }
        heading[headingText] += 1

        // Create current entry
        const liEl = document.createElement("li")
        const aEl = document.createElement("a")
        const hId = `${headingText}-${heading[headingText]}`
        aEl.innerHTML = hEl.innerHTML
        aEl.href= `#${hId}`
        hEl.id = hId
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

    outline.appendChild(root)
}
