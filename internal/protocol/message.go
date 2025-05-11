package protocol

type ReadFileArgs struct {
	Filename string
}

type ReadFileReply struct {
	Filename string
	Data     []byte
}
