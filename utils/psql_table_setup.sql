CREATE TABLE activities 
(
	id				SERIAL PRIMARY KEY, 
	activity_type	text, 
	body			text, 
	target_name		text, 
	target_url		text, 
	created_at		timestamp
);

CREATE TABLE measurements
(
	id				SERIAL PRIMARY KEY, 
	sensor_type		smallint, 
	value			real, 
	created_at		timestamp
);
