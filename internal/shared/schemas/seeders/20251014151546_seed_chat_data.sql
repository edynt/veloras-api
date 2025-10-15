-- +goose Up
-- +goose StatementBegin
-- Insert chat conversations and messages
WITH conversation_data AS (
    INSERT INTO conversations DEFAULT VALUES RETURNING id
),
conversation_participants_data AS (
    INSERT INTO conversation_participants (conversation_id, user_id)
    SELECT 
        c.id,
        u.id
    FROM conversation_data c
    CROSS JOIN users u
    WHERE u.email IN ('john.doe@example.com', 'seller1@example.com')
    RETURNING conversation_id, user_id
),
chat_messages_data AS (
    INSERT INTO chat_messages (conversation_id, sender_id, content, type, is_read)
    SELECT 
        cp.conversation_id,
        cp.user_id,
        msg_data.content,
        msg_data.type,
        msg_data.is_read
    FROM conversation_participants_data cp
    CROSS JOIN (
        VALUES 
        ('Hello! I''m interested in this product. Is it still available?', 'text', false),
        ('Yes, it''s available! Would you like to know more details?', 'text', false),
        ('What''s the condition of the item?', 'text', false),
        ('It''s in excellent condition, barely used. I can send you more photos if needed.', 'text', false),
        ('That would be great! Can you also tell me about the warranty?', 'text', false),
        ('Sure! It comes with 1-year manufacturer warranty. I''ll send the photos now.', 'text', false),
        ('Perfect! I''ll take it. How do we proceed with payment?', 'text', false),
        ('Great! You can pay via Momo or bank transfer. I''ll send you the details.', 'text', false)
    ) AS msg_data(content, type, is_read)
    WHERE cp.user_id IN (
        SELECT id FROM users WHERE email IN ('john.doe@example.com', 'seller1@example.com')
    )
    RETURNING conversation_id, sender_id, content, type, is_read
)
SELECT 1; -- Dummy select to complete the CTE

-- Insert more conversations
WITH conversation_data2 AS (
    INSERT INTO conversations DEFAULT VALUES RETURNING id
),
conversation_participants_data2 AS (
    INSERT INTO conversation_participants (conversation_id, user_id)
    SELECT 
        c.id,
        u.id
    FROM conversation_data2 c
    CROSS JOIN users u
    WHERE u.email IN ('jane.smith@example.com', 'seller2@example.com')
    RETURNING conversation_id, user_id
),
chat_messages_data2 AS (
    INSERT INTO chat_messages (conversation_id, sender_id, content, type, is_read)
    SELECT 
        cp.conversation_id,
        cp.user_id,
        msg_data.content,
        msg_data.type,
        msg_data.is_read
    FROM conversation_participants_data2 cp
    CROSS JOIN (
        VALUES 
        ('Hi! I saw your product listing. Is the price negotiable?', 'text', false),
        ('Hello! Yes, I''m open to reasonable offers. What did you have in mind?', 'text', false),
        ('Would you consider 15% off the listed price?', 'text', false),
        ('That''s a bit low for me. How about 10% off?', 'text', false),
        ('Deal! When can I pick it up?', 'text', false),
        ('I can arrange pickup this weekend. I''ll send you the address.', 'text', false)
    ) AS msg_data(content, type, is_read)
    WHERE cp.user_id IN (
        SELECT id FROM users WHERE email IN ('jane.smith@example.com', 'seller2@example.com')
    )
    RETURNING conversation_id, sender_id, content, type, is_read
)
SELECT 1; -- Dummy select to complete the CTE
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
-- Delete all chat data
DELETE FROM chat_messages;
DELETE FROM conversation_participants;
DELETE FROM conversations;
-- +goose StatementEnd
