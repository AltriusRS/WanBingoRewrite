-- WAN Bingo Database Schema

-- Drop tables if they exist (in reverse dependency order)
DROP TABLE IF EXISTS messages CASCADE;
DROP TABLE IF EXISTS tile_confirmations CASCADE;
DROP TABLE IF EXISTS boards CASCADE;
DROP TABLE IF EXISTS show_tiles CASCADE;
DROP TABLE IF EXISTS tiles CASCADE;
DROP TABLE IF EXISTS shows CASCADE;
DROP TABLE IF EXISTS sessions CASCADE;
DROP TABLE IF EXISTS players CASCADE;

-- Players table
CREATE TABLE players
(
    id           VARCHAR(10) PRIMARY KEY,
    did          VARCHAR(20) UNIQUE NOT NULL,
    display_name VARCHAR(100)       NOT NULL,
    avatar       TEXT,
    settings     JSONB                    DEFAULT '{}'::jsonb,
    score        INTEGER                  DEFAULT 0,
    created_at   TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at   TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    deleted_at   TIMESTAMP WITH TIME ZONE
);

-- Create index on Discord ID for faster lookups
CREATE INDEX idx_players_did ON players (did);

-- Sessions table
CREATE TABLE sessions
(
    id         VARCHAR(32) PRIMARY KEY,
    player_id  VARCHAR(10) REFERENCES players (id) ON DELETE CASCADE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    expires_at TIMESTAMP WITH TIME ZONE NOT NULL,
    deleted_at TIMESTAMP WITH TIME ZONE
);

-- Create index on player_id for faster lookups
CREATE INDEX idx_sessions_player_id ON sessions (player_id);
-- Create index on expires_at for cleanup queries
CREATE INDEX idx_sessions_expires_at ON sessions (expires_at);

-- Shows table
CREATE TABLE shows
(
    id                VARCHAR(10) PRIMARY KEY,
    youtube_id        VARCHAR(20) UNIQUE,
    scheduled_time    TIMESTAMP WITH TIME ZONE,
    actual_start_time TIMESTAMP WITH TIME ZONE,
    thumbnail         TEXT,
    metadata          JSONB                    DEFAULT '{}'::jsonb,
    created_at        TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at        TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    deleted_at        TIMESTAMP WITH TIME ZONE
);

-- Create index on scheduled_time for querying upcoming/past shows
CREATE INDEX idx_shows_scheduled_time ON shows (scheduled_time);
CREATE INDEX idx_shows_youtube_id ON shows (youtube_id);

-- Tiles table
CREATE TABLE tiles
(
    id         VARCHAR(10) PRIMARY KEY,
    title      VARCHAR(200) NOT NULL,
    category   VARCHAR(50),
    last_drawn TIMESTAMP WITH TIME ZONE,
    weight     FLOAT8       NOT NULL    DEFAULT 0.30 + floor(random() * 36) * 0.02,
    score      FLOAT8       NOT NULL    DEFAULT floor(random() * 46 + 5),
    created_by VARCHAR(10)  REFERENCES players (id) ON DELETE SET NULL,
    settings   JSONB                    DEFAULT '{}'::jsonb,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP WITH TIME ZONE
);

-- Create indexes for common queries
CREATE INDEX idx_tiles_category ON tiles (category);
CREATE INDEX idx_tiles_created_by ON tiles (created_by);
CREATE INDEX idx_tiles_last_drawn ON tiles (last_drawn);

-- Show Tiles junction table
CREATE TABLE show_tiles
(
    show_id    VARCHAR(10) REFERENCES shows (id) ON DELETE CASCADE,
    tile_id    VARCHAR(10) REFERENCES tiles (id) ON DELETE CASCADE,
    weight     float8                   DEFAULT 1,
    score      float8                   DEFAULT 0,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP WITH TIME ZONE,
    PRIMARY KEY (show_id, tile_id)
);

-- Create indexes for common queries
CREATE INDEX idx_show_tiles_show_id ON show_tiles (show_id);
CREATE INDEX idx_show_tiles_tile_id ON show_tiles (tile_id);

-- Boards table
CREATE TABLE boards
(
    id                      VARCHAR(10) PRIMARY KEY,
    player_id               VARCHAR(10) REFERENCES players (id) ON DELETE CASCADE,
    show_id                 VARCHAR(10) REFERENCES shows (id) ON DELETE CASCADE,
    tiles                   TEXT[] NOT NULL,
    winner                  BOOLEAN                  DEFAULT FALSE,
    total_score             float8                   DEFAULT 0,
    potential_score         float8                   DEFAULT 0,
    regeneration_diminisher float8                   DEFAULT 0,
    created_at              TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at              TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    deleted_at              TIMESTAMP WITH TIME ZONE,
    UNIQUE (player_id, show_id)
);

-- Create indexes for common queries
CREATE INDEX idx_boards_player_id ON boards (player_id);
CREATE INDEX idx_boards_show_id ON boards (show_id);
CREATE INDEX idx_boards_winner ON boards (winner);

-- Tile Confirmations table
CREATE TABLE tile_confirmations
(
    id                VARCHAR(10) PRIMARY KEY,
    show_id           VARCHAR(10) REFERENCES shows (id) ON DELETE CASCADE,
    tile_id           VARCHAR(10) REFERENCES tiles (id) ON DELETE CASCADE,
    confirmed_by      VARCHAR(10) REFERENCES players (id) ON DELETE SET NULL,
    context           TEXT,
    confirmation_time TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    created_at        TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at        TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    deleted_at        TIMESTAMP WITH TIME ZONE
);

-- Create indexes for common queries
CREATE INDEX idx_tile_confirmations_show_id ON tile_confirmations (show_id);
CREATE INDEX idx_tile_confirmations_tile_id ON tile_confirmations (tile_id);
CREATE INDEX idx_tile_confirmations_confirmed_by ON tile_confirmations (confirmed_by);
CREATE INDEX idx_tile_confirmations_time ON tile_confirmations (confirmation_time);

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

-- Create a function to update the updated_at timestamp
CREATE OR REPLACE FUNCTION update_updated_at_column()
    RETURNS TRIGGER AS
$$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ language 'plpgsql';

-- Create triggers to auto-update updated_at timestamps
CREATE TRIGGER update_players_updated_at
    BEFORE UPDATE
    ON players
    FOR EACH ROW
EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_sessions_updated_at
    BEFORE UPDATE
    ON sessions
    FOR EACH ROW
EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_shows_updated_at
    BEFORE UPDATE
    ON shows
    FOR EACH ROW
EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_tiles_updated_at
    BEFORE UPDATE
    ON tiles
    FOR EACH ROW
EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_show_tiles_updated_at
    BEFORE UPDATE
    ON show_tiles
    FOR EACH ROW
EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_boards_updated_at
    BEFORE UPDATE
    ON boards
    FOR EACH ROW
EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_tile_confirmations_updated_at
    BEFORE UPDATE
    ON tile_confirmations
    FOR EACH ROW
EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_messages_updated_at
    BEFORE UPDATE
    ON messages
    FOR EACH ROW
EXECUTE FUNCTION update_updated_at_column();

-- Add comments to document the schema
COMMENT ON TABLE players IS 'Stores user account information';
COMMENT ON TABLE sessions IS 'Stores user session tokens for authentication';
COMMENT ON TABLE shows IS 'Stores information about WAN show episodes';
COMMENT ON TABLE tiles IS 'Stores bingo tile definitions';
COMMENT ON TABLE show_tiles IS 'Junction table linking tiles to specific shows with dynamic scoring';
COMMENT ON TABLE boards IS 'Stores player bingo boards for each show';
COMMENT ON TABLE tile_confirmations IS 'Records when tiles are confirmed during a show';
COMMENT ON TABLE messages IS 'Records all the messages sent during a given show cycle';

COMMENT ON COLUMN players.did IS 'Discord user ID from OAuth authentication';
COMMENT ON COLUMN players.settings IS 'JSON object containing user preferences (chat_name_color, pronouns, sound_on_mention, interface_language, etc)';
COMMENT ON COLUMN tiles.settings IS 'JSON object containing tile configuration (needs_context, has_timer, timer_duration, etc)';
COMMENT ON COLUMN shows.metadata IS 'JSON object containing additional show data (hosts, sponsors, duration, etc)';
COMMENT ON COLUMN boards.tiles IS 'Array of tile IDs representing the player''s bingo board layout';


INSERT INTO public.shows (id, youtube_id, scheduled_time, actual_start_time, thumbnail, metadata, created_at,
                          updated_at, deleted_at)
VALUES ('Y2kz75uBC8', 'YVHXYqMPyzc', '2025-10-11 00:30:00.000000 +00:00', '2025-10-11 00:05:06.000000 +00:00',
        'https://pbs.floatplane.com/stream_thumbnails/5c13f3c006f1be15e08e05c0/733054221374526_1760139634263.jpeg', '{
    "title": "Piracy Is Dangerous And Harmful",
    "fp_vod": "w3A5fKcfTi"
  }', '2025-10-10 23:46:40.000000 +00:00', '2025-10-16 19:12:58.692094 +00:00', null);



-- Insert the "deleted user" placeholder account
-- This account is used to attribute data from deleted users
INSERT INTO players (id, did, display_name, avatar, settings, score, created_at, updated_at)
VALUES ('DELETED',
        '0',
        '[Deleted User]',
        NULL,
        '{
          "system_account": true
        }'::jsonb,
        0,
        CURRENT_TIMESTAMP,
        CURRENT_TIMESTAMP),
       ('SYSTEM',
        '1',
        '[System]',
        NULL,
        '{
          "system_account": true
        }'::jsonb,
        0,
        CURRENT_TIMESTAMP,
        CURRENT_TIMESTAMP)
;



INSERT INTO tiles (id, title, category, last_drawn)
VALUES ('pYhro7iTSQ', 'Linus or Luke or Dan sighs', 'Events', NOW()),
       ('3AXPH39mmU', 'Linus ignores Luke to change the topic', 'Linus', NOW()),
       ('kJdmP6gDSi', 'Linus hates on Twitch Chat', 'Linus', NOW()),
       ('ji3eCBLPrx', 'Linus Facepalms', 'Linus', NOW()),
       ('AQrAcWkqsg', 'Intro/Outro run accidentally', 'Set/Production', NOW()),
       ('W9No_ylM54', 'The microphone gets hit', 'Set/Production', NOW()),
       ('OPdeRADtmK', '*Wow, I feel old...*', 'Events', NOW()),
       ('jXurzTPKY_', 'Camera Not Focused', 'Set/Production', NOW()),
       ('GdOM4Jx7ZI', 'Luke was Wrong', 'Luke', NOW()),
       ('WH6tzRJ5PG', 'Linus Was Wrong', 'Linus', NOW()),
       ('U3SEvAZT4E', 'Colton Quit / Fired joke', 'Events', NOW()),
       ('kxR5uknbsW', 'Linus Drops Something', 'Linus', NOW()),
       ('LdaAhl25W4', 'Nvidia News!', 'Topics', NOW()),
       ('BaD8f-Qmpg', 'AMD News!', 'Topics', NOW()),
       ('RubFqyA_Zs', 'Intel News!', 'Topics', NOW()),
       ('lovN9OEyXu', 'Apple News!', 'Topics', NOW()),
       ('1zVBE4wLA9', 'New Sponsor!', 'Sponsors', NOW()),
       ('3LHG-WTzBz', 'Screenshare has No Audio', 'Set/Production', NOW()),
       ('fY2u7DoLZ7', 'Audio Clipping', 'Set/Production', NOW()),
       ('sBbnG8b1O_', 'Literally one super topic until sponsor spot', 'Sponsors', NOW()),
       ('nc9JZl99ME', 'Linus leaves the other host alone', 'Linus', NOW()),
       ('JAgGIaNECU', 'Video output not connected to laptop', 'Set/Production', NOW()),
       ('lND6zOabtj', 'Console Topic for the peasantry', 'Topics', NOW()),
       ('eb8uFT5T_K', 'Luke says *That''s Hilarious!*', 'Luke', NOW()),
       ('17U6OdT6QW', 'Someone messes with the set', 'Set/Production', NOW()),
       ('7zeqP0zUVd', 'Linus: ''We''ve got a great show for you today!''', 'Linus', NOW()),
       ('NZUFuGT4dg', 'No actual news before sponsor spot', 'Sponsors', NOW()),
       ('uRE41fjRbK', 'Linus doesn''t censor while swearing', 'Linus', NOW()),
       ('cjAULhjZBL', 'Motion-Sickness Camera', 'Set/Production', NOW()),
       ('69LrHZvYtE', 'Super-Zoomed Camera', 'Set/Production', NOW()),
       ('MdM9kwiAna', 'Hello Dan', 'Dan', NOW()),
       ('i27fmnr_tu', 'LTT Store Plug', 'Events', NOW()),
       ('u5jjt2ruMP', 'LTT Water Bottle drank from', 'Events', NOW()),
       ('JDJuwnV-OI', 'Banana For Scale', 'Events', NOW()),
       ('ftsUAzY_Db', 'Linus Hot Take', 'Linus', NOW()),
       ('A6NucRK2Gj', 'Linus Roasts a Company', 'Linus', NOW()),
       ('kJvJpGqFBV', 'Linus'' parenting stories', 'Linus', NOW()),
       ('KsG42jyeJu', 'Mispronunciation of a word/phrase', 'Events', NOW()),
       ('I3GA3iJPjt', 'Stream Dies', 'Events', NOW()),
       ('tj69a_JZfe', 'Costumes!', 'Events', NOW()),
       ('BfaqFYztlR', '4+ Hour WAN Show (What a Champ!)', 'Events', NOW()),
       ('ASPY0Qy1tZ', 'Luke laughs uncomfortably loud', 'Luke', NOW()),
       ('2wN544Mxvc', 'Luke talks about AI', 'Luke', NOW()),
       ('KKbzU2rbiJ', 'Accidentally shared confidential info', 'Events', NOW()),
       ('W-motMP4Zs', 'Dan ignores Linus', 'Dan', NOW()),
       ('lqnGNVc9nt', 'Hit me Dan', 'Dan', NOW()),
       ('c7BjRI697F', 'WANShow.bingo is mentioned', 'Events', NOW()),
       ('TBmEkuWwpR', 'Not my problem anymore, talk to the new ceo.', 'Events', NOW()),
       ('r9Mnf_ME0O', 'Linus: it''s not my fault', 'Linus', NOW()),
       ('8pce6MxeUu', 'Sponsored by dbrand!', 'Sponsors', NOW()),
       ('BlxEmS0DoW', 'Sponsored by SquareSpace!', 'Sponsors', NOW()),
       ('73KLQbGb8q', 'Linus eats/drinks something', 'Linus', NOW()),
       ('sljhxujZPi', 'Linus complains about YouTube', 'Linus', NOW()),
       ('trAzkzKWLx', 'Live Tech Product Unboxing', 'Events', NOW()),
       ('mfTNmsg6UI', 'Talking over Audio', 'Set/Production', NOW()),
       ('WeV8Fp5qTk', 'Trust me Bro', 'Events', NOW()),
       ('Qp6vz_-EAD', 'Linus is out of touch', 'Linus', NOW()),
       ('-Q5HcciRTE', 'Too long (15 min+) on a topic', 'Events', NOW()),
       ('jjhqqMEYt7', 'Too long (5 mins+) on a merch message', 'Events', NOW()),
       ('kBkO1ou5Mg', 'Innuendo', 'Events', NOW()),
       ('jEU-ZVrpBR', '*oh yeah I guess we should do sponsor spots*', 'Events', NOW()),
       ('BV1VTXBrjr', 'Linus rants while Luke interacts with chat bored', 'Linus', NOW()),
       ('DEVW39oIdq', 'Linus changes camera angle for dan ''I''ve got it''', 'Linus', NOW()),
       ('SSd2LAc9FQ', 'Linus: ''Since you put me on the spot''', 'Linus', NOW()),
       ('qMU9ISD_9Q', 'New merch launch', 'Events', NOW()),
       ('HhLMpqN2FA', 'Putting off Merch Messages ''We''ll get them in a sec''', 'Events', NOW()),
       ('bEFxAZvdfT', 'DAN tries to talk but is muted', 'Dan', NOW()),
       ('4ArJxjLxtv', 'DING', 'Events', NOW()),
       ('Yyb48SI0jU', 'Dennis overboard sponsorspot', 'Sponsors', NOW()),
       ('ukjQJ5_-CK', 'Going back to previous topic/merch message', 'Events', NOW()),
       ('hQ3vMD9xNl', '*The Hackening*', 'Events', NOW()),
       ('KeThW3a1tm', 'Luke struggles to pick a topic', 'Luke', NOW()),
       ('23Cv4gW6DY', '*Spicy Take*', 'Events', NOW()),
       ('X_335PYYBn', 'Dan cam with no Dan', 'Dan', NOW()),
       ('KVeHSSKAWk', 'Dan goes and gets snacks for Linus', 'Dan', NOW()),
       ('Avkh2dEj_R', 'Linus ''turns off'' Dan', 'Dan', NOW()),
       ('VKWvCKvn0s', '*Where was I going with this?*', 'Linus', NOW()),
       ('A4EhjSc-Fe', 'Luke doing a concern', 'Luke', NOW()),
       ('qvQmuF6SnJ', 'Linus rants for over 2 mins without input from Luke', 'Linus', NOW()),
       ('J6KHwBchwS', 'Dan goes AFK', 'Dan', NOW()),
       ('W97lSUnnJb', 'Mentions another creator', 'Events', NOW()),
       ('jI3iC_XkFz', 'Linus says: ''Look, the thing is''', 'Linus', NOW()),
       ('NW4TF_QlJS', 'Google News', 'Topics', NOW()),
       ('CEk26e6DFO', 'Videogame Topic', 'Topics', NOW()),
       ('kU6Q-jECe2', 'We''re Hiring', 'Events', NOW()),
       ('ALK9YM8Hm2', 'Linus calls someone live on show', 'Linus', NOW()),
       ('wcdAehN0u-', 'Linus explains how to human', 'Linus', NOW()),
       ('oZuNEsBKSV', 'Linus gets signed out of something', 'Linus', NOW()),
       ('wYu9kCUSEG', 'Linus mocks French people', 'Linus', NOW()),
       ('sxxPZxny2w', 'Luke talks about movies', 'Luke', NOW()),
       ('VX3l5edoJa', 'Dan mocks a file extension', 'Dan', NOW()),
       ('cDJ-tz58hc', 'Linus, Luke, or Dan make a terrible joke/pun', 'Events', NOW()),
       ('0_R_ptyCeC', 'Not Financial Advice', 'Events', NOW()),
       ('HbDcO3mpTN', 'Floatplane / Labs Preview', 'Topics', NOW()),
       ('etmSz-0MgR', 'Rapid-fire merch messages', 'Events', NOW()),
       ('AL0_lZxh-Q', 'Someone in the audience buys a domain', 'Events', NOW()),
       ('Xkq3e-6Efn', 'Linus Theft/Legal Tips', 'Events', NOW()),
       ('juZ0aAds-9', 'Linus talks about an upcoming product/video', 'Linus', NOW()),
       ('Ws6Y4nd9KB', 'Special Guest on Stream (Anyone but DLL)', 'Set/Production', NOW());