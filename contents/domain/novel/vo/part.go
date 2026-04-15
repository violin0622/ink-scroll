package vo

type Part struct {
	// Seq must be a positive number.
	Seq int

	// Name is optional
	Name string

	Chapters []Chapter
}
