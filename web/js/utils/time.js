
function normalizeTimestamp(timestamp) {
    const value = Number(timestamp);
    if (!Number.isFinite(value) || value <= 0) {
        return 0;
    }

    // Backend timestamps in this app are usually Unix seconds, but some
    // client-side models already carry millisecond or nanosecond values.
    // Normalize all of them to milliseconds.
    if (value > 1e15) {
        return Math.floor(value / 1e6);
    }
    return value < 1e12 ? value * 1000 : value;
}

export default function formatDateTime(timestamp) {
    const now = Date.now();
    const createdAt = normalizeTimestamp(timestamp);
    const diff = now - createdAt;

    const seconds = Math.floor(diff / 1000);

    if (seconds < 60) {
        return `${seconds}s ago`;
    }

    const minutes = Math.floor(seconds / 60);

    if (minutes < 60) {
        return `${minutes}m ago`;
    }

    const hours = Math.floor(minutes / 60);

    if (hours < 24) {
        return `${hours}h ago`;
    }

    const days = Math.floor(hours / 24);

    if (days < 30) {
        return `${days}d ago`;
    }

    const months = Math.floor(days / 30);

    if (months < 12) {
        return `${months}mo ago`;
    }

    const years = Math.floor(days / 365);

    return `${years}y ago`;
}
