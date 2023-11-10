package meta

import (
	"strings"

	"google.golang.org/grpc/metadata"
)

// MetadataTextMap implements opentracing TextMap Reader, Writer.
type MetadataTextMap struct {
	metadata.MD
}

func (m MetadataTextMap) Set(k, v string) {
	k = strings.ToLower(k)
	m.MD.Append(k, v)
}

func (m MetadataTextMap) ForeachKey(handler func(k, v string) error) error {
	for k, vs := range m.MD {
		for _, v := range vs {
			if err := handler(k, v); err != nil {
				return err
			}
		}
	}

	return nil
}
