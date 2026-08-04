import auth from "../pages/auth";
import chat from "../pages/chat";
import HomePage from "../pages/home";
import pagePost from "../pages/postDaitaile";

let Routes = {
    "/":HomePage,
    "/post":pagePost,
    "/chat":chat,
    "/auth":auth
}
function router(){
    const path = window.location.pathname;
    const page = routes[path];
    const app = document.getElementById("app");
    if (page){
        app.innerHTML = page
    }else{
        app.innerHTML = "<h1>404</h1>"
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