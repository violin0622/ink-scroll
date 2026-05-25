package vo

import (
	"fmt"

	"github.com/oklog/ulid/v2"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/encoding/protojson"
)

func Example() {
	e := volumeHasChapterError(ulid.Make(), ulid.Make(), 3)
	fmt.Println(e)

	bytes, _ := protojson.MarshalOptions{
		Multiline: true,
		Indent:    `  `,
	}.Marshal(status.Convert(e).Proto())
	fmt.Println(string(bytes))
	// Output: ok
}
