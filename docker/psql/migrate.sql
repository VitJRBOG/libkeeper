DROP TABLE IF EXISTS version;

DROP TABLE IF EXISTS note;

CREATE TABLE note (
    id SERIAL NOT NULL UNIQUE,
    c_date TEXT,
    PRIMARY KEY(id)
);

CREATE TABLE version (
    id SERIAL NOT NULL UNIQUE,
    full_text TEXT,
    c_date TEXT,
    checksum TEXT,
    note_id INTEGER,
    PRIMARY KEY(id),
    CONSTRAINT fk_note FOREIGN KEY(note_id) REFERENCES note(id)
);