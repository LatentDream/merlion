CREATE TABLE notes (
    note_id	 UUID PRIMARY KEY,
    title	 TEXT NOT NULL,
    content	 TEXT,
    tags	 TEXT[],
    is_favorite	 BOOLEAN NOT NULL DEFAULT false,
    is_work_log  BOOLEAN NOT NULL DEFAULT false,
    is_trash     BOOLEAN NOT NULL DEFAULT false,
    created_at   TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at   TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);
