DROP TABLE IF EXISTS note_category;

DROP TABLE IF EXISTS version;

DROP TABLE IF EXISTS note;

DROP TABLE IF EXISTS category;

DROP TABLE IF EXISTS icon;

CREATE TABLE icon (
    id SERIAL NOT NULL UNIQUE,
    type TEXT NOT NULL,
    path TEXT NOT NULL UNIQUE,
    PRIMARY KEY(id)
);

CREATE TABLE category (
    id SERIAL NOT NULL UNIQUE,
    name TEXT NOT NULL UNIQUE,
    immutable INTEGER NOT NULL DEFAULT 0,
    icon_id INTEGER,
    PRIMARY KEY(id)
);

CREATE TABLE note (
    id SERIAL NOT NULL UNIQUE,
    title TEXT,
    c_date TEXT NOT NULL,
    PRIMARY KEY(id)
);

CREATE TABLE version (
    id SERIAL NOT NULL UNIQUE,
    full_text TEXT,
    c_date TEXT NOT NULL,
    checksum TEXT NOT NULL,
    note_id INTEGER NOT NULL,
    PRIMARY KEY(id),
    CONSTRAINT fk_note FOREIGN KEY(note_id) REFERENCES note(id)
);

CREATE TABLE note_category (
    id SERIAL NOT NULL UNIQUE,
    note_id INTEGER NOT NULL,
    category_id INTEGER NOT NULL,
    PRIMARY KEY(id),
    CONSTRAINT fk_nc_note FOREIGN KEY(note_id) REFERENCES note(id) ON DELETE CASCADE,
    CONSTRAINT fk_nc_category FOREIGN KEY(category_id) REFERENCES category(id) ON DELETE CASCADE
);

INSERT INTO icon(type, path) VALUES
    ('category', 'noicon.png'),
    ('category', 'all.png'),
    ('category', 'uncategorised.png'),
    ('category', 'trashed.png'),
    ('category', 'todo.png');

INSERT INTO category(name, immutable, icon_id) VALUES
    ('Uncategorised', 1, (SELECT id FROM icon WHERE path LIKE 'uncategorised%')),
    ('Trashed', 1, (SELECT id FROM icon WHERE path LIKE 'trashed%'));