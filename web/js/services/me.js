import {fetchMe, fetchLogout} from "./../api/auth.js"
import { navigate } from "../router/router.js";
import { resetPosts } from "./posts.js";


export let me = {}


export default async function CreateMe(path){
    let response = await fetchMe()
    if (response.code == 200) {
        me = {...response.body}
        if (path == "/login" || path == "/register" ){
            navigate("/");
        }
    }else{
        if (path != "/login" && path != "/register" ){
            navigate("/login");
        }
    }
}

export async function logoutMe(){
    let response = await fetchLogout()
    if (response.code != 200 && response.code != 401) {
        console.log("logout failed:", response.body)
        return false
    }
    me = {}
    resetPosts()
    return true
}

