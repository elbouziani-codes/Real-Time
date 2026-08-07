/**
 * Post model — pure data structure (no rendering, no I/O, no seed data).
 * Seed/demo datasets live in js/services/seed.js.
 *
 * Fields:
 *  - id         unique identifier
 *  - title      post heading
 *  - content    body text
 *  - category   category name (e.g. "Programming")
 *  - author     author display name
 *  - likes      like count
 *  - comments   comment count
 *  - createdAt  timestamp / relative time label
 */
export default function createPostModel({
    id = 0,
    title = '',
    content = '',
    category = '',
    author = '',
    likes = 0,
    comments = 0,
    createdAt = '',
} = {}) {
    return {
        id,
        title,
        content,
        category,
        author,
        likes,
        comments,
        createdAt,
    };
}
