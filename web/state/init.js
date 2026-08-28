


export class State {
		constructor() {
			this.user = new User();		
			this.postsCursor = null;
			this.messagesCursor = null;
		};

		render() {
				
		}
}


export class PaginatedCollection {
		constructor(fetchPage) {
				this.fetchPage = fetchPage; 		
				this.items = []; 
				this.cursor = [];
				this.hasMore = true;	
		};	
		await fetchMore() {
			if (!this.hasMore) return;
			{items, nextCursor}	 = await this.fetchPage(this.cursor);
			this.items = this.items.push(...items); 
			if (nextCursor == null) {
				this.hasMore = false;	
				return;
			};
			this.cursor = nextCursor;
		};
} 


export class ChatState {
		constructor() {
			this.users = await fetchUsers();	
			this.queue = [];
			this.cursor = null;
			if (this.users.length > 0 ) {
				this.cursor = this.users[this.users.length-1]; 		
			}	
			this.SelectedChat = null;
		};
		MoreUsers() {
			if (this.cursor == null) {
				return;	
			}	
			const users = await fetchUsers(this.cursor);
			this.users.push(users);	
		};
}

export class MessagesCollections extends PaginatedCollections {
		constructor(user) {
				this.user = user;
				super()	
		};
		moreMessages() {
			if (!this.hasMore) return;
			{items, nextCursor}	 = await fetchUsers(this.cursor, this.user);
			this.items = this.items.push(...items); 
			if (nextCursor == null) {
				this.hasMore = false;	
				return;
			};
			this.cursor = nextCursor;	
		};
}

export class PostsCollections extends PaginatedCollections {
		constructor() {	
				super()	
		};
		morePosts() {
			if (!this.hasMore) return;
			{items, nextCursor}	 = await fetchPosts(this.cursor);
			this.items = this.items.push(...items); 
			if (nextCursor == null) {
				this.hasMore = false;	
				return;
			};
			this.cursor = nextCursor;	
		};
};

export class UsersCollections extends PaginatedCollections {
		constructor() {	
				super()	
				this.queue = [];
		};
		moreUsers() {
			if (!this.hasMore) return;
			{items, nextCursor}	 = await fetchUsers(this.cursor);
			this.items = this.items.push(...items); 
			// here I must update online queue// 
			if (nextCursor == null) {
				this.hasMore = false;	
				return;
			};
			this.cursor = nextCursor;	
		};
}

export class CommentsCollections extends PaginatedCollections {
		constructor(postID) {	
				this.postID = postID;
				super()	
		};
		moreComments() {
			if (!this.hasMore) return;
			{items, nextCursor}	 = await fetchComments(this.cursor, this.postID);
			this.items = this.items.push(...items); 
			if (nextCursor == null) {
				this.hasMore = false;	
				return;
			};
			this.cursor = nextCursor;	
		};
};

export class User {
	constructor(me = false) {
		let user = await fetchUser(me);  	
		this.id = user.id
		this.nickName = user.nickName;
		this.lastName = user.lastName;
		this.firstName = user.firstName;
		this.online = me;
	};	
}
