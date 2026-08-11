
import createPostModel from "./../models/Post.js"
import {fetchPost, fetchCreatePost} from "./../api/posts.js"
import {sendAllPost, reactToPost} from "./../services/posts.js"
import {reactionState} from "./../components/home/post.js"
import throttle from "./../utils/helpers.js"
import {config} from "./../config/config.js"
import {navigate} from "./../router/router.js"

export function HomeListener() {
    const createPost = document.querySelector(".createPost");
    const cancelCreatePost = document.getElementById("cancelCreatePost");
    const modalOverlay = document.querySelector(".modal-overlay");
    const submitCreatePost = document.querySelector(".modal-overlay .modal-submit-btn");

    createPost?.addEventListener("click", () => {
        modalOverlay.className = "modal-overlay visible"
    })
    cancelCreatePost?.addEventListener("click", () => {
        modalOverlay.className = "modal-overlay hidden"
    })
    submitCreatePost?.addEventListener("click",  submitPost)

    // Delegated on the feed so posts appended by the infinite scroll react too.
    document.querySelector(".feed")?.addEventListener("click", reactListener)
}

async function reactListener(e) {
    const card = e.target.closest(".post-card");
    const postID = card?.dataset.postId;
    if (!postID) return;

    const button = e.target.closest(".like-btn, .comment-btn");
    // Anywhere else on the card opens the post details page.
    if (!button) {
        navigate("/postDetails?id=" + postID);
        return;
    }

    const isLike = button.classList.contains("like-btn");
    const post = await reactToPost(postID, isLike);
    if (!post) return;

    const likeBtn = card.querySelector(".like-btn");
    const disLikeBtn = card.querySelector(".comment-btn");
    const {liked, disliked} = reactionState(post);
    likeBtn.textContent = `/\\ ${post.Likes}`;
    disLikeBtn.textContent = `\\/ ${post.DisLikes}`;
    likeBtn.classList.toggle("active", liked);
    disLikeBtn.classList.toggle("active", disliked);
}

export async function HomeScrollListener() {
    const throttledLoadPosts = throttle(() => {
        sendAllPost(true);
    }, 1500);
    let lastOffset = 0 
    window.addEventListener("scroll", async () => {
        const scrollTop = window.scrollY;
        const windowHeight = window.innerHeight;
        const documentHeight = document.documentElement.scrollHeight;

        if (scrollTop + windowHeight >= documentHeight - 100 && lastOffset != config.offsetPost) {
            throttledLoadPosts()
            lastOffset = config.offsetPost
        }
    });
}

async function submitPost(e){
    e.preventDefault()

    const data = {
        title:document.querySelector(".modal-input").value,
        category_ids:Array.from(document.querySelectorAll(".modal-overlay .category-item input")).filter((e) => {return e.checked}).map((e) => {return e.id}),
        content:document.querySelector(".modal-textarea").value,
    }
    // fetch 
    let Response =  await fetchCreatePost(data)
}