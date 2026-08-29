import { Navbar } from "./../components/navbar.js";
import Sidebar from "../components/home/Sidebar.js";
import Feed from "./../components/home/Feed.js";
import createPost from "../components/home/creatPost.js";
//import { me } from "../services/me.js";

export  default async function homePage(state) {
    return `
        <section class="home">
            ${Navbar(state.user)}
            <div class="home-layout">
               ${Sidebar(state)}
                ${Feed(state.postsCollections.get())}
            </div>
            ${createPost()}
        </section>
    `;
}
