import User from './User.js';

const LIST_CLASSES = {
    item: 'users-list',
    message: 'last-messages-list',
};


export default function UserList(users = [], { variant = 'item', listClass = '', listId = '' } = {}) {
    const containerClass = listClass || LIST_CLASSES[variant] || LIST_CLASSES.item;
    const idAttr = listId ? ` id="${listId}"` : '';
    const itemsHtml = users.map((user) => User(user, { variant })).join('');

    return `<div class="${containerClass}"${idAttr}>${itemsHtml}</div>`;
}
