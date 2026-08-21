import { auth } from "./../pages/auth.js";
import authListener from "./../listeners/auth.js";
import {CreateMe, me} from "./../services/me.js";
import {homeListenerUser} from "./../listeners/users.js"
import chatListener from "./../listeners/chat.js"
import {connectSocket} from "./../websocket/socket.js"
import {initChatSocket} from "./../websocket/chat.js"
import chat from "../pages/chat.js";
import endPage from "../pages/end.js";
import endListener from "../listeners/end.js";

import { navBarListener } from "./../components/navbar.js";
import HomePage from "../pages/home.js";
import {seedAllUsers} from "./../services/user.js"
import {HomeListener, HomeScrollListener} from "./../listeners/home.js";

import pagePost from "../pages/postDetails.js";
import PostDetailsListener from "./../listeners/postDetails.js";
import { seedCategories } from "./../services/categories.js";
import { sendAllPost } from "./../services/posts.js";

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
    app.innerHTML = `
            <h1>
                404
            </h1>
        `;
  }
}

function navigate(path) {
  if (window.location.pathname === path) return;
  history.pushState({}, "", path);
  router();
}

async function Fetching(path) {
  await CreateMe(path);
  // CreateMe redirects unauthenticated visitors to /login; never seed data (or
  // open the socket) for a session that is about to be redirected.
  if (path == "/login" || path == "/register") return;
  if (!me?.ID?.Value) return;

  switch (path) {
    case "/":
      await seedAllUsers();
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
  initChatSocket();
}

async function Listening(path) {
  switch (path) {
    case "/login":
    case "/register":
      authListener();
      break;
    case "/":
      navBarListener();
      homeListenerUser()
      HomeListener();
      await HomeScrollListener()
      break;
    case "/chat":
      navBarListener();
      chatListener();
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
