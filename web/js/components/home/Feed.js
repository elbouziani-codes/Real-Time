import PostCard from './post.js';
//import { DEFAULT_POST } from "../../services/posts.js";


export default function Feed(posts) {
    return `
        <main id="feed">
            <h2>📰 Latest Posts</h2>
            ${posts.map((post) => PostCard(post)).join('')}
        </main>
    `;
}
