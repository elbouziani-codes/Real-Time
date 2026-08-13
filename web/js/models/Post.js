
export default function createPostModel({
    id = 0,
    title = '',
    content = '',
    category = '',
    author = '',
    likes = 0,
    comments = 0,
    createdAt = '',
} = {}) {
    return {
        id,
        title,
        content,
        category,
        author,
        likes,
        comments,
        createdAt,
    };
}
