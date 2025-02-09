
-- +migrate Up
CREATE TABLE collage_sections (
    db_id SERIAL PRIMARY KEY,
    id UUID UNIQUE NOT NULL,
    image_id UUID NOT NULL,
    collage_id UUID NOT NULL,
    section INTEGER NOT NULL,
    created_at DATE NOT NULL,
    updated_at DATE NOT NULL,
    FOREIGN KEY (collage_id) REFERENCES collages(id)
);

-- +migrate Down
DROP TABLE collage_sections;
