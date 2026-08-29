


export async function fetchUsers(cursor = null) {
	const params = new URLSearchParams();
	try {
		if (cursor)  params.append("cursor", cursor);
		const response  = await fetch(`/api/users?${params.toString()}`, {method: "GET"});		
		if (!response.ok) return {items: [], nextCursor: null}
		const users = await response.json();
		return {items: users, nextCursor: users.length ? users.at(-1).ID : null};	
	} catch(error) {
		return {items: [], nextCursor: null};
	};
};

export async function reactToPost(postID, isLike) {
	const body = JSON.stringify({
		parent_id: postID,
		is_like: isLike
	});
	try {
		const response  = await fetch(`/api/reactions`, {method: "POST", body: body});		
		if (!response.ok) return false; 
		return true;	
	} catch(error) {
		console.log(error)
		return false;
	};
};

export async function fetchCategories() {
	try {
		const response  = await fetch(`/api/categories`, {method: "GET"});		
		if (!response.ok) return []
		const categories = await response.json();
		return categories;	
	} catch(error) {
		return [];
	};
};


export async function fetchMe() {
	try {
		const response  = await fetch(`/api/me`, {method: "GET"});		
		if (!response.ok) return null;
		const me = await response.json();
		return me;	
	} catch(error) {
		return null;
	};
};



export async function fetchPosts(cursor = null, filters) {
	const params = new URLSearchParams();
	console.log(filters)
	try {
		if (cursor)  params.append("cursor", cursor);
		if (filters) {
			filters.categories.forEach((c) => {
				params.append("category", c);	
			});	
			if (filters.liked) {
				params.append("liked", "true");	
			}
			
		}
		const response  = await fetch(`/api/posts?${params.toString()}`, {method: "GET"});		
		if (!response.ok) return {items: [], nextCursor: null};
		const posts = await response.json();
		
		return {items: posts, nextCursor: posts.length ? posts.at(-1).ID : null};	
	} catch(error) {
		return {items: [], nextCursor: null};
	};
};


export async function fetchComments(postID, cursor = null) {
	const params = new URLSearchParams({postID: postID});
	try {
		if (cursor)  params.append("cursor", cursor);
		const response  = await fetch(`/api/comments?${params.toString()}`, {method: "GET"});		
		if (!response.ok) return {items: [], nextCursor: null};
		const comments = await response.json();
		return {items: comments, nextCursor: comments.length ? comments.at(-1).ID : null};
	} catch(error) {
		return {items: [], nextCursor: null};
	};
};


export async function fetchMessages(userID, cursor = null) {
	const params = new URLSearchParams({userID: userID});
	try {
		if (cursor)  params.append("cursor", cursor);
		const response  = await fetch(`/api/messages?${params.toString()}`, {method: "GET"});		
		if (!response.ok) return {items: [], nextCursor: null};
		const messages = await response.json();
		return {items: messages, nextCursor: messages.length ? messages.at(-1).ID : null};
	} catch(error) {
		return {items: [], nextCursor: null};
	};
};
