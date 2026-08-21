import { DEFAULT_USERS, DEFAULT_RECENT_MESSAGES } from "./../services/user.js";
import { me } from "./../services/me.js";
import { navigate } from "./../router/router.js";

// Holds the conversation being opened: the authenticated user and the other
// participant. The server resolves their shared room.
export let chat = { UserA: "", UserB: "" };

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

            chat = {
                UserA: me.ID?.Value ?? "",
                UserB: friendId,
            };
            navigate("/chat");
        });
    });
}
