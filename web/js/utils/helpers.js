
export function escapeHTML(value = "") {
    return String(value)
        .replaceAll("&", "&amp;")
        .replaceAll("<", "&lt;")
        .replaceAll(">", "&gt;")
        .replaceAll('"', "&quot;")
        .replaceAll("'", "&#39;");
}

export default function throttle(fn, wait) {
  let lastCall = 0;
  let timeout = null;

  return (...args) => {
    const now = Date.now();
    const remaining = wait - (now - lastCall);

    if (remaining <= 0) {
      // enough time has passed — fire immediately (leading edge)
      clearTimeout(timeout);
      timeout = null;
      lastCall = now;
      fn(...args);
    } else if (!timeout) {
      // schedule a trailing call so the last scroll near the edge isn't dropped
      timeout = setTimeout(() => {
        lastCall = Date.now();
        timeout = null;
        fn(...args);
      }, remaining);
    }
  };
}

