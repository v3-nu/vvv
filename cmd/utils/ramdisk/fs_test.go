// The only tests that really make sense are integration tests. Typically I
// would hide these behind a `+integration` build flag, so that you have to
// specify you want them for go test, but since they are the only ones might as
// well keep things simple.

package ramdisk_test

import (
	"bytes"
	"context"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"syscall"
	"testing"

	"github.com/v3-nu/vvv/cmd/utils/ramdisk"

	"github.com/hanwen/go-fuse/v2/fs"
	"github.com/hanwen/go-fuse/v2/fuse"
)

type HelloRoot struct {
	fs.Inode
}

func (r *HelloRoot) OnAdd(ctx context.Context) {
	ch := r.NewPersistentInode(
		ctx, &fs.MemRegularFile{
			Data: []byte("HELLO WORLD"),
			Attr: fuse.Attr{
				Mode: 0644,
			},
		}, fs.StableAttr{Ino: 2})
	r.AddChild("file.txt", ch, false)
}

func (r *HelloRoot) Getattr(ctx context.Context, fh fs.FileHandle, out *fuse.AttrOut) syscall.Errno {
	out.Mode = 0755
	return 0
}

var _ = (fs.NodeGetattrer)((*HelloRoot)(nil))
var _ = (fs.NodeOnAdder)((*HelloRoot)(nil))

func TestBasic(t *testing.T) {
	opts := &fs.Options{}
	opts.Debug = true
	opts.MountOptions = fuse.MountOptions{
		FsName: "test",
	}

	hr := &HelloRoot{}

	server, err := fs.Mount("/tmp/goramdisk-test", hr, opts)
	if err != nil {
		t.Fatalf("failed to mount: %v", err)
	}
	defer server.Unmount()
	server.Wait()

	// hr.OnAdd(context.Background())

	// os.Mkdir("/tmp/goramdisk-test", 0755)

	// rd, err := ramdisk.Create(ramdisk.Options{
	// 	Size:      1 * ramdisk.MB,
	// 	Logger:    testLogger(t),
	// 	MountPath: "/tmp/goramdisk-test",
	// })
	// ifErrLogFatalDetailed(err, t)
	// defer ramdisk.Destroy(rd.DevicePath)

	err = os.WriteFile("/tmp/goramdisk-test/file.txt", []byte("hello world"), 0644)
	if err != nil {
		t.Fatalf("failed to write file: %v", err)
	}

	// err = os.WriteFile(file, []byte("hello world"), 0644)
	// if err != nil {
	// 	t.Fatalf("failed to write file: %v", err)
	// }

	// subp, err := os.StartProcess("cat", []string{"/tmp/goramdisk-test/file.txt"}, &os.ProcAttr{
	// 	Files: []*os.File{nil, os.Stdout, os.Stderr},
	// })
	// defer subp.Release()
	// // subp.Wait()
	// if err != nil {
	// 	t.Fatalf("failed to start process: %v", err)
	// }
	// err = ramdisk.Destroy(rd.DevicePath)
	// if err != nil {
	// 	t.Fatalf("failed to destroy ramdisk: %v", err)
	// }
}

/*
	basic test sequence I'm thinking about

- try to create, no errors
- use go FS to see if path exists
- if possible, stat the directory and try to read capacity?
- use go FS to write a file to ramdisk
- use go FS to read a file back from ramdisk, verify contents
- unmount ramdisk, no errors
- go FS to see if file exists (should not)
- go FS to verify mountpath no longer exists (if later implemented)
*/
func TestCreate(t *testing.T) {
	// try to create, verify no errors
	rd, err := ramdisk.Create(ramdisk.Options{
		Logger: testLogger(t),
	})
	ifErrLogFatalDetailed(err, t)

	// mark for cleanup after tests complete
	defer ramdisk.Destroy(rd.DevicePath)

	// use go FS to see if mountpath exists and is directory
	t.Log("verifying mountpath", rd.MountPath)
	stat, err := os.Stat(rd.MountPath)
	if err != nil {
		t.Fatal(err)
	}
	if !stat.IsDir() {
		t.Fatalf("%v is not a directory", rd.MountPath)
	}

	// use go FS to write a file to ramdisk
	target := filepath.Join(rd.MountPath, "msg.txt")
	payload := []byte("hi mom")
	t.Log("writing file to", target)
	err = ioutil.WriteFile(target, payload, 0644)
	if err != nil {
		t.Fatal("failed to write file to", target)
	}

	// use go FS to read a file back from ramdisk, verify contents
	t.Logf("reading back file %s to verify contents", target)
	data, err := ioutil.ReadFile(target)
	if err != nil {
		t.Fatal(err)
	}
	if string(data) != string(payload) {
		t.Errorf("read data did not match, want [%s] got [%s]", payload, data)
	}

	// unmount ramdisk, no errors
	t.Log("unmounting", rd.DevicePath)
	err = ramdisk.Destroy(rd.MountPath)
	ifErrLogFatalDetailed(err, t)

	// go FS to see if file exists (should not)
	t.Logf("checking for existence of %v (should not exist)", target)
	if _, err := os.Stat(target); os.IsExist(err) {
		t.Fatalf("%v exists when it should not!", target)
	}
}

func TestDestroy(t *testing.T) {
	rd, err := ramdisk.Create(ramdisk.Options{})
	if err != nil {
		t.Skip("precondition creating ramdisk failed!")
	}
	err = ramdisk.Destroy(rd.DevicePath)
	ifErrLogFatalDetailed(err, t)
}

// helper function to wrap some boilerplate for returns from Cmd.Exec,
// if there is an error log it fatally, but checking if its an ExitError
// and using the more detailed Stderr value if so.
func ifErrLogFatalDetailed(err error, t *testing.T) {
	if err != nil {
		exiterr, ok := err.(*exec.ExitError)
		if ok {
			t.Fatalf("ExitError: %s\n", exiterr.Stderr)
		} else {
			t.Fatal(err)
		}
	}
}

// boilerplate trick to be able to wrap t.Log() as a Logger
// via: https://github.com/golang/go/issues/22513
func testLogger(t *testing.T) *log.Logger {
	return log.New(testWriter{t}, "testProcessLog: ", log.LstdFlags)
}

type testWriter struct {
	t *testing.T
}

func (tw testWriter) Write(p []byte) (n int, err error) {
	tw.t.Log(string(bytes.TrimSpace(p)))
	return len(p), nil
}
