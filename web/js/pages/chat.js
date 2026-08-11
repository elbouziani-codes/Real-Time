import { Navbar } from "./../components/navbar.js";
import Sidebar from "../components/home/Sidebar.js";
import Feed from "../components/home/Feed.js";

export default function chat({ nav = {} } = {}) {
    return `
        <section class="home">
            ${Navbar(nav)}
            
        </section>
    `;
}
