CREATE TABLE images (
	id TEXT PRIMARY KEY,
	created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
	updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
	image_created_at DATETIME NOT NULL
);

CREATE TABLE tags (
	id INTEGER PRIMARY KEY,
	name TEXT NOT NULL,
	image_id TEXT NOT NULL,
	created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,

	FOREIGN KEY(image_id) REFERENCES images(id)
);

CREATE UNIQUE INDEX tags_name_image_id ON tags (name, image_id);
CREATE UNIQUE INDEX tags_image_id_name ON tags (image_id, name);
CREATE INDEX tags_name ON tags (name);
CREATE INDEX tags_image_id ON tags (image_id);
