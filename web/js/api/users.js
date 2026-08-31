import { config } from "./../config/config.js";

let usersLoading = false;
let usersHasMore = true;

export function resetUsersPaging() {
    usersLoading = false;
    usersHasMore = true;
    config.UsersCursor = null;
}

// One page of contacts, ordered by last message then nickname (server side).
// Advances the cursor; returns [] when exhausted or on failure.
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
            return [];
        }

        const users = await response.json();
        if (!Array.isArray(users)) {
            throw new Error("Invalid users response");
        }
        const nextCursor = users.at(-1)?.ID;
        if (!nextCursor || nextCursor === config.UsersCursor) usersHasMore = false;
        else config.UsersCursor = nextCursor;

        if (users.length < config.Userslimit) usersHasMore = false;

        return users;
    } finally {
        usersLoading = false;
    }
}

// One profile by id, used when a user nobody paged in yet comes online.
export async function fetchUser(id) {
    const response = await fetch(`/api/users/${id}`);
    if (!response.ok) return null;
    return response.json();
}
