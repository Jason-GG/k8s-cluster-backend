package enums

type LoadBalancerStatus string

const (
	StatusActive   LoadBalancerStatus = "active"
	StatusLocked   LoadBalancerStatus = "locked"
	StatusInActive LoadBalancerStatus = "inactive"
)

type AssociateType string

const (
	AssociateTypeDefaultBackend AssociateType = "default_backend"
	AssociateTypeVServerGroup   AssociateType = "virtual_server_group"
)

type AssociateID string
