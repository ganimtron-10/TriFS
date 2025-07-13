package protocol

type ReadFileRequestArgs struct {
	Filename string
}

type ReadFileRequestReply struct {
	WorkerUrls []string
}

type ReadFileArgs struct {
	Filename string
}

type ReadFileReply struct {
	Filename string
	Data     []byte
}

type HeartBeatArgs struct {
	Address    string
	FileHashes map[string]struct{}
}

type HeartBeatReply struct {
}

type WriteFileRequestArgs struct {
	Filename string
}

type WriteFileRequestReply struct {
	WorkerUrls []string
}

type WriteFileArgs struct {
	Filename string
	Data     []byte
}

type WriteFileReply struct {
}
