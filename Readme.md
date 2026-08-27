# TABLES

## Storing Images

I don't want to build a "wordpress" backend to make CRUD to every single elemnent of the page, that's not the purpose.

The web content is actually quite static once it is posted. The main purpose of this backend would be to streamline the image loading for every section and article; by storing every image in a Cloudfare S3 server or simmilar, keep each image metadata in a postgresSql DB, and then retrieving the applicable images for each section.

For the time I just set up ONE SINGLE TABLE for storing all images metadata.

That is the first purpose of this backend, to allow retrieval of all images of the front end of the site.

I Might add in the future some additional tables for managing users or perhaps a relationship table between:

- web sections
- image used by those web sections

---

- Get available display_indexes
- Update all display indexes for a given context and season
- images display position (int) MUST BE UNIQUE for their given context and season (if the context is season)
