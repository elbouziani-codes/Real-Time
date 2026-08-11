import { Navbar } from "./../components/navbar.js";
import Sidebar from "../components/home/Sidebar.js";
import Feed from "./../components/home/Feed.js";
import createPost from "../components/home/creatPost.js";
import { navigate } from "./../router/router.js";
import { me } from "../services/me.js";

export default async function HomePage({ sidebar = {}, feed = {} } = {}) {
    return `
        <section class="home">
            ${Navbar(me)}
            <div class="home-layout">
                ${Sidebar(sidebar)}
                ${Feed()}
            </div>
            ${createPost()}
        </section>
    `;
}
