import PostCard from './post.js';
import { DEFAULT_POSTS } from '../services/seed.js';

/**
 * Feed
 * Main posts feed: heading + a list of PostCard components.
 * All data comes from Post models (defaults to the model seed).
 */
export default function Feed({ title = '📰 Latest Posts', posts = DEFAULT_POSTS } = {}) {
    const postsHtml = posts.map((post) => PostCard(post)).join('');

    return `
        <main class="feed">
            <h2>${title}</h2>
            ${postsHtml}
        </main>
    `;
}
