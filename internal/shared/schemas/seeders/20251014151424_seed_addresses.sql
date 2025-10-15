-- +goose Up
-- +goose StatementBegin
-- Insert addresses for users
INSERT INTO addresses (user_id, street, city, state, zip_code, country, is_default)
SELECT 
    u.id,
    address_data.street,
    address_data.city,
    address_data.state,
    address_data.zip_code,
    address_data.country,
    address_data.is_default
FROM users u
CROSS JOIN (
    VALUES 
    -- Regular users addresses
    ('123 Le Loi Street', 'Ho Chi Minh City', 'Ho Chi Minh', '700000', 'Vietnam', true),
    ('456 Nguyen Hue Boulevard', 'Ho Chi Minh City', 'Ho Chi Minh', '700000', 'Vietnam', true),
    ('789 Dong Khoi Street', 'Ho Chi Minh City', 'Ho Chi Minh', '700000', 'Vietnam', true),
    ('321 Pasteur Street', 'Ho Chi Minh City', 'Ho Chi Minh', '700000', 'Vietnam', true),
    ('654 Vo Van Tan Street', 'Ho Chi Minh City', 'Ho Chi Minh', '700000', 'Vietnam', true),
    ('987 Cach Mang Thang 8', 'Ho Chi Minh City', 'Ho Chi Minh', '700000', 'Vietnam', true),
    ('147 Hai Ba Trung Street', 'Ho Chi Minh City', 'Ho Chi Minh', '700000', 'Vietnam', true),
    ('258 Ly Tu Trong Street', 'Ho Chi Minh City', 'Ho Chi Minh', '700000', 'Vietnam', true),
    ('369 Tran Hung Dao Street', 'Ho Chi Minh City', 'Ho Chi Minh', '700000', 'Vietnam', true),
    ('741 Nguyen Thi Minh Khai', 'Ho Chi Minh City', 'Ho Chi Minh', '700000', 'Vietnam', true),
    
    -- Seller addresses
    ('123 Tech Street, District 1', 'Ho Chi Minh City', 'Ho Chi Minh', '700000', 'Vietnam', true),
    ('456 Fashion Avenue, District 3', 'Ho Chi Minh City', 'Ho Chi Minh', '700000', 'Vietnam', true),
    ('789 Garden Road, District 7', 'Ho Chi Minh City', 'Ho Chi Minh', '700000', 'Vietnam', true),
    ('321 Sports Plaza, District 2', 'Ho Chi Minh City', 'Ho Chi Minh', '700000', 'Vietnam', true),
    ('654 Book Lane, District 5', 'Ho Chi Minh City', 'Ho Chi Minh', '700000', 'Vietnam', true),
    
    -- Admin addresses
    ('999 Admin Tower, District 1', 'Ho Chi Minh City', 'Ho Chi Minh', '700000', 'Vietnam', true),
    ('888 Super Admin Plaza, District 1', 'Ho Chi Minh City', 'Ho Chi Minh', '700000', 'Vietnam', true),
    
    -- Vietnamese users addresses
    ('111 Nguyen Trai Street', 'Ho Chi Minh City', 'Ho Chi Minh', '700000', 'Vietnam', true),
    ('222 Le Van Sy Street', 'Ho Chi Minh City', 'Ho Chi Minh', '700000', 'Vietnam', true),
    ('333 Pham Van Dong Street', 'Ho Chi Minh City', 'Ho Chi Minh', '700000', 'Vietnam', true),
    ('444 Truong Chinh Street', 'Ho Chi Minh City', 'Ho Chi Minh', '700000', 'Vietnam', true),
    ('555 Hoang Van Thu Street', 'Ho Chi Minh City', 'Ho Chi Minh', '700000', 'Vietnam', true),
    
    -- More users addresses
    ('666 Ba Trieu Street', 'Ho Chi Minh City', 'Ho Chi Minh', '700000', 'Vietnam', true),
    ('777 Dien Bien Phu Street', 'Ho Chi Minh City', 'Ho Chi Minh', '700000', 'Vietnam', true),
    ('888 Xo Viet Nghe Tinh', 'Ho Chi Minh City', 'Ho Chi Minh', '700000', 'Vietnam', true),
    ('999 Nguyen Van Cu Street', 'Ho Chi Minh City', 'Ho Chi Minh', '700000', 'Vietnam', true),
    ('1010 3 Thang 2 Street', 'Ho Chi Minh City', 'Ho Chi Minh', '700000', 'Vietnam', true)
) AS address_data(street, city, state, zip_code, country, is_default)
WHERE u.email IN (
    'john.doe@example.com', 'jane.smith@example.com', 'mike.wilson@example.com', 'sarah.johnson@example.com', 'david.brown@example.com',
    'lisa.davis@example.com', 'robert.miller@example.com', 'emily.garcia@example.com', 'james.martinez@example.com', 'jennifer.anderson@example.com',
    'seller1@example.com', 'seller2@example.com', 'seller3@example.com', 'seller4@example.com', 'seller5@example.com',
    'admin@example.com', 'superadmin@example.com',
    'nguyen.van.a@example.com', 'tran.thi.b@example.com', 'le.van.c@example.com', 'pham.thi.d@example.com', 'hoang.van.e@example.com',
    'alice.johnson@example.com', 'bob.smith@example.com', 'carol.williams@example.com', 'daniel.jones@example.com', 'eve.brown@example.com'
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
-- Delete all addresses
DELETE FROM addresses;
-- +goose StatementEnd
