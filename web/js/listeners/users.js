import { DEFAULT_USERS, DEFAULT_RECENT_MESSAGES } from "./../services/user.js";
import { me } from "./../services/me.js";
import { getRoomId } from "./../services/messages.js";
import { navigate } from "./../router/router.js";

// Holds the conversation being opened: UserA is the authenticated user, UserB
// the other participant, RoomID the conversation uuid (null until the
// WebSocket reveals it — the users list does not carry it).
export let chat = { UserA: "", UserB: "", RoomID: null, RequestApi: false };

export function setChat(next) {
    chat = { ...next };
}

export function homeListenerUser() {
    const users = document.querySelectorAll(".user-item, .last-message-item");

    users.forEach((user) => {
        user.addEventListener("click", () => {
            const friendId = user.dataset.userId;
            if (!friendId) return;

            const friend = [...DEFAULT_RECENT_MESSAGES, ...DEFAULT_USERS].find(
                (candidate) => String(candidate.id?.Value ?? candidate.id) === String(friendId),
            );
            if (!friend) {
                console.warn("User not found in list:", friendId);
                return;
            }

            const roomId = getRoomId(friendId);
            chat = {
                UserA: me.ID?.Value ?? "",
                UserB: friendId,
                RoomID: roomId,
                RequestApi: roomId != null,
            };
            navigate("/chat");
        });
    });
}
