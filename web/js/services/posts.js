
//import createPost from '../models/Post.js';
import {fetchPost} from "./../api/posts.js"


export let DEFAULT_POST = [];

export async function sendAllPost() {
    let response = await fetchPost();

    if (response.code == 200){
        if (!Array.isArray(response.body)) {
            DEFAULT_POST = ["not fond post"];
        }else{
            DEFAULT_POST = response.body
        }
    }else{
        DEFAULT_POST = ["error in fetch Categories"]
    }
}
// export default function CreatePostInDom(post){
//     let Posts = Array.from(document.querySelectorAll(".feed article"))

//     if(Posts.length == 1 && Posts[0].textContent.trim() === "not fond post"){    
//         Posts.shift()
//     }
//     Posts.unshift(post)
//     DEFAULT_POST = Posts
//     document.querySelector(".feed").innerHTML = DEFAULT_POST
// }