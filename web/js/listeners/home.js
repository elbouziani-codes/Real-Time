import { fetchCreatePost } from "./../api/posts.js";
import { reactionState } from "./../components/home/post.js";
import throttle from "./../utils/helpers.js";
import { config } from "./../config/config.js";
import { navigate } from "./../router/router.js";
import {state } from "/main.js"
import { refreshFeed } from "/handlers/feed.js"
import { reactToPost } from "/fetcher.js"



let windowScrollHandler = null;

export function initFeedListeners() {
		HomeListener();
};

export function HomeListener() {
    const createPost = document.querySelector(".createPost");
    const cancelCreatePost = document.getElementById("cancelCreatePost");
    const modalOverlay = document.querySelector(".modal-overlay");
    const modalForm = modalOverlay?.querySelector("form");

    createPost?.addEventListener("click", () => {
        modalOverlay.className = "modal-overlay visible";
    });
    cancelCreatePost?.addEventListener("click", () => {
        modalOverlay.className = "modal-overlay hidden";
    });
    modalForm?.addEventListener("submit", submitPost);
    document.getElementById("feed")?.addEventListener("click", reactListener);
    document.querySelector(".sidebar")?.addEventListener("change", filterListener);
}

async function filterListener(e) {
    if (!e.target.matches(".category-item input[type=checkbox]")) return;

    const boxes = Array.from(document.querySelectorAll(".sidebar .category-item input[type=checkbox]"));
    const categories = boxes
        .filter((input) => input.checked && input.id != "likedFilter")
        .map((input) => input.id);
    const liked = document.getElementById("likedFilter")?.checked || false;

    state.postsCollections.updateFilters(categories, liked);
		
    await state.postsCollections.morePosts(); 
	refreshFeed();
};

async function reactListener(e) {
    const card = e.target.closest(".post-card");
    const postID = card?.dataset.postId;
    if (!postID) return;

    const button = e.target.closest(".like-btn, .comment-btn");
    if (!button) {
		// to do later
        //navigate("/postDetails?id=" + postID);
        return;
    }

    const isLike = button.classList.contains("like-btn");

    const post = await reactToPost(postID, isLike);
	console.log(post)

    if (!post) return;

    const likeBtn = card.querySelector(".like-btn");
    const disLikeBtn = card.querySelector(".comment-btn");
    const { liked, disliked } = reactionState(post);
    likeBtn.textContent = `👍 ${post.Likes}`;
    disLikeBtn.textContent = `👎 ${post.DisLikes}`;
    likeBtn.classList.toggle("active", liked);
    disLikeBtn.classList.toggle("active", disliked);
}





function currentStateKey() {
    return config.postsCursor + "|" + JSON.stringify(config.postsFilters);
}

export async function HomeScrollListener() {
    let lastScrollKey = null;

    const throttledLoadPosts = throttle(() => {
        lastScrollKey = currentStateKey();
        sendAllPost();
    }, 1500);

    // The window listener persists across SPA mounts: detach the previous one
    // so revisiting the home page never stacks duplicate scroll handlers.
    if (windowScrollHandler) {
        window.removeEventListener("scroll", windowScrollHandler);
    }
    windowScrollHandler = () => {
        const scrollTop = window.scrollY;
        const windowHeight = window.innerHeight;
        const documentHeight = document.documentElement.scrollHeight;

        if (scrollTop + windowHeight >= documentHeight - 100 && lastScrollKey != currentStateKey()) {
            throttledLoadPosts();
        }
    };
    window.addEventListener("scroll", windowScrollHandler);
}

async function submitPost(e) {
    e.preventDefault();

    const modal = document.querySelector(".modal-overlay");
    const status = modal?.querySelector(".modal-status");
    const submit = modal?.querySelector(".modal-submit-btn");
    if (status) status.textContent = "";

    const data = {
        title: document.querySelector(".modal-input").value.trim(),
        category_ids: Array.from(document.querySelectorAll(".modal-overlay .category-item input"))
            .filter((box) => box.checked)
            .map((box) => box.id),
        content: document.querySelector(".modal-textarea").value.trim(),
    };

    if (!data.title || !data.content) {
        if (status) status.textContent = "Title and content are required.";
        return;
    }
    if (data.category_ids.length === 0) {
        if (status) status.textContent = "Pick at least one category.";
        return;
    }

    if (submit) submit.disabled = true;
    try {
        const response = await fetchCreatePost(data);
        if (response.code == 401) {
            navigate("/login");
            return;
        }
        if (response.code != 200) {
            const message = typeof response.body == "string" && response.body.trim()
                ? response.body.trim()
                : "Could not create the post. Please try again.";
            if (status) status.textContent = message;
            return;
        }

        if (modal) modal.className = "modal-overlay hidden";
        document.querySelector(".modal-input").value = "";
        document.querySelector(".modal-textarea").value = "";
        document.querySelectorAll(".modal-overlay .category-item input").forEach((box) => { box.checked = false; });

        // The new post is at the top of the next page: reload from the start.
        applyPostFilters(config.postsFilters);
        await sendAllPost();
    } finally {
        if (submit) submit.disabled = false;
    }
}
