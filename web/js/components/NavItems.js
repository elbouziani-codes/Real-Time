
export default function NavItems() {
    const onChat = window.location.pathname == "/chat";

    return `
        <div class="nav-links">
            <button class="${onChat ? "" : "active "}homePage">🏠 Home</button>
            <button class="${onChat ? "active " : ""}chatPage">💬 Chat</button>
            ${onChat ? "" : `<button class="createPost">➕ Create Post</button>`}
        </div>
    `;
}
