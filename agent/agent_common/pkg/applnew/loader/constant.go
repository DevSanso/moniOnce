package loader

/**
select query
*/
const (
	_SELECT_OBJECT_CONFIG_QUERY = " SELECT \"key\", \"value\" FROM objects_config WHERE object_id = $1 and category = 'config' "
	_SELECT_OBJECT_SYNC_QUERY = " SELECT \"key\", \"value\" FROM objects_config WHERE object_id = $1 and category = 'sync' "
	_SELECT_OBJECT_FLAG_QUERY = " SELECT \"key\", \"value\" FROM objects_config WHERE object_id = $1 and category = 'flag' "
	_UPDATE_OBJECT_FLAG_QUERY = " UPDATE objects_config SET \"value\" = $1 WHERE object_id = $1 "
)
