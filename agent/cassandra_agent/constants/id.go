package constants

type DataID int 

const (
	ConnTypeNodeTool DataID = iota
	ConnTypeCQLTool
)

const (
	NodeToolTpStatsData DataID = iota
	NodeToolInfoData
)