

export default async function fetchCategories() {

    try{
        const response = await fetch("/api/categories", {
            method: "GET",
        });

        const categories = await response.json();
        const code = response.ok
        return {code:code , body:categories}
    }catch(error){
        (error)
        return {code:500 , body:"Error in request"}
    }
}