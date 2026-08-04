const DEFAULT_ITEMS = [
    { label: '🏠 Home', active: true },
    { label: '💬 Chat' },
    { label: '➕ Create Post' },
];

/** Single navigation button; only the active item carries the `active` class. */
function renderNavButton({ label = '', active = false }) {
    const activeAttribute = active ? ' class="active"' : '';
    return `<button${activeAttribute}>${label}</button>`;
}

/**
 * NavItems
 * Navigation buttons block of the navbar.
 * Data-driven: each item is { label, active }.
 */
export default function NavItems() {
    const buttonsHtml = DEFAULT_ITEMS.map(renderNavButton).join('');

    return `<div class="nav-links">${buttonsHtml}</div>`;
}
