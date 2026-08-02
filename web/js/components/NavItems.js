const DEFAULT_ITEMS = [
    { label: '🏠 Home', active: true },
    { label: '💬 Chat' },
    { label: '➕ Create Post' },
];

/**
 * NavItems
 * Navigation buttons block of the navbar.
 * Data-driven: each item is { label, active }.
 */
export default function NavItems({ items = DEFAULT_ITEMS } = {}) {
    const itemsHtml = items.map(({ label = '', active = false }) => `<button${active ? ' class="active"' : ''}>${label}</button>`).join('');
    return `<div class="nav-links">${itemsHtml}</div>`;
}
