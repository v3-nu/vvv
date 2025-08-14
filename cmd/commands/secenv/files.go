package secenv

import (
	"github.com/hanwen/go-fuse/v2/fs"
)

type HelloRoot struct {
	fs.Inode
}

// func (r *HelloRoot) OnAdd(ctx context.Context) {
// 	// ch := r.NewPersistentInode(
// 	// 	ctx, &fs.MemRegularFile{
// 	// 		Data: []byte("HELLO WORLD"),
// 	// 		Attr: fuse.Attr{
// 	// 			Mode: 0644,
// 	// 		},
// 	// 	}, fs.StableAttr{Ino: 2})
// 	// r.AddChild("file.txt", ch, false)
// }

// func (r *HelloRoot) Getattr(ctx context.Context, fh fs.FileHandle, out *fuse.AttrOut) syscall.Errno {
// 	// out.Mode = 0755
// 	return 0
// }

// var _ = (fs.NodeGetattrer)((*HelloRoot)(nil))
// var _ = (fs.NodeOnAdder)((*HelloRoot)(nil))

// func TestBasic(t *testing.T) {
// 	opts := &fs.Options{}
// 	opts.Debug = true
// 	opts.MountOptions = fuse.MountOptions{
// 		FsName: "test",
// 	}

// 	hr := &HelloRoot{}

// 	server, err := fs.Mount("/tmp/goramdisk-test", hr, opts)
// 	if err != nil {
// 		t.Fatalf("failed to mount: %v", err)
// 	}
// 	defer server.Unmount()
// 	server.Wait()

// 	// hr.OnAdd(context.Background())

// 	// os.Mkdir("/tmp/goramdisk-test", 0755)

// 	// rd, err := ramdisk.Create(ramdisk.Options{
// 	// 	Size:      1 * ramdisk.MB,
// 	// 	Logger:    testLogger(t),
// 	// 	MountPath: "/tmp/goramdisk-test",
// 	// })
// 	// ifErrLogFatalDetailed(err, t)
// 	// defer ramdisk.Destroy(rd.DevicePath)

// 	err = os.WriteFile("/tmp/goramdisk-test/file.txt", []byte("hello world"), 0644)
// 	if err != nil {
// 		t.Fatalf("failed to write file: %v", err)
// 	}

// 	// err = os.WriteFile(file, []byte("hello world"), 0644)
// 	// if err != nil {
// 	// 	t.Fatalf("failed to write file: %v", err)
// 	// }

// 	// subp, err := os.StartProcess("cat", []string{"/tmp/goramdisk-test/file.txt"}, &os.ProcAttr{
// 	// 	Files: []*os.File{nil, os.Stdout, os.Stderr},
// 	// })
// 	// defer subp.Release()
// 	// // subp.Wait()
// 	// if err != nil {
// 	// 	t.Fatalf("failed to start process: %v", err)
// 	// }
// 	// err = ramdisk.Destroy(rd.DevicePath)
// 	// if err != nil {
// 	// 	t.Fatalf("failed to destroy ramdisk: %v", err)
// 	// }
// }
