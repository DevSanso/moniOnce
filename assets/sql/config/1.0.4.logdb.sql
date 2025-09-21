create table logdb (
	logdb_id bigint,
	logdb_type varchar(6),
	host varchar(256),
	port int,
	user varchar(256),
	password varchar(256),
	dbname varchar(256),
	
	"version" numeric
);