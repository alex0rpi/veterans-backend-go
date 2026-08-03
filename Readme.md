# TABLES

## Storing Images

For the time being we'll only have ONE SINGLE TABLE for storing all images metadata.

That is the first purpose of this backend, to allow retrieval of all images of the front end of the site.

This table would have the following columns:

IMAGE_METADATA

- id
- storage_key
- original_filename
- mime_type
- created_at
- height
- width
- filesize
- blur_key
- small_key
- medium_key
- large_key

## Storing administratos

Additionally we'll have an admin table:

ADMIN_USERS

- id
- email
- password
- token_hash
- expires_at
- created_at
