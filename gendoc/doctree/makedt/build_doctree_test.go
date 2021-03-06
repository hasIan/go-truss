package makedt

import (
	"reflect"
	"testing"

	// This has to be imported because it modifies the state of `proto` by
	// registering the `google.api.http` extension, allowing us to specify it
	// in the sources below.
	_ "github.com/TuneLab/go-truss/gendoc/doctree/makedt/googlethirdparty"
	//"github.com/davecgh/go-spew/spew"

	"github.com/TuneLab/go-truss/gendoc/doctree"
	"github.com/golang/protobuf/proto"
	descriptor "github.com/golang/protobuf/protoc-gen-go/descriptor"
)

func TestNewFile(t *testing.T) {
	src := `
		name: "path/to/example.proto",
		package: "example"
		message_type <
			name: "StringMessage"
			field <
				name: "string"
				number: 1
				label: LABEL_OPTIONAL
				type: TYPE_STRING
			>
		>
		service <
			name: "ExampleService"
			method <
				name: "Echo"
				input_type: "StringMessage"
				output_type: "StringMessage"
				options <
					[google.api.http] <
						post: "/v1/example/echo"
						body: "*"
					>
				>
			>
		>
	`
	var fd descriptor.FileDescriptorProto
	if err := proto.UnmarshalText(src, &fd); err != nil {
		t.Fatalf("proto.UnmarshalText(%s, &fd) failed with %v; want success", src, err)
	}

	dt := doctree.MicroserviceDefinition{}
	newFile, err := NewFile(&fd, &dt)
	if err != nil {
		t.Fatalf("Error creating new file: %v\n", err)
	}

	msg := doctree.ProtoMessage{
		Fields: []*doctree.MessageField{
			&doctree.MessageField{
				Label:  "LABEL_OPTIONAL",
				Number: 1,
			},
		},
	}
	msg.SetName("StringMessage")
	msg.Fields[0].SetName("string")
	msg.Fields[0].Type.SetName("TYPE_STRING")

	f := &doctree.ProtoFile{
		Messages: []*doctree.ProtoMessage{
			&msg,
		},
		Services: []*doctree.ProtoService{
			&doctree.ProtoService{
				FullyQualifiedName: ".example.ExampleService",
				Methods: []*doctree.ServiceMethod{
					&doctree.ServiceMethod{
						RequestType:  &msg,
						ResponseType: &msg,
					},
				},
			},
		},
	}
	f.SetName("path/to/example.proto")
	f.Services[0].SetName("ExampleService")
	f.Services[0].Methods[0].SetName("Echo")

	if got, want := newFile, f; !reflect.DeepEqual(got, want) {
		t.Errorf("doctree.ProtoFile = \n%#v, want = \n%#v\n", got, want)
	}
}
