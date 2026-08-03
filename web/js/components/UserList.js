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
 */
export default function UserList(users = [], { variant = 'item' } = {}) {
    const listClass = LIST_CLASSES[variant] || LIST_CLASSES.item;
    const itemsHtml = users.map((user) => User(user, { variant })).join('');

    return `<div class="${listClass}">${itemsHtml}</div>`;
}
