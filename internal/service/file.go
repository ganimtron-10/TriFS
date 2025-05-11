package service

type Message struct {
	Code string
}

type ReadFileRequest struct {
	*Message
	Filename string
}

type ReadFileResponse struct {
	*Message
	Filename string
	Data     []byte
}

type FileService struct{}

func (service *FileService) Read(request *ReadFileRequest, response *ReadFileResponse) error {
	response.Filename = request.Filename
	response.Data = []byte{0, 1, 2, 3, 4, 5}
	return nil
}
