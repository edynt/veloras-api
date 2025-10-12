-- name: CreateConversation :one
INSERT INTO conversations DEFAULT VALUES RETURNING *;

-- name: AddConversationParticipant :one
INSERT INTO conversation_participants (conversation_id, user_id)
VALUES ($1, $2)
ON CONFLICT (conversation_id, user_id) DO NOTHING
RETURNING *;

-- name: GetConversation :one
SELECT c.*, 
       array_agg(cp.user_id) as participant_ids
FROM conversations c
LEFT JOIN conversation_participants cp ON c.id = cp.conversation_id
WHERE c.id = $1
GROUP BY c.id, c.created_at, c.updated_at;

-- name: GetConversationByParticipants :one
SELECT c.*
FROM conversations c
WHERE c.id IN (
    SELECT conversation_id 
    FROM conversation_participants 
    WHERE user_id = ANY($1::integer[])
    GROUP BY conversation_id 
    HAVING COUNT(DISTINCT user_id) = array_length($1::integer[], 1)
)
LIMIT 1;

-- name: ListUserConversations :many
SELECT c.*, 
       array_agg(cp.user_id) as participant_ids,
       cm.content as last_message_content,
       cm.created_at as last_message_at,
       cm.sender_id as last_message_sender_id
FROM conversations c
LEFT JOIN conversation_participants cp ON c.id = cp.conversation_id
LEFT JOIN LATERAL (
    SELECT content, created_at, sender_id
    FROM chat_messages 
    WHERE conversation_id = c.id 
    ORDER BY created_at DESC 
    LIMIT 1
) cm ON true
WHERE c.id IN (
    SELECT conversation_id 
    FROM conversation_participants 
    WHERE conversation_participants.user_id = $1
)
GROUP BY c.id, c.created_at, c.updated_at, cm.content, cm.created_at, cm.sender_id
ORDER BY COALESCE(cm.created_at, c.created_at) DESC
LIMIT $2 OFFSET $3;

-- name: CreateChatMessage :one
INSERT INTO chat_messages (
    conversation_id, sender_id, content, type, attachments
) VALUES (
    $1, $2, $3, $4, $5
) RETURNING *;

-- name: ListChatMessages :many
SELECT cm.*, 
       u.first_name || ' ' || u.last_name as sender_name,
       u.email as sender_email
FROM chat_messages cm
LEFT JOIN users u ON cm.sender_id = u.id
WHERE cm.conversation_id = $1
ORDER BY cm.created_at ASC
LIMIT $2 OFFSET $3;

-- name: MarkMessagesAsRead :exec
UPDATE chat_messages 
SET is_read = true
WHERE conversation_id = $1 AND sender_id != $2;

-- name: GetUnreadMessageCount :one
SELECT COUNT(*) as unread_count
FROM chat_messages cm
JOIN conversation_participants cp ON cm.conversation_id = cp.conversation_id
WHERE cp.user_id = $1 AND cm.sender_id != $1 AND cm.is_read = false;

-- name: GetConversationUnreadCount :one
SELECT COUNT(*) as unread_count
FROM chat_messages cm
JOIN conversation_participants cp ON cm.conversation_id = cp.conversation_id
WHERE cp.user_id = $1 AND cm.conversation_id = $2 AND cm.sender_id != $1 AND cm.is_read = false;

-- name: DeleteConversation :exec
DELETE FROM conversations WHERE id = $1;

-- name: RemoveConversationParticipant :exec
DELETE FROM conversation_participants 
WHERE conversation_id = $1 AND user_id = $2;

-- name: GetConversationWithDetails :one
SELECT c.*, 
       array_agg(
           json_build_object(
               'user_id', cp.user_id,
               'joined_at', cp.created_at
           )
       ) as participants
FROM conversations c
LEFT JOIN conversation_participants cp ON c.id = cp.conversation_id
WHERE c.id = $1
GROUP BY c.id, c.created_at, c.updated_at;

-- name: UpdateConversation :one
UPDATE conversations 
SET updated_at = NOW()
WHERE id = $1
RETURNING *;

-- name: AddParticipant :one
INSERT INTO conversation_participants (conversation_id, user_id)
VALUES ($1, $2)
ON CONFLICT (conversation_id, user_id) DO NOTHING
RETURNING *;

-- name: RemoveParticipant :exec
DELETE FROM conversation_participants 
WHERE conversation_id = $1 AND user_id = $2;

-- name: GetConversationParticipants :many
SELECT cp.*, u.first_name || ' ' || u.last_name as user_name, u.email as user_email
FROM conversation_participants cp
LEFT JOIN users u ON cp.user_id = u.id
WHERE cp.conversation_id = $1;

-- name: GetChatMessageWithDetails :one
SELECT cm.*, 
       u.first_name || ' ' || u.last_name as sender_name,
       u.email as sender_email
FROM chat_messages cm
LEFT JOIN users u ON cm.sender_id = u.id
WHERE cm.id = $1;

-- name: ListConversationMessages :many
SELECT cm.*, 
       u.first_name || ' ' || u.last_name as sender_name,
       u.email as sender_email
FROM chat_messages cm
LEFT JOIN users u ON cm.sender_id = u.id
WHERE cm.conversation_id = $1
ORDER BY cm.created_at ASC
LIMIT $2 OFFSET $3;

-- name: MarkMessageAsRead :exec
UPDATE chat_messages 
SET is_read = true
WHERE id = $1;

-- name: MarkConversationAsRead :exec
UPDATE chat_messages 
SET is_read = true
WHERE conversation_id = $1 AND sender_id != $2;

-- name: DeleteChatMessage :exec
DELETE FROM chat_messages WHERE id = $1;

-- name: GetChatStats :one
SELECT 
    COUNT(DISTINCT c.id) as total_conversations,
    COUNT(cm.id) as total_messages,
    COUNT(CASE WHEN cm.is_read = false THEN 1 END) as unread_messages
FROM conversations c
LEFT JOIN chat_messages cm ON c.id = cm.conversation_id;