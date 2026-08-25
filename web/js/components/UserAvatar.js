import { me } from "../services/me.js";

export default function UserAvatar(user = {}, { sizeClass = '', className = '' } = {}) {
     if (me.NickName == user.NickName){
        user.online = true;
        user.letter = user.NickName[0];
    } 
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
