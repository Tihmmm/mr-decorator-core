package errors

type DownloadError struct {
	Err string
}

func (e *DownloadError) Error() string {
	return "can't download artifact: " + e.Err
}
