// NOTE: file name kept lowercase to match the existing project convention
// (navbar.js, modal.js, toast.js...). Import as `import PostCard from './post.js'`
// — file paths are case-sensitive on Linux.

/**
 * PostCard
 * Single post card in the feed. All data comes from a Post model:
 * { title, author, content, category, likes, comments, createdAt }.
 */
export default function PostCard(post = {}) {
    const {
        title = '',
        author = '',
        content = '',
        category = '',
        likes = 0,
        comments = 0,
        createdAt = '',
    } = post;

    return `
        <article class="post-card">
            <div class="post-header">
                <div>
                    <h3>${title}</h3>
                    <span>by <span class="user-info__name">${author}</span></span>
                </div>
                <span class="post-time">${createdAt}</span>
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
