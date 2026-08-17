
export default function createMessage({
    id = 0,
    roomId = 0,
    senderId = 0,
    content = '',
    createdAt = 0,
} = {}) {
    return {
        id,
        roomId,
        senderId,
        content,
        createdAt,
    };
}
