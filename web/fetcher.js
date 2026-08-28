


export async function fetchUsers(cursor = null) {
	const params = new URLSearchParams();
	try {
		if (cursor)  params.append("cursor", cursor);
		const repoonse  = await fetch(`/api/users?${params.toString()}`, {method: "GET"});		
		if (!response.ok) return []
		const users = await response.json();
		return {users, users.at(-1).ID};	
	} catch(error) {
		return {[], null};
	};
};

export async function fetchCategories() {
	try {
		const repoonse  = await fetch(`/api/categories`, {method: "GET"});		
		if (!response.ok) return []
		const categories = await response.json();
		return categories;	
	} catch(error) {
		return {[], null};
	};
};

export async function fetchPosts(cursor = null) {
	const params = new URLSearchParams();
	try {
		if (cursor)  params.append("cursor", cursor);
		const repoonse  = await fetch(`/api/posts?${params.toString()}`, {method: "GET"});		
		if (!response.ok) return [];
		const posts = await response.json();
		return {posts, posts.at(-1).ID};	
	} catch(error) {
		return {[], null};
	};
};


export async function fetchComments(postID, cursor = null) {
	const params = new URLSearchParams({postID: postID});
	try {
		if (cursor)  params.append("cursor", cursor);
		const repoonse  = await fetch(`/api/posts?${params.toString()}`, {method: "GET"});		
		if (!response.ok) return [];
		const comments = await response.json();
		return {comments, comments.at(-1).ID};	
	} catch(error) {
		return {[], null};
	};
};


export async function fetchMessages(userID, cursor = null) {
	const params = new URLSearchParams({userID: userID});
	try {
		if (cursor)  params.append("cursor", cursor);
		const repoonse  = await fetch(`/api/posts?${params.toString()}`, {method: "GET"});		
		if (!response.ok) return [];
		const comments = await response.json();
		return {comments, comments.at(-1).ID};	
	} catch(error) {
		return {[], null};
	};
};
