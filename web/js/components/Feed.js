import PostCard from './post.js';
import {DEFAULT_POST} from "./../services/posts.js"
/**
 * Feed
 * Main posts feed: heading + a list of PostCard components.
 * All data comes from Post models (defaults to the model seed).
 */
export default function Feed() {
    const postsHtml = DEFAULT_POST.map((e) => {return PostCard(e) }).join('');

    return `
        <main class="feed">
            <h2>📰 Latest Posts</h2>
            ${postsHtml}
        </main>
    `;
}
