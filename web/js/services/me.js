import { fetchMe, fetchLogout } from "./../api/auth.js";
import { navigate } from "../router/router.js";
import { resetPosts } from "./posts.js";
import { disconnectSocket } from "../websocket/socket.js";

export let me = {};

export async function CreateMe(path) {
    let response = await fetchMe();
    if (!response) response = { code: 0, body: "Network error" };

    if (response.code == 200) {
        me = { ...response.body };
        if (path == "/login" || path == "/register") {
            navigate("/");
        }
    } else {
        if (path != "/login" && path != "/register") {
            navigate("/login");
        }
    }
}

export async function logoutMe() {
    const response = await fetchLogout();
    if (response.code != 200 && response.code != 401) {
        ("logout failed:", response.body);
        return false;
    }
    me = {};
    resetPosts();
    disconnectSocket();
    return true;
}
