import Navbar from './navbar.js';
import Sidebar from './Sidebar.js';
import Feed from './Feed.js';

/**
 * HomePage
 * Top-level composition: Navbar + Sidebar + Feed inside .home.
 * A single call renders the whole page as an HTML string.
 */
export default function HomePage({ nav = {}, sidebar = {}, feed = {} } = {}) {
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
