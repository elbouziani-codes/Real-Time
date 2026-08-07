import Navbar from './../components/navbar.js';
import Sidebar from './../components/Sidebar.js';
import Feed from './../components/Feed.js';
import fetchPost from "./../api/posts.js"
/**
 * HomePage
 * Top-level composition: Navbar + Sidebar + Feed inside .home.
 * A single call renders the whole page as an HTML string.
 *
 * Model-driven props (each component falls back to its model seed):
 *  - nav      { navItems, user: User model }
 *  - sidebar  { users, categories, filters, recentMessages }
 *  - feed     { title, posts: Post model[] }
 */

async function sendAllPost() {
    let result = await fetchPost();

    if (!Array.isArray(result)) {
        return [];
    }

    return result;
}
export default async function HomePage({ nav = {}, sidebar = {}, feed = {} } = {}) {
    let Posts = await sendAllPost()
    return `
        <section class="home">
            ${Navbar(nav)}
            <div class="home-layout">
                ${Sidebar(sidebar)}
                ${FeefetchPost()d(Posts)}
            </div>
        </section>
    `;
}