
export default function createUser({
    id = 0,
    name = '',
    handle = '',
    letter = '',
    avatarClass = 'avatar--mine',
    onlineStatus = '',
    lastMessage = '',
    unreadCount = 0,
    createdAt = '',
} = {}) {
    return {
        id,
        name,
        handle,
        letter: letter || (name ? name.charAt(0).toUpperCase() : ''),
        avatarClass,
        onlineStatus,
        lastMessage,
        unreadCount,
        createdAt,
    };
}
