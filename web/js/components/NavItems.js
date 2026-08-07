



/** Single navigation button; only the active item carries the `active` class. */

/**
 * NavItems
 * Navigation buttons block of the navbar.
 * Data-driven: each item is { label, active }.
 */
export default function NavItems() {
    let buttonsHtml = `
    <button class ="active homePage">🏠 Home</button>
    <button class ="chatPage">💬 Chat</button>
    <button class ="createPost">➕ Create Post</button>
    `
    if (window.location.pathname == "/chat"){
        buttonsHtml = `
            <button class ="homePage">🏠 Home</button>
            <button class ="active chatPage">💬 Chat</button>
    `
    console.log(buttonsHtml)
    }

    return `<div class = "nav-links">${buttonsHtml}</div>`;
}
