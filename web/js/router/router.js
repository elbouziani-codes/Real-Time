import { auth } from "./../pages/auth.js";
import chat from "../pages/chat.js";
import endPage from "../pages/end.js";
import ErrorPage from "../pages/notFound.js";
import pagePost from "../pages/postDetails.js";
import HomePage from "../pages/home.js";

import endListener from "../listeners/end.js";
import {chatListener ,onSideBareScroll}  from "./../listeners/chat.js";
import authListener from "./../listeners/auth.js";
import PostDetailsListener from "./../listeners/postDetails.js";
import { HomeListener, HomeScrollListener } from "./../listeners/home.js";

import { navBarListener } from "./../components/navbar.js";

import { CreateMe, me } from "./../services/me.js";
import { seedAllUsers, resetUsers } from "./../services/user.js";
import { seedCategories } from "./../services/categories.js";
import { sendAllPost} from "./../services/posts.js";

import { connectSocket } from "./../websocket/socket.js";


const routes = {
    "/": HomePage,
    "/postDetails": pagePost,
    "/chat": chat,
    "/end": endPage,
    "/login": auth,
    "/register": auth,
};

async function router() {
    const path = window.location.pathname;
    const page = routes[path];
    const app = document.getElementById("app");
    if (page) {
        await Fetching(path);
        if (window.location.pathname !== path) return;
        app.innerHTML = await page();
        await Listening(path);
    } else {
        app.innerHTML = ErrorPage({});
    }
}

function navigate(path) {
    if (window.location.pathname === path) return;
    history.pushState({}, "", path);
    router();
}

async function Fetching(path) {
    await CreateMe(path);

    if (path == "/login" || path == "/register") return;
    if (!me?.ID) return;

    switch (path) {
        case "/":
            await sendAllPost();
            await seedCategories();
            break;
        case "/chat":
            await seedAllUsers();
            break;
        default:
            return null;
    }
    connectSocket();
}

async function Listening(path) {
    switch (path) {
        case "/login":
        case "/register":
            authListener();
            break;
        case "/":
            navBarListener();
            HomeListener();
            await HomeScrollListener();
            break;
        case "/chat":
            navBarListener();
            chatListener();
            await onSideBareScroll();
            break;
        case "/end":
            endListener();
            break;
        case "/postDetails":
            navBarListener();
            await PostDetailsListener();
            break;
    }
}

window.addEventListener("popstate", router);

export { router, navigate };
