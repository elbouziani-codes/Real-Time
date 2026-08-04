import User from './User.js';

const LIST_CLASSES = {
    item: 'users-list',
    message: 'last-messages-list',
    conversation: 'conversations-list',
};

/**
 * UserList
 * Renders a list of User models.
 * Only the surrounding container changes between variants:
 *  - 'item'         -> .users-list
 *  - 'message'      -> .last-messages-list
 *  - 'conversation' -> .conversations-list
 *
 * `listClass` overrides the container class when the markup pairs a row
 * variant with a different wrapper (e.g. message rows inside .users-list).
 */
export default function UserList(users = [], { variant = 'item', listClass = '' } = {}) {
    const containerClass = listClass || LIST_CLASSES[variant] || LIST_CLASSES.item;
    const itemsHtml = users.map((user) => User(user, { variant })).join('');

    return `<div class="${containerClass}">${itemsHtml}</div>`;
}
