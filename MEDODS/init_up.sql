CREATE TABLE users(id BIGSERIAL PRIMARY KEY,guid text,email text);
CREATE TABLE refresh_tokens(
    user_id BIGINT PRIMARY KEY REFERENCES users(id) ON DELETE CASCADE UNIQUE,
    refresh_token text,
    available boolean,
    created_at time);

INSERT INTO users(guid,email)
VALUES ('b758e00b-9c42-4f2b-a84d-0af47b937d17','test@gmail.com');
INSERT INTO refresh_tokens(user_id,refresh_token,available)
VALUES (1,'',false)