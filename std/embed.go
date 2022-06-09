package std

import (
	"bufio"
	"embed"
	"fmt"
	"gscript/proto"
)

//go:embed *.gsproto
var ProtoFiles embed.FS

const ProtoSuffix string = ".gsproto"

var StdLibs = map[string]uint32{
	"Buffer": 0,
	"fs":     1,
	"os":     2,
}

func ReadProto(lib string) (proto.Proto, error) {
	file, err := ProtoFiles.Open(lib + ProtoSuffix)
	if err != nil {
		return proto.Proto{}, fmt.Errorf("can not find std libarary bytes code: %v", err)
	}
	defer file.Close()
	return proto.ReadProto(bufio.NewReader(file))
}

func ReadProtos() ([]proto.Proto, error) {
	protos := make([]proto.Proto, len(StdLibs))
	for lib, num := range StdLibs {
		p, err := ReadProto(lib)
		if err != nil {
			return nil, err
		}
		protos[num] = p
	}
	return protos, nil
}
