
import UserAvatar from "./../UserAvatar.js";
import { authorModel } from "./PostDetails.js";
import formatDateTime from "./../../utils/time.js";
import { escapeHTML } from "./../../utils/helpers.js";

/**
 * Comments
 * Comment block of the post details page: the form the backend already accepts
 * (POST /api/comments) plus the comments of this post (GET /api/comments/{id}).
 *
 * A CommentInfo carries ID, Author, Content and CreatedAt. Its like counters are
 * not filled by the backend query, so no reaction control is rendered here.
 */

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

// CommentsList is the inner html of the list only, so the listener can refresh
// the comments after a new one is posted without rebuilding the form.
export function CommentsList(comments = []) {
    if (!comments.length) {
        return `<p class="empty-state">No comment yet. Be the first one.</p>`;
    }
    return comments.map((e) => { return Comment(e); }).join("");
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
