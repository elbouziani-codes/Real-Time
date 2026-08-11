// NOTE: file name kept lowercase to match the existing project convention
// (navbar.js, modal.js, toast.js...). Import as `import PostCard from './post.js'`
// — file paths are case-sensitive on Linux.

import { ZERO_UUID } from "./../../config/config.js";



// reactionState reports which of the two buttons the signed in user currently
// holds. LikeInfo carries that reaction, and a nil id means "no reaction yet".
export function reactionState(post = {}) {
    const {LikeInfo = {}} = post;
    const reacted = Boolean(LikeInfo.ID?.Value) && LikeInfo.ID.Value != ZERO_UUID;
    return {liked: reacted && LikeInfo.IsLike == true, disliked: reacted && LikeInfo.IsLike == false};
}

export default function PostCard(post = {}) {
    if (typeof(post) == "string"){
        return `
        <article class="post-card">
            <h3>not fond post</h3>
        </article>
        `
    }
    const {
        ID = {},
        Title = '',
        Author = {},
        Content = '',
        Categories = [],
        Likes = 0,
        DisLikes = 0,
        CreatedAt = '',
    } = post;
    const createdAt = new Date(CreatedAt * 1000).toLocaleString();
    const {liked, disliked} = reactionState(post);
    const likeClass = liked ? " active" : "";
    const disLikeClass = disliked ? " active" : "";

    return `
        <article class="post-card" data-post-id="${ID.Value ?? ''}">
            <div class="post-header">
                <div>
                    <h3>${Title}</h3>
                    <span>by <span class="user-info__name">${Author.NickName}</span></span>
                </div>
                <span class="post-time">${createdAt}</span>
            </div>
            <p>${Content}</p>
            <div class="post-categories"> ${Categories.map((e) => { return `<span class="post-category">${e.title}</span>`; }).join("")} </div>
            <div class="post-actions">
                <button class="like-btn${likeClass}">/\\ ${Likes}</button>
                <button class="comment-btn${disLikeClass}">\\/ ${DisLikes}</button>
            </div>
        </article>
    `;
}
