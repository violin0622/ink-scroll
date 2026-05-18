package vo

type Volume struct {
	ID ID
	// Seq must be a positive number.
	Seq         int
	Title       string
	Description string
}
