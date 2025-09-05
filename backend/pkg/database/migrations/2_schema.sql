-- +goose Up
CREATE TABLE IF NOT EXISTS albums (
    id VARCHAR(36) NOT NULL,
    title TEXT NOT NULL,
    description TEXT,
    creator VARCHAR(255) NOT NULL,
    created_at DATETIME NOT NULL,
    updated_at DATETIME NOT NULL,
    PRIMARY KEY (id)
);
CREATE TABLE IF NOT EXISTS images (
    id VARCHAR(36) NOT NULL,
    PRIMARY KEY (id)
);
CREATE TABLE IF NOT EXISTS album_images (
    id VARCHAR(36) NOT NULL,
    album_id VARCHAR(36) NOT NULL,
    image_id VARCHAR(36) NOT NULL,
    PRIMARY KEY (id),
    FOREIGN KEY (album_id) REFERENCES albums(id),
    FOREIGN KEY (image_id) REFERENCES images(id)
);