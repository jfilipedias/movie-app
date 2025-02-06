CREATE DATABASE movie CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;

USE movie;

CREATE TABLE IF NOT EXISTS movies (
  id VARCHAR(255),
  title VARCHAR(255),
  description TEXT,
  director VARCHAR(255)
);

CREATE TABLE IF NOT EXISTS ratings (
  record_id VARCHAR(255),
  record_type VARCHAR(255),
  user_id VARCHAR(255),
  value INT
);
