package protocol

type ReadFileArgs struct {
	Filename string
}

type ReadFileReply struct {
	Filename string
	Data     []byte
}

type HeartBeatArgs struct {
	Address string
}

type HeartBeatReply struct {
}

type WriteFileRequestArgs struct {
	Filename string
}

type WriteFileRequestReply struct {
	WorkerUrl string
}

type WriteFileArgs struct {
	Filename string
	Data     []byte
}

type WriteFileReply struct {
	WorkerUrl string
}
