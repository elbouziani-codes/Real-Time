


// escapeHTML makes user written text safe to drop inside an innerHTML template,
// so a post or a comment cannot inject markup into the page.
export function escapeHTML(value = "") {
    return String(value)
        .replaceAll("&", "&amp;")
        .replaceAll("<", "&lt;")
        .replaceAll(">", "&gt;")
        .replaceAll('"', "&quot;")
        .replaceAll("'", "&#39;");
}

export default function throttle(callback, delay) {
    let lastTime = 0;

    return (...args) => {
        const now = Date.now();

        if (now - lastTime >= delay) {
            lastTime = now;
            callback(...args);
        }
    };
}