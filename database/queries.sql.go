package db 


import (
	"context"
)



const creatChatQuery = `
	INSERT into conversations (
		id 
	) VALUES (
		$1
);
`
const getChatIDQuery = `
	SELECT cp1.conversation_id 		
	FROM conversation_participants cp1
	JOIN conversation_participants cp2
	ON cp1.conversation_id = cp2.conversation_id
		WHERE cp1.user_id = ? 
		AND cp2.user_id = ?  
`


func (q *Queries) createChat(id []byte, ctx context.Context) (err)  {
	err := q.ExecContext(ctx, createChatQuery, id)	
	return err
}  


