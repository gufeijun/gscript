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

var StdLibMap = map[string]uint32{
	"Buffer": 0,
	"fs":     1,
	"os":     2,
}

var stdLibs = []string{"Buffer", "fs", "os"}

func GetLibNameByProtoNum(num uint32) string {
	return stdLibs[num]
}

func GetLibProtoNumByName(name string) (uint32, error) {
	num, ok := StdLibMap[name]
	if !ok {
		return 0, fmt.Errorf("invalid std libarary: %s", name)
	}
	return num, nil
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
	protos := make([]proto.Proto, len(StdLibMap))
	for lib, num := range StdLibMap {
		p, err := ReadProto(lib)
		if err != nil {
			return nil, err
		}
		protos[num] = p
	}
	return protos, nil
}
