package image

type Reader interface {
	Read(filePath string) (*Image, error)
}
