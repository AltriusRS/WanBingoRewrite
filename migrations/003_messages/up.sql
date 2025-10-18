-- Messages schema

CREATE TABLE messages
(
    id         VARCHAR(10) PRIMARY KEY,
    show_id    VARCHAR(10) REFERENCES shows (id) ON DELETE CASCADE       NOT NULL,
    player_id  VARCHAR(10) REFERENCES players (id) ON DELETE SET DEFAULT NOT NULL DEFAULT 'DELETED',
    contents   TEXT                                                      NOT NULL,
    system     BOOLEAN                                                   NOT NULL DEFAULT FALSE,
    replying   VARCHAR(10) references messages (id),
    created_at TIMESTAMP WITH TIME ZONE                                           DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE                                           DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP WITH TIME ZONE
);

CREATE INDEX idx_messages_show_id ON messages (show_id);
CREATE INDEX idx_messages_player_id ON messages (player_id);

CREATE TRIGGER update_messages_updated_at
    BEFORE UPDATE
    ON messages
    FOR EACH ROW
EXECUTE FUNCTION update_updated_at_column();

COMMENT ON TABLE messages IS 'Records all the messages sent during a given show cycle';