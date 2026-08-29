import {initFeedListeners} from "./../js/listeners/home.js" 
import homePage  from "./../js/pages/home.js" 
import Feed from "./../js/components/home/Feed.js";
import {state} from "/main.js";




export async function handleFeed() {
	await state.postsCollections.morePosts() 
	console.log(state.postsCollections.get())
	state.app.innerHTML = await homePage(state);
	state.feed = document.getElementById("feed")
	initFeedListeners();
};

export async function refreshFeed() {
	state.feed.innerHTML = await Feed(state.postsCollections.get());
	initFeedListeners();
};
