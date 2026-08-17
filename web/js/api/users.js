import { config } from "./../config/config.js";

let usersLoading = false;
let usersHasMore = true;

function userCursor(user) {
    const id = user?.ID ?? user?.id;
    if (typeof id === "string") return id;
    return id?.Value ?? id?.value ?? null;
}

export function resetUsersPaging() {
    usersLoading = false;
    usersHasMore = true;
    config.UsersCursor = null;
}

export async function fetchAllUsers() {
    if (usersLoading || !usersHasMore) return [];

    usersLoading = true;

    const params = new URLSearchParams({ limit: String(config.Userslimit) });
    if (config.UsersCursor) params.set("cursor", config.UsersCursor);

    try {
        const response = await fetch(`/api/users?${params.toString()}`, {
            method: "GET",
        });

        if (!response.ok) {
            // A 401 (e.g. the session was just invalidated by a logout racing
            // an in-flight navigation) is handled by the router's auth redirect;
            // report the page as empty instead of crashing the console.
            console.warn(`Failed to fetch users: ${response.status}`);
            return [];
        }

        const users = await response.json();

        if (!Array.isArray(users)) {
            throw new Error("Invalid users response");
        }

        if (users.length > 0) {
            const nextCursor = userCursor(users[users.length - 1]);
            if (!nextCursor || nextCursor === config.UsersCursor) {
                usersHasMore = false;
            } else {
                config.UsersCursor = nextCursor;
            }
        }

        if (users.length < config.Userslimit) {
            usersHasMore = false;
        }

        return users;
    } finally {
        usersLoading = false;
    }
}
