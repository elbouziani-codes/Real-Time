const Category = {
    
}
export default function createCategory({
    id = 0,
    name = '',
    icon = '',
    colorClass = '',
} = {}) {
    return {
        id,
        name,
        icon,
        colorClass,
    };
}
