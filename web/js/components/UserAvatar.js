export default function UserAvatar(user = {}, { sizeClass = '', className = '' } = {}) {
    const { letter = '', avatarClass = '', online = false } = user;

    const avatarClasses = ['user-avatar', sizeClass, avatarClass, className].filter(Boolean).join(' ');
    const statusClass = online ? 'online' : 'offline';
    const statusHtml = `<span class="user-avatar__status ${statusClass}"></span>`;
    return `
        <div class="user-avatar-wrapper">
            <div class="${avatarClasses}">${letter}</div>
            ${statusHtml}
        </div>
    `;
}
