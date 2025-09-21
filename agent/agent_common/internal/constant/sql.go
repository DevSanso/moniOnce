package constant

/**
select query
*/
const (
	SELECT_OBJECT_CONFIG_QUERY = "SELECT category || '_' || \"key\", \"value\" FROM objects_config WHERE object_id = $1 "
)
