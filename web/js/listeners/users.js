import { DEFAULT_USERS, DEFAULT_RECENT_MESSAGES } from "./../services/user.js";
import { me } from "./../services/me.js";
import { roomIds } from "./../services/messages.js";
import { navigate } from "./../router/router.js";

// Holds the conversation being opened: UserA is the authenticated user, UserB
// the other participant, RoomID the conversation uuid (null until the
// WebSocket reveals it — the users list does not carry it).
export let chat = { UserA: "", UserB: "", RoomID: null, RequestApi: false };

export function setChat(next) {
    chat = { ...next };
}

export function homeListenerUser() {
    const users = document.querySelectorAll(".user-item");

    users.forEach((user) => {
        user.addEventListener("click", () => {
            // "item" rows render .user-item-name, "message" rows .last-msg-name.
            const nameEl = user.querySelector(".user-item-name, .last-msg-name");
            const NickName = nameEl?.innerHTML;
            if (!NickName) return;

            let userB = DEFAULT_USERS.filter((e) => { return e.name == NickName; });
            if (userB.length == 1) {
                chat = {
                    UserA: me.ID.Value,
                    UserB: userB[0].id.Value,
                    RoomID: roomIds.get(userB[0].id.Value) ?? null,
                    RequestApi: false,
                };
                navigate("/chat");
                return;
            }

            userB = DEFAULT_RECENT_MESSAGES.filter((e) => { return e.name == NickName; });
            if (userB.length == 1) {
                chat = {
                    UserA: me.ID.Value,
                    UserB: userB[0].id.Value,
                    RoomID: roomIds.get(userB[0].id.Value) ?? null,
                    RequestApi: true,
                };
                navigate("/chat");
                return;
            }

            console.warn("User not found in list:", NickName);
        });
    });
}
