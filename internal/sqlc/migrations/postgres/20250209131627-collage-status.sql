-- +migrate Up
CREATE TYPE collage_status AS ENUM (
    'pending',
    'processing',
    'ready',
    'failed'
);

ALTER TABLE collages
ADD COLUMN status collage_status NOT NULL DEFAULT 'pending';

-- +migrate Down
ALTER TABLE collages
DROP COLUMN status;

DROP TYPE collage_status;