import Navbar from './../components/navbar.js';
import Sidebar from './../components/Sidebar.js';
import Feed from './../components/Feed.js';

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
export default function chat({ nav = {}, sidebar = {}, feed = {} } = {}) {
    return `
        <section class="home">
            ${Navbar(nav)}
            <div class="home-layout">
                ${Sidebar(sidebar)}
                ${Feed(feed)}
            </div>
        </section>
    `;
}