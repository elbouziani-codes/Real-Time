import createCategory from '../models/Category.js';
import fetchCategories from './../api/categories.js';


export let DEFAULT_CATEGORIES = [];


export async function seedCategories(){
    let response = await fetchCategories()

    if (response.code == true){
        if (!Array.isArray(response.body)) {
                DEFAULT_POST = ["not fond categories"];
        }else{
            DEFAULT_CATEGORIES = response.body
        }
    }else{
        DEFAULT_CATEGORIES = ["error in fetch Categories"]
    }
}