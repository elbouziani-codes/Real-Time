import User from './User.js';

const LIST_CLASSES = {
    item: 'users-list',
    message: 'last-messages-list',
};


export default function UserList(users = [], { variant = 'item', listClass = '' } = {}) {
    const containerClass = listClass || LIST_CLASSES[variant] || LIST_CLASSES.item;
    const itemsHtml = users.map((user) => User(user, { variant })).join('');

    return `<div class="${containerClass}">${itemsHtml}</div>`;
}
