import { Navbar } from "./../components/navbar.js";
import ChatPage from "../components/chat/ChatPage.js";
import { me } from "../services/me.js";

export default function chat() {
    return `
        ${Navbar(me)}
        ${ChatPage()}
    `;
}
