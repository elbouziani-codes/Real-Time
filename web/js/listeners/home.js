
import createPostModel from "./../models/Post.js"
import {fetchPost, fetchCreatePost} from "./../api/posts.js"
import {sendAllPost, reactToPost, applyPostFilters} from "./../services/posts.js"
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
    document.querySelector(".feed")?.addEventListener("click", reactListener)
    document.querySelector(".sidebar")?.addEventListener("change", filterListener)
}

async function filterListener(e) {
    if (!e.target.matches(".category-item input[type=checkbox]")) return;

    const boxes = Array.from(document.querySelectorAll(".sidebar .category-item input[type=checkbox]"));
    const categories = boxes
        .filter((input) => input.checked && input.id != "likedFilter")
        .map((input) => input.id);
    const liked = document.getElementById("likedFilter")?.checked || false;

    applyPostFilters({categories, liked});

    const dispatchedKey = currentStateKey();

    const feed = document.querySelector(".feed");
    if (feed) feed.innerHTML = `<h2>📰 Latest Posts</h2>`;
    await sendAllPost(true);

    lastScrollKey = dispatchedKey;
}

async function reactListener(e) {
    const card = e.target.closest(".post-card");
    const postID = card?.dataset.postId;
    if (!postID) return;

    const button = e.target.closest(".like-btn, .comment-btn");
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
    likeBtn.textContent = `👍 ${post.Likes}`;
    disLikeBtn.textContent = `👎 ${post.DisLikes}`;
    likeBtn.classList.toggle("active", liked);
    disLikeBtn.classList.toggle("active", disliked);
}


let lastScrollKey = null;

function currentStateKey() {
    return config.postsCursor + "|" + JSON.stringify(config.postsFilters);
}

export async function HomeScrollListener() {
    const throttledLoadPosts = throttle(() => {
        lastScrollKey = currentStateKey();
        sendAllPost(true);
    }, 1500);

    window.addEventListener("scroll", async () => {
        const scrollTop = window.scrollY;
        const windowHeight = window.innerHeight;
        const documentHeight = document.documentElement.scrollHeight;

        if (scrollTop + windowHeight >= documentHeight - 100 && lastScrollKey != currentStateKey()) {
            throttledLoadPosts()
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