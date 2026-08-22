import { escapeHTML } from "./helpers.js";

// Transient notifications stacked in the top-right corner. Used for WebSocket
// events that concern conversations the user is not looking at.

let container = null;

function toastContainer() {
    if (!container || !document.body.contains(container)) {
        container = document.createElement("div");
        container.className = "toast-container";
        document.body.appendChild(container);
    }
    return container;
}

// Shows one toast. onClick makes it clickable (e.g. jump to the conversation);
// toasts dismiss themselves after timeout ms.
export function showToast({ title = "", message = "", timeout = 5000, onClick = null } = {}) {
    const toast = document.createElement("div");
    toast.className = "toast";
    toast.innerHTML =
        (title ? `<strong class="toast-title">${escapeHTML(title)}</strong>` : "") +
        (message ? `<span class="toast-message">${escapeHTML(message)}</span>` : "");

    let timer = null;
    const dismiss = () => {
        clearTimeout(timer);
        toast.classList.add("toast--leaving");
        setTimeout(() => toast.remove(), 200);
    };

    if (onClick) {
        toast.classList.add("toast--clickable");
        toast.addEventListener("click", () => {
            dismiss();
            onClick();
        });
    }

    timer = setTimeout(dismiss, timeout);
    toastContainer().appendChild(toast);
}
