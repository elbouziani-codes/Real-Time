
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

