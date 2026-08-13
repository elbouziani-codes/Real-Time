import PostCard from './post.js';
import { DEFAULT_POST } from "../../services/posts.js";


export default function Feed() {
    return `
        <main class="feed">
            <h2>📰 Latest Posts</h2>
            ${DEFAULT_POST.map((post) => PostCard(post)).join('')}
        </main>
    `;
}
