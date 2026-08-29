import {State, User} from "./state/init.js"
import { handleAuth} from "./handlers/login.js"
import { handleFeed } from "./handlers/feed.js"


export let state = null;

async function main() {
	state = new State();	
	await state.getUser(); 

	if (state.user == null) {
	 	handleAuth();		
	} else {
		await state.getCategories();
		await state.Ws();
		handleFeed(state)	
	};	
}

await main();
