/**
 * SidebarSection
 * Generic titled block used by sidebar sections (Online, Offline, Recent Messages...).
 * Returns a fragment: an h4 title (optional) followed by the body markup.
 */
export default function SidebarSection({ title = '', titleClass = '', body = '' } = {}) {
    const titleHtml = title ? `<h4${titleClass ? ` class="${titleClass}"` : ''}>${title}</h4>` : '';

    return `${titleHtml}${body}`;
}
