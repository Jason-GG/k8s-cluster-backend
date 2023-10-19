package enums

type PortType int

const (
	PortApplication PortType = 200
	PortDebug                = 300
	PortSchedule             = 400
)

const (
	_ = iota
	DEV
	TEST
	PRE
	PROD
	ST
)

const (
	TypeReserved  = "reserve"
	TypeDroppable = "droppable"
)

func EnvString(enum int) string {
	switch enum {
	case DEV:
		return "dev"
	case TEST:
		return "test"
	case PRE:
		return "pre"
	case PROD:
		return "prod"
	case ST:
		return "st"
	}
	return ""
}
