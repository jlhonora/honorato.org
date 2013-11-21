/* Drop a type and its references: */
DROP TYPE activity_type CASCADE;

/* Apply a file: */
psql -U pgmainuser --file=utils/psql_table_setup.sql pgmaindb

