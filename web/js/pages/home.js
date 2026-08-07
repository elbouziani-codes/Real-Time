import {Navbar} from './../components/navbar.js';
import Sidebar from './../components/Sidebar.js';
import Feed from './../components/Feed.js';
import {fetchPost, fetchCreatePost} from "./../api/posts.js"
import creatPost from "./../components/creatPost.js"
import {navigate} from "./../router/router.js"
import createPostModel from "./../models/Post.js"




async function sendAllPost() {
    let result = await fetchPost();

    if (!Array.isArray(result)) {
        return ["not fond post"];
    }

    return result;
}
export async function HomePage({ nav = {}, sidebar = {}, feed = {} } = {}) {
    let Posts = await sendAllPost()
    DEFAULT_POST = [...DEFAULT_POST, ...Posts]
    return `
        <section class="home">
            ${Navbar(nav)}
            <div class="home-layout">
                ${Sidebar(sidebar)}
                ${Feed(Posts)}
            </div>
            ${creatPost()}
        </section>
    `;
}


export function HomeListener() {
    const createPost = document.querySelector(".createPost");
    const cancelCreatePost = document.getElementById("cancelCreatePost");
    const modalOverlay = document.querySelector(".modal-overlay");

    createPost?.addEventListener("click", () => {
        modalOverlay.className = "modal-overlay visible"
    })
    cancelCreatePost?.addEventListener("click", () => {
        modalOverlay.className = "modal-overlay hidden"
    })
}



async function submitPost(e){
    e.preventDefault()

    const data = {
        title:document.querySelector(".modal-input").value,
        category_id:Array.from(document.querySelectorAll(".modal-overlay .category-item input")).filter((e) => {return e.checked}).map((e) => {return e.id}),
        content:document.querySelector("modal-textarea").value,
    }

    // fetch 
    let Response =  await fetchCreatePost(data)
    if(Response.code == 200){
        PostCard(createPostModel(Response.body))
    }
}