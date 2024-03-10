package xray

type FileLocation int8

const (
	Remote FileLocation = iota + 1
	Local
	Both
)

type OperationType int8

const (
	Create OperationType = iota + 1
	Update
	Drop
)
