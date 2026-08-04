const Category = {
    
}
export default function createCategory({
    id = 0,
    name = '',
    icon = '',
    count = 0,
    colorClass = '',
} = {}) {
    return {
        id,
        name,
        icon,
        count,
        colorClass,
    };
}
