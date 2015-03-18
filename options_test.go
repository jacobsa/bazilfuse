package bazilfuse_test

import (
	"os"
	"runtime"
	"testing"

	"github.com/jacobsa/bazilfuse/fs/fstestutil"
)

func init() {
	fstestutil.DebugByDefault()
}

func TestMountOptionFSName(t *testing.T) {
	if runtime.GOOS == "freebsd" {
		t.Skip("FreeBSD does not support FSName")
	}
	t.Parallel()
	const name = "FuseTestMarker"
	mnt, err := fstestutil.MountedT(t, fstestutil.SimpleFS{fstestutil.Dir{}},
		fuse.FSName(name),
	)
	if err != nil {
		t.Fatal(err)
	}
	defer mnt.Close()

	info, err := fstestutil.GetMountInfo(mnt.Dir)
	if err != nil {
		t.Fatal(err)
	}
	if g, e := info.FSName, name; g != e {
		t.Errorf("wrong FSName: %q != %q", g, e)
	}
}

func testMountOptionFSNameEvil(t *testing.T, evil string) {
	if runtime.GOOS == "freebsd" {
		t.Skip("FreeBSD does not support FSName")
	}
	t.Parallel()
	var name = "FuseTest" + evil + "Marker"
	mnt, err := fstestutil.MountedT(t, fstestutil.SimpleFS{fstestutil.Dir{}},
		fuse.FSName(name),
	)
	if err != nil {
		t.Fatal(err)
	}
	defer mnt.Close()

	info, err := fstestutil.GetMountInfo(mnt.Dir)
	if err != nil {
		t.Fatal(err)
	}
	if g, e := info.FSName, name; g != e {
		t.Errorf("wrong FSName: %q != %q", g, e)
	}
}

func TestMountOptionFSNameEvilComma(t *testing.T) {
	if runtime.GOOS == "darwin" {
		// see TestMountOptionCommaError for a test that enforces we
		// at least give a nice error, instead of corrupting the mount
		// options
		t.Skip("TODO: OS X gets this wrong, commas in mount options cannot be escaped at all")
	}
	testMountOptionFSNameEvil(t, ",")
}

func TestMountOptionFSNameEvilSpace(t *testing.T) {
	testMountOptionFSNameEvil(t, " ")
}

func TestMountOptionFSNameEvilTab(t *testing.T) {
	testMountOptionFSNameEvil(t, "\t")
}

func TestMountOptionFSNameEvilNewline(t *testing.T) {
	testMountOptionFSNameEvil(t, "\n")
}

func TestMountOptionFSNameEvilBackslash(t *testing.T) {
	testMountOptionFSNameEvil(t, `\`)
}

func TestMountOptionFSNameEvilBackslashDouble(t *testing.T) {
	// catch double-unescaping, if it were to happen
	testMountOptionFSNameEvil(t, `\\`)
}

func TestMountOptionSubtype(t *testing.T) {
	if runtime.GOOS == "darwin" {
		t.Skip("OS X does not support Subtype")
	}
	if runtime.GOOS == "freebsd" {
		t.Skip("FreeBSD does not support Subtype")
	}
	t.Parallel()
	const name = "FuseTestMarker"
	mnt, err := fstestutil.MountedT(t, fstestutil.SimpleFS{fstestutil.Dir{}},
		fuse.Subtype(name),
	)
	if err != nil {
		t.Fatal(err)
	}
	defer mnt.Close()

	info, err := fstestutil.GetMountInfo(mnt.Dir)
	if err != nil {
		t.Fatal(err)
	}
	if g, e := info.Type, "fuse."+name; g != e {
		t.Errorf("wrong Subtype: %q != %q", g, e)
	}
}

// TODO test LocalVolume

// TODO test AllowOther; hard because needs system-level authorization

func TestMountOptionAllowOtherThenAllowRoot(t *testing.T) {
	t.Parallel()
	mnt, err := fstestutil.MountedT(t, fstestutil.SimpleFS{fstestutil.Dir{}},
		fuse.AllowOther(),
		fuse.AllowRoot(),
	)
	if err == nil {
		mnt.Close()
	}
	if g, e := err, fuse.ErrCannotCombineAllowOtherAndAllowRoot; g != e {
		t.Fatalf("wrong error: %v != %v", g, e)
	}
}

// TODO test AllowRoot; hard because needs system-level authorization

func TestMountOptionAllowRootThenAllowOther(t *testing.T) {
	t.Parallel()
	mnt, err := fstestutil.MountedT(t, fstestutil.SimpleFS{fstestutil.Dir{}},
		fuse.AllowRoot(),
		fuse.AllowOther(),
	)
	if err == nil {
		mnt.Close()
	}
	if g, e := err, fuse.ErrCannotCombineAllowOtherAndAllowRoot; g != e {
		t.Fatalf("wrong error: %v != %v", g, e)
	}
}

type unwritableFile struct{}

func (f unwritableFile) Attr() fuse.Attr { return fuse.Attr{Mode: 0000} }

func TestMountOptionDefaultPermissions(t *testing.T) {
	if runtime.GOOS == "freebsd" {
		t.Skip("FreeBSD does not support DefaultPermissions")
	}
	t.Parallel()
	mnt, err := fstestutil.MountedT(t,
		fstestutil.SimpleFS{
			fstestutil.ChildMap{"child": unwritableFile{}},
		},
		fuse.DefaultPermissions(),
	)

	if err != nil {
		t.Fatal(err)
	}
	defer mnt.Close()

	// This will be prevented by kernel-level access checking when
	// DefaultPermissions is used.
	f, err := os.OpenFile(mnt.Dir+"/child", os.O_WRONLY, 0000)
	if err == nil {
		f.Close()
		t.Fatal("expected an error")
	}
	if !os.IsPermission(err) {
		t.Fatalf("expected a permission error, got %T: %v", err, err)
	}
}
