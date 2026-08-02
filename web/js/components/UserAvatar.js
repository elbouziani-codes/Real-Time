/**
 * UserAvatar
 * Reusable avatar block: letter/image + optional online status dot + optional badge.
 * Renders inside a `.user-avatar-wrapper` so status/badge are positioned correctly.
 */
export default function UserAvatar({ٍcolorClass = '', sizeClass = '', className = '', status = '', badge = ''} = {}) {
    const avatarClass = ['user-avatar', sizeClass, colorClass, className].filter(Boolean).join(' ');
    const content = image ? `<img src="${image}" alt="${letter || 'avatar'}">` : letter;

    const statusHtml = status ? `<span class="user-avatar__status ${status}"></span>` : '';
    const badgeHtml = badge ? `<span class="unread-badge">${badge}</span>` : '';

    return `
        <div class="user-avatar-wrapper">
            <div class="${avatarClass}">${content}</div>
            ${statusHtml}
            ${badgeHtml}
        </div>
    `;
}
