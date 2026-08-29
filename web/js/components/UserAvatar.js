//import { me } from "../services/me.js";

export default function UserAvatar(user) {
    const status = user.Online  ? "online" : "offline" 
    return `
        <div class="user-avatar-wrapper">
            <div class="user-avatar">${user.NickName[0]}</div>
			<span class="user-avatar__status ${status}"></span>
        </div>
    `;
}
