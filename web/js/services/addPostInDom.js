




export default function CreatePostInDom(post){
    let Posts = Array.from(document.querySelectorAll(".feed article"))
    if(Posts.length == 1 && Posts[0].textContent.trim() === "not fond post"){    
        Posts.shift()
    }
    Posts.unshift(post)
    DEFAULT_POST = Posts
}