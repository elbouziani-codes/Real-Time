/**
 * UserAvatar
 * Letter-only avatar (the project has no user images) + optional
 * online/offline/away status dot + optional unread badge.
 *
 * Data comes from a User model; the optional second argument only
 * carries presentation classes (size / custom class), never data.
 */
export default function UserAvatar(user = {}, { sizeClass = '', className = '' } = {}) {
    const { letter = '', avatarClass: avatarColorClass = '', onlineStatus = '', unreadCount = 0 } = user;

    const avatarClass = ['user-avatar', sizeClass, avatarColorClass, className].filter(Boolean).join(' ');
    const statusHtml = onlineStatus ? `<span class="user-avatar__status ${onlineStatus}"></span>` : '';
    const badgeHtml = unreadCount ? `<span class="unread-badge">${unreadCount}</span>` : '';

    return `
        <div class="user-avatar-wrapper">
            <div class="${avatarClass}">${letter}</div>
            ${statusHtml}
            ${badgeHtml}
        </div>
    `;
}
