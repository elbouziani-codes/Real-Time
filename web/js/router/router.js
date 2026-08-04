import auth from "../pages/auth.js";
import chat from "../pages/chat.js";
import HomePage from "../pages/home.js";
import pagePost from "../pages/postDaitaile.js";


const routes = {

    "/": HomePage,

    "/post": pagePost,

    "/chat": chat,

    "/auth": auth

};



function router(){

    const path = window.location.pathname;

    const page = routes[path];

    const app = document.getElementById("app");


    if(page){

        app.innerHTML = page();

    }else{

        app.innerHTML = `
            <h1>
                404
            </h1>
        `;

    }

}



function navigate(path){

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