import PostCard from './post.js';

/**
 * Feed
 * Main posts feed: heading + a list of PostCard components.
 */
export default function Feed({ title = '📰 Latest Posts', posts = [] } = {}) {
    const postsHtml = posts.map((post) => PostCard(post)).join('');

    return `
        <main class="feed">
            <h2>${title}</h2>
            ${postsHtml}
        </main>
    `;
}
