```go

// Type Domain -- /core/domain/attribute.go
type Attribute struct {
	common.Persistent
}

// Use-Cases -- /core/operators/port.go
type AttributeService interface {}

type AttributeOperator interface {}

type AttributeRepository interface {}

// Service -- /core/services/attribute.go
type attributeService struct {}

func NewAttributeService(ports.AttributeRepository, ports.AttributeOperator) ports.AttributeService

// Operator -- /core/operators/attribute.go
type attributeOperator struct {}

func NewAttributeOperator(*controller.Controller) ports.AttributeOperator

// Repository -- /core/repository/attribute.go
type attributeRepository struct {}

func NewAttributeRepository(*gorm.DB) ports.AttributeRepository

// Routes -- /port/routes/attribute.go
type attributeRouter struct {}

func NewAttributeRouter(ports.AttributeService) ports.AttributeRepository

```