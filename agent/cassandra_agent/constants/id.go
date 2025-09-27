package constants

type DataID int 

const (
	ConnTypeNodeTool DataID = iota
	ConnTypeCQLTool
	ConnTypeAgent
)

const (
	NodeToolTpStatsData DataID = iota
	NodeToolInfoData
	CQLTracesSessions
	CQLRunningQuerys
	AgentHostCpu
	AgentHostMemory
)