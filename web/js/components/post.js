// NOTE: file name kept lowercase to match the existing project convention
// (navbar.js, modal.js, toast.js...). Import as `import PostCard from './post.js'`
// — file paths are case-sensitive on Linux.

/**
 * PostCard
 * Single post card in the feed.
 * Replaces the old global CreateCartPost with a clean ES module
 * that returns HTML (one component, no side effects).
 */
export default function PostCard({
    title = '',
    author = '',
    time = '',
    content = '',
    category = '',
    likes = 0,
    comments = 0,
} = {}) {
    return `
        <article class="post-card">
            <div class="post-header">
                <div>
                    <h3>${title}</h3>
                    <span>by <span class="user-info__name">${author}</span></span>
                </div>
                <span class="post-time">${time}</span>
            </div>
            <p>${content}</p>
            <span class="post-category">${category}</span>
            <div class="post-actions">
                <button class="like-btn">❤️ ${likes}</button>
                <button class="comment-btn">💬 ${comments}</button>
            </div>
        </article>
    `;
}
