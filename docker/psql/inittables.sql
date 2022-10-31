DROP TABLE IF EXISTS versions;

DROP TABLE IF EXISTS notes;

CREATE TABLE notes (
    id SERIAL NOT NULL,
    title TEXT,
    c_date NUMERIC,
    PRIMARY KEY(id)
);

CREATE TABLE versions (
    id SERIAL NOT NULL,
    text TEXT,
    c_date NUMERIC,
    ch_sum TEXT,
    note_id INTEGER,
    PRIMARY KEY(id),
    CONSTRAINT fk_note
        FOREIGN KEY(note_id)
            REFERENCES notes(id)
);