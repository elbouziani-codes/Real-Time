
export default function formatDateTime(seconds) {
    const time = Number(seconds);
    if (!time) return "";
    return new Date(time * 1000).toLocaleString();
}
