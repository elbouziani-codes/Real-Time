import PostCard from './post.js';

/**
 * Feed
 * Main posts feed: heading + a list of PostCard components.
 * All data comes from Post models (defaults to the model seed).
 */
export default function Feed(posts) {
    const postsHtml = posts.map((e) => {return PostCard(e) }).join('');

    return `
        <main class="feed">
            <h2>📰 Latest Posts</h2>
            ${postsHtml}
        </main>
    `;
}
