CREATE TYPE activity_type AS (activity_type text, body text, target_name text, target_url text, created_at timestamp);
CREATE TABLE activities OF activity_type;

CREATE TYPE measurement_type AS (sensor_type smallint, value real, created_at timestamp);
CREATE TABLE measurements OF measurement_type;
