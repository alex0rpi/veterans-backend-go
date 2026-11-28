# Prospect actions

## Future endpoints

### People endpoints

- Endpoint to upload a new board member
  - name: `uploadBoardMember`
  - incorporaría formData con una o + imagenes y los datos de la persona.
  - Crea registros en la tabla de imágenes y en la futura tabla de boardMembers.

- Endpoint to classify an existing board member into a historical figure
  - name: `migrateBoardMember`
  - Dado un identificador de miembro y una fecha de finación, trasladar sus datos a la tabla de historical members y borrar su registro en boardMembers

- Endpoint to upload a new historical figure
  - name: `uploadHistoricalMemebr`
  - incorporaría formData con una o + imagenes y los datos de la persona.
  - Crea registros en la tabla de imágenes y en la futura tabla de historicalMembers.

- Endpoints to:
  - Add an additional foto to an existing member: `addNewFotoToMember`
  - Remove an image from a existing memebr: `deleteMemberFoto`

> Cada tabla, boardMembers, y historicalMembers, debe incorporar un campo/columna que indique qué imagen es la principal mostrar en primer lugar.

### Season endpoints

- Endpoint para crear una nueva temporada
  - `UploadNewSeason`
  - Tendrá campos algunos obligatorios otros opcionales:
  - Secció professional (títol, text, imatges)
  - Secció bases i filial (títol, text, imatges)
  - Secció veterans (títol, text, imatges)

- Endpoint per a actualitzar una temporada
  - Per exemple per enriquirla amb noves informacions, imatges, correccions de text/títol, o bé marcant una opció "visible/no visible"
  - `UpdateSeason`

- Endpoint per a esborrar una temporada. `deleteSeason` <-- OJU

### News/articles endpoints

PUBLIC

- Endpoint to show all existing: `getAllPosts`
- Endpoint to get a single post: `getSinglePost`

PROTECTED

- Endpoint to upload a new post: `createPost`
  - Potser un formInput amb 1 o+ imatges, un títol i un text.
- Endpoint to delete a post: `deletePost`
- Endpoint to update an existing post: `updatePost`
  - Potser un formInput amb 1 o+ imatges, un títol i un text.
- Endpoint to list existing posts: `getPostList`

> Potser caldrà afegir un nou imageContext: mediaPost

## Prospect migrations to consider

### Add board members table

```sql
-- +goose Up
CREATE TABLE board_members (
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
ADD COLUMN board_member_id INTEGER REFERENCES board_members(id) ON DELETE
SET NULL;
-- +goose Down
ALTER TABLE image_metadata DROP COLUMN IF EXISTS board_member_id;
DROP TABLE IF EXISTS board_members;
```

### Add historical members table

```sql
-- +goose Up
CREATE TABLE historical_members (
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
ADD COLUMN historical_member_id INTEGER REFERENCES historical_members(id) ON DELETE
SET NULL;
-- +goose Down
ALTER TABLE image_metadata DROP COLUMN IF EXISTS historical_member_id;
DROP TABLE IF EXISTS historical_members;
```

### Add season data table

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
    author VARCHAR(255) NOT NULL,
    visible BOOLEAN DEFAULT TRUE NOT NULL
);
-- +goose Down
DROP TABLE IF EXISTS season_data;
```

### Add socialPost table

```sql
-- +goose Up
CREATE TABLE social_posts (
    id SERIAL PRIMARY KEY,
    author VARCHAR(255) NOT NULL,
    title TEXT NOT NULL,
    content TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ,
    visible BOOLEAN DEFAULT TRUE NOT NULL
);
ALTER TABLE image_metadata
ADD COLUMN social_post_id INTEGER REFERENCES social_posts(id)
ON DELETE CASCADE;

-- REVISITAR ESTO por si hay que crear más índices
-- Cada post tiene su propia secuencia de posiciones
CREATE UNIQUE INDEX unique_position_per_post ON image_metadata (
    social_post_id,
    display_position
) WHERE media_context = 'mediaPost' AND social_post_id IS NOT NULL;

-- +goose Down
ALTER TABLE image_metadata DROP COLUMN IF EXISTS social_post_id;
DROP TABLE IF EXISTS social_posts;
```
