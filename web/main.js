import {State, User} from "./state/init.js"

let CurrentState = null;

function main() {
	CurrentState = new State();	
	if (CurrentState.user == null) {
		// navigate to login	
	};	
}
