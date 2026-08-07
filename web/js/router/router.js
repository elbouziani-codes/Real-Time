import { authListener, auth } from "./../pages/auth.js";
import chat from "../pages/chat.js";
import { navBarListener } from "./../components/navbar.js"
import {HomePage, HomeListener} from "../pages/home.js";
import pagePost from "../pages/postDaitaile.js";
import {getProfile} from "./../api/auth.js"


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
        app.innerHTML = await page();
        if(path === "/login" || path === "/register"){
            authListener();
        }else if(path === "/"){
            navBarListener()

            HomeListener()
        }else if(path === "/chat"){
            navBarListener()
        }

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



window.addEventListener(
    "popstate",
    router
);


export {
    router,
    navigate
};