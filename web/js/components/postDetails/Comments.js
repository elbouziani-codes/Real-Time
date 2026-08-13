import UserAvatar from "./../UserAvatar.js";
import { authorModel } from "./PostDetails.js";
import formatDateTime from "./../../utils/time.js";
import { escapeHTML } from "./../../utils/helpers.js";


function Comment(comment = {}) {
    const { Author = {}, Content = '', CreatedAt = 0 } = comment;
    const author = authorModel(Author);

    return `
        <div class="post-details__comment">
            ${UserAvatar(author, { sizeClass: 'user-avatar--md', className: 'post-details__comment-avatar' })}
            <div class="post-details__comment-body">
                <div class="post-details__comment-header">
                    <span class="post-details__comment-author">${escapeHTML(author.name)}</span>
                    <span class="post-details__comment-time">${formatDateTime(CreatedAt)}</span>
                </div>
                <p class="post-details__comment-text">${escapeHTML(Content)}</p>
            </div>
        </div>
    `;
}


export function CommentsList(comments = []) {
    if (!comments.length) {
        return `
            <div class="empty-state">
                <span class="empty-state__icon">💬</span>
                <p class="empty-state__title">No comments yet</p>
                <p class="empty-state__text">Be the first to start the conversation.</p>
            </div>
        `;
    }
    return comments.map((comment) => Comment(comment)).join("");
}

export default function Comments(comments = []) {
    return `
        <section class="post-details__comments">
            <h2 class="post-details__comments-title">💬 Comments <span class="post-details__comments-count">(${comments.length})</span></h2>

            <form class="post-details__comment-form">
                <textarea class="post-details__comment-input" placeholder="Write a comment..." rows="3" maxlength="4096"></textarea>
                <p class="post-details__comment-status"></p>
                <button type="submit" class="post-details__comment-submit">Post Comment</button>
            </form>

            <div class="post-details__comments-list">
                ${CommentsList(comments)}
            </div>
        </section>
    `;
}
