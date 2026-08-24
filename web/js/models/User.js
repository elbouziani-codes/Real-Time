export default function createUser({
    id = 0,
    name = '',
    handle = '',
    letter = '',
    avatarClass = 'avatar--mine',
    online = false,
    lastMessage = '',
    lastMessageAt = 0,
} = {}) {
    return {
        id,
        name,
        handle,
        letter: letter || (name ? name.charAt(0).toUpperCase() : ''),
        avatarClass,
        online,
        lastMessage,
        lastMessageAt: Number(lastMessageAt),
    };
}
