package akizuki

type NewPageDetector interface {
	NewPages(urls []string) ([]string, error)
	AddPages(urls []string) error
}
