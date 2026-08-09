import { authListener, auth } from "./../pages/auth.js";
import chat from "../pages/chat.js";
import { navBarListener } from "./../components/navbar.js"
import {HomePage, HomeListener} from "../pages/home.js";
import pagePost from "../pages/postDaitaile.js";
import {getProfile} from "./../api/auth.js"
import {seedCategories} from "./../services/categories.js"
import {sendAllPost} from "./../services/posts.js"

const routes = {

    "/": HomePage,

    "/post": pagePost,

    "/chat": chat,

    "/login": auth,

    "/register": auth

};



async function router(){
    const path = window.location.pathname;
    const page = routes[path];
    const app = document.getElementById("app");
    await getProfile(path)  

    if(page){
        await Fetching(path)
        app.innerHTML = await page();
        await Listening(path)

    }else{

        app.innerHTML = `
            <h1>
                404
            </h1>
        `;

    }
    

}



function navigate(path){

    if(window.location.pathname === path)
        return;

    history.pushState(
        {},
        "",
        path
    );
    router();
}


async function Fetching(path) {
    switch (path) {
        case "/":
            await sendAllPost();
            await seedCategories();
        default:
            console.log("fetch")
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
            console.log("deee")
            navBarListener()
            HomeListener()
        break;
        case "/chat":
            navBarListener()
        break;
    }
}
window.addEventListener(
    "popstate",
    router
);


export {
    router,
    navigate
};