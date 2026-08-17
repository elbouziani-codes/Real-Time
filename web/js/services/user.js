import createUser from "../models/User.js";
import { fetchAllUsers, resetUsersPaging } from "../api/users.js";

export const DEFAULT_USERS = [];
export const DEFAULT_RECENT_MESSAGES = [];

function normalizeUser(user = {}) {
    return createUser({
        id: user.ID ?? user.id ?? 0,
        name: user.NickName,
        handle: "@"+user.NickName,
        letter: user.Letter ?? user.letter ?? "",
        avatarClass: user.AvatarClass ?? user.avatarClass ?? "avatar--mine",
        onlineStatus: user.OnlineStatus ?? user.onlineStatus ?? "",
        lastMessage: user.LastMessage ?? user.lastMessage ?? "",
        unreadCount: user.UnreadCount ?? user.unreadCount ?? 0,
        createdAt: user.CreatedAt ?? user.createdAt ?? "",
    });
}

function formatLastMessageTime(lastMessageAt) {
    if (!lastMessageAt) return "";
    return String(lastMessageAt);
}

function sortUsersAlphabetically(users = []) {
    return [...users].sort((a, b) => String(a.name ?? "").localeCompare(String(b.name ?? ""), undefined, {
        sensitivity: "base",
        numeric: true,
    }));
}

function resortDefaultUsers() {
    DEFAULT_USERS.sort((a, b) => String(a.name ?? "").localeCompare(String(b.name ?? ""), undefined, {
        sensitivity: "base",
        numeric: true,
    }));
}

function sortRecentMessagesByLatestMessage(messages = []) {
    return [...messages].sort((a, b) => {
        const timeA = Number(a.createdAt ?? 0);
        const timeB = Number(b.createdAt ?? 0);
        return timeB - timeA;
    });
}

export async function seedAllUsers() {
    resetUsersPaging();
    DEFAULT_USERS.length = 0;
    DEFAULT_RECENT_MESSAGES.length = 0;

    const users = await fetchAllUsers();
    const normalizedUsers = users.map(normalizeUser);
    const { usersWithoutMessages, recentMessages } = filterUsers(normalizedUsers, users);
    DEFAULT_USERS.push(...sortUsersAlphabetically(usersWithoutMessages));
    resortDefaultUsers();
    DEFAULT_RECENT_MESSAGES.push(...sortRecentMessagesByLatestMessage(recentMessages));

    return {
        users: DEFAULT_USERS,
        recentMessages: DEFAULT_RECENT_MESSAGES,
    };
}

export async function loadMoreUsers() {
    const users = await fetchAllUsers();
    if (!users.length) return false;

    const normalizedUsers = users.map(normalizeUser);
    const { usersWithoutMessages, recentMessages } = filterUsers(normalizedUsers, users);

    DEFAULT_USERS.push(...sortUsersAlphabetically(usersWithoutMessages));
    resortDefaultUsers();
    DEFAULT_RECENT_MESSAGES.push(...sortRecentMessagesByLatestMessage(recentMessages));

    return true;
}

export function filterUsers(defaultUsers = [], rawUsers = defaultUsers) {
    const usersWithoutMessages = [];
    const recentMessages = [];

    for (let index = 0; index < defaultUsers.length; index += 1) {
        const user = defaultUsers[index];
        const rawUser = rawUsers[index] ?? user;
        const lastMessageAt = rawUser.LastMessageAt ?? rawUser.lastMessageAt ?? 0;
        const hasMessages = Number(lastMessageAt) > 0;

        if (hasMessages) {
            recentMessages.push({
                ...user,
                lastMessage: rawUser.lastMessage ?? rawUser.LastMessage ?? "",
                createdAt: formatLastMessageTime(lastMessageAt),
            });
            continue;
        }

        usersWithoutMessages.push(user);
    }

    return { usersWithoutMessages, recentMessages };
}
