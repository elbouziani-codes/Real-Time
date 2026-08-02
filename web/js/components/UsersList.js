import User from './User.js';

const LIST_CLASSES = {
    item: 'users-list',
    message: 'last-messages-list',
    conversation: 'conversations-list',
};

/**
 * UsersList
 * Renders a list of User components.
 * Only the surrounding container changes between variants:
 *  - 'item'         -> .users-list
 *  - 'message'      -> .last-messages-list
 *  - 'conversation' -> .conversations-list
 */
export default function UsersList({ variant = 'item', users = [] } = {}) {
    const listClass = LIST_CLASSES[variant] || LIST_CLASSES.item;
    const itemsHtml = users.map((user) => User({ variant, ...user })).join('');

    return `<div class="${listClass}">${itemsHtml}</div>`;
}
