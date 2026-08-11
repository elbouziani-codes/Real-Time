

// Backend timestamps are unix seconds (posts.created_at, comments.created_at).
// Same formatting the feed card already uses, kept in one place so the post
// details page and the feed read alike.
export default function formatDateTime(seconds) {
    const time = Number(seconds);
    if (!time) return "";
    return new Date(time * 1000).toLocaleString();
}
