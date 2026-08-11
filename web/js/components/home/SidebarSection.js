function renderTitle(title, titleClass) {
    if (!title) return '';

    const classAttribute = titleClass ? ` class="${titleClass}"` : '';
    return `<h4${classAttribute}>${title}</h4>`;
}


export default function SidebarSection({ title = '', titleClass = '', body = '' } = {}) {
    return `${renderTitle(title, titleClass)}${body}`;
}
