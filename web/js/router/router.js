import { auth } from "./../pages/auth.js";
import authListener from "./../listeners/auth.js";
import CreateMe from "./../services/me.js";

import chat from "../pages/chat.js";

import { navBarListener } from "./../components/navbar.js";
import HomePage from "../pages/home.js";
import {HomeListener, HomeScrollListener} from "./../listeners/home.js";

import pagePost from "../pages/postDetails.js";
import PostDetailsListener from "./../listeners/postDetails.js";
import { seedCategories } from "./../services/categories.js";
import { sendAllPost } from "./../services/posts.js";

const routes = {
  "/": HomePage,

  "/postDetails": pagePost,

  "/chat": chat,

  "/login": auth,

  "/register": auth,
};

async function router() {
  const path = window.location.pathname;
  const page = routes[path];
  const app = document.getElementById("app");

  if (page) {
    await Fetching(path);
    // Fetching may redirect (an unauthenticated user is sent to /login), and
    // that redirect already rendered its own page.
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
  switch (path) {
    case "/":
      await sendAllPost();
      await seedCategories();
    default:
      return null;
  }
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
      await HomeScrollListener()
      break;
    case "/chat":
      navBarListener();
      break;
    case "/postDetails":
      navBarListener();
      await PostDetailsListener();
      break;
  }
}
window.addEventListener("popstate", router);

export { router, navigate };
