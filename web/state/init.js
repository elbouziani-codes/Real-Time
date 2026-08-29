import {fetchMe, fetchCategories, fetchPosts } from "./../fetcher.js"


export class State {
		constructor() {
			this.user = null;	
			this.socket = null;
			this.app = document.getElementById("app")
			this.feed = null;
			this.postsCollections = new PostsCollections;
			this.commentsCollections = null;
			this.messagesCollections = null;
		};

		async getUser() {
			const user = await fetchMe();		
			if (!user) return null
			this.user = new User(user, true)
		};

		async getCategories() {
			this.categories = await fetchCategories();	
		}

		async Ws() {
			 	
		}
		
}


export class PaginatedCollections {
		constructor() {
				this.items = []; 
				this.cursor = null;
				this.hasMore = true;	
		};	
		get() {
			return this.items;	
		}
		
} 




export class MessagesCollections extends PaginatedCollections {
		constructor(user) {
				this.user = user;
				super()	
		};
		async moreMessages() {
			if (!this.hasMore) return;
			const {items, nextCursor }= await fetchUsers(this.cursor, this.user);
			this.items.push(...items); 
			if (nextCursor == null) {
				this.hasMore = false;	
				return;
			};
			this.cursor = nextCursor;	
		};
}

export class PostsCollections extends PaginatedCollections {
		constructor(categories , liked = false) {		
				super();
				this.filters = null;
		};

		updateFilters(categories, liked) {
			this.items = [];
			this.cursor = null;
			this.hasMore = true;
			this.filters = {categories: categories, liked: liked};
		};

		async morePosts() {
			if (!this.hasMore) return;
			const {items, nextCursor} = await fetchPosts(this.cursor, this.filters);
			this.items.push(...items); 
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
		async moreUsers() {
			if (!this.hasMore) return;
			const {items, nextCursor}	 = await fetchUsers(this.cursor);
			this.items.push(...items);
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
		async moreComments() {
			if (!this.hasMore) return;
			const {items, nextCursor}	 = await fetchComments(this.postID, this.cursor);
			this.items.push(...items);
			if (nextCursor == null) {
				this.hasMore = false;	
				return;
			};
			this.cursor = nextCursor;
		};
};

export class User {
	constructor(user, me = false) {
		this.ID = user.ID;
		this.NickName = user.NickName;
		this.LastName = user.LastName;
		this.FirstName = user.FirstName;
		this.Online = me;
	};	
}
