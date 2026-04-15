package vo

type Volume struct {
	// Seq must be a positive number.
	Seq   int
	Name  string
	Parts []Part
}
