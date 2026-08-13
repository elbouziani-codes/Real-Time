
export default function UserAvatar(user = {}, { sizeClass = '', className = '' } = {}) {
    const { letter = '', avatarClass = '', onlineStatus = '', unreadCount = 0 } = user;

    const avatarClasses = ['user-avatar', sizeClass, avatarClass, className].filter(Boolean).join(' ');
    const statusHtml = onlineStatus ? `<span class="user-avatar__status ${onlineStatus}"></span>` : '';
    const badgeHtml = unreadCount ? `<span class="unread-badge">${unreadCount}</span>` : '';

    return `
        <div class="user-avatar-wrapper">
            <div class="${avatarClasses}">${letter}</div>
            ${statusHtml}
            ${badgeHtml}
        </div>
    `;
}
