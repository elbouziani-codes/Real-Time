/** Optional `<h4>` title; the class attribute is omitted when no class is given. */
function renderTitle(title, titleClass) {
    if (!title) return '';

    const classAttribute = titleClass ? ` class="${titleClass}"` : '';
    return `<h4${classAttribute}>${title}</h4>`;
}

/**
 * SidebarSection
 * Generic titled block used by sidebar sections (Online, Offline, Recent Messages...).
 * Returns a fragment: an h4 title (optional) followed by the body markup.
 */
export default function SidebarSection({ title = '', titleClass = '', body = '' } = {}) {
    return `${renderTitle(title, titleClass)}${body}`;
}
