# Prospect migrations to consider

## Add board members tale and historical members table

```sql
-- +goose Up
CREATE TABLE board_members_data (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    role VARCHAR(255) NOT NULL,
    role_esp VARCHAR(255) NOT NULL,
    player_seasons VARCHAR(255) NOT NULL,
    birth_date DATE,
    birth_country VARCHAR(255),
    birth_city VARCHAR(255),
    anecdote TEXT,
    anecdote_esp TEXT,
    trajectory TEXT,
    trajectory_esp TEXT,
    other_comment TEXT,
    other_comment_esp TEXT
);
ALTER TABLE image_metadata
ADD COLUMN board_member_id INTEGER REFERENCES board_members_data(id) ON DELETE
SET NULL;

CREATE TABLE historical_members_data (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    former_role VARCHAR(255) NOT NULL,
    former_role_esp VARCHAR(255) NOT NULL,
    player_seasons VARCHAR(255) NOT NULL,
    birth_date DATE,
    passing_date DATE,
    birth_country VARCHAR(255),
    birth_city VARCHAR(255),
    passing_country VARCHAR(255),
    passing_city VARCHAR(255),
    anecdote TEXT,
    anecdote_esp TEXT,
    trajectory TEXT,
    trajectory_esp TEXT,
    other_comment TEXT,
    other_comment_esp TEXT
);
ALTER TABLE image_metadata
ADD COLUMN historical_member_id INTEGER REFERENCES historical_members_data(id) ON DELETE
SET NULL;
-- +goose Down
ALTER TABLE image_metadata DROP COLUMN IF EXISTS board_member_id;
ALTER TABLE image_metadata DROP COLUMN IF EXISTS historical_member_id;
DROP TABLE IF EXISTS board_members_data;
DROP TABLE IF EXISTS historical_members_data;
```

## Add season data table

```sql
-- +goose Up
CREATE TABLE season_data (
    id SERIAL PRIMARY KEY,
    season TEXT NOT NULL UNIQUE,
    title_pro VARCHAR(255) NOT NULL,
    title_bases VARCHAR(255) NOT NULL,
    title_vets VARCHAR(255) NOT NULL,
    content_pro TEXT NOT NULL,
    content_bases TEXT NOT NULL,
    content_vets TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    author VARCHAR(255) NOT NULL
);
SET NULL;
-- +goose Down
DROP TABLE IF EXISTS season_data;
```
