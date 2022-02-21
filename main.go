package main

import (
	"flag"
	"log"
	"os"
	"syscall"
	"unsafe"

	"github.com/paskozdilar/go-v4l2/v4l2"
)

func main() {
	// Configure logger
	log.SetFlags(log.Flags() | log.Lshortfile)

	// Parse device path
	var devPath string
	flag.StringVar(&devPath, "dev-path", "/dev/video0", "Video4linux device path")
	flag.Parse()

	// Open device
	f, err := os.Create(devPath)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()
	log.Println("- open")
	fd := f.Fd()

	// Query capabilities
	capability := v4l2.Capability{}
	err = ioctl(fd, v4l2.Vidioc_QueryCap, uintptr(unsafe.Pointer(&capability)))
	if err != nil {
		log.Fatal(err)
	}
	log.Println("- query capabilities")
	if (capability.Capabilities & v4l2.Cap_VideoCapture) == 0 {
		log.Fatal(devPath, "is not a capture device")
	}
	if (capability.Capabilities & v4l2.Cap_Streaming) == 0 {
		log.Fatal(devPath, "does not support streaming I/O")
	}

	// Reset crop to default
	cropCap := v4l2.CropCap{
		Type: v4l2.BufType_VideoCapture,
	}
	err = ioctl(fd, v4l2.Vidioc_CropCap, uintptr(unsafe.Pointer(&cropCap)))
	if err != nil {
		log.Println("Ignoring error [CropCap]:", err)
		// ignore
	} else {
		log.Println("- query crop cap")
		crop := v4l2.Crop{
			Type: v4l2.BufType_VideoCapture,
			C:    cropCap.DefRect,
		}
		err = ioctl(fd, v4l2.Vidioc_SCrop, uintptr(unsafe.Pointer(&crop)))
		if err != nil {
			log.Println("Ignoring error [CropCap]:", err)
			// ignore
		} else {
			log.Println("- set crop cap")
		}
	}

	// Negotiate format
	format := v4l2.Format{
		Type: v4l2.BufType_VideoCapture,
		Fmt: *((*[200]byte)(unsafe.Pointer(&v4l2.PixFormat{
			Width:       1920,
			Height:      1080,
			PixelFormat: v4l2.PixFmt_Mjpeg,
		}))),
	}
	err = ioctl(fd, v4l2.Vidioc_SFmt, uintptr(unsafe.Pointer(&format)))
	if err != nil {
		log.Fatal(err)
	}
	pixFormat := *((*v4l2.PixFormat)(unsafe.Pointer(&format.Fmt)))
	log.Println("- negotiated format:")
	log.Println("  + width =", pixFormat.Width)
	log.Println("  + height =", pixFormat.Height)

	// Initialize DMABUF
	reqBuf := v4l2.RequestBuffers{
		Count:  4,
		Type:   v4l2.BufType_VideoCapture,
		Memory: v4l2.Memory_Dmabuf,
	}
	err = ioctl(fd, v4l2.Vidioc_Reqbufs, uintptr(unsafe.Pointer(&reqBuf)))
	if err != nil {
		log.Fatal(err)
	}
	if reqBuf.Count < 2 {
		log.Fatal("Deriver returned buffer count", reqBuf.Count)
	}
	log.Println("- requested buffers:", reqBuf.Count)
	defer func() {
		err := ioctl(fd, v4l2.Vidioc_Reqbufs, uintptr(unsafe.Pointer(&v4l2.RequestBuffers{
			Count:  0,
			Type:   v4l2.BufType_VideoCapture,
			Memory: v4l2.Memory_Dmabuf,
		})))
		if err != nil {
			log.Println("Cleanup error:", err)
		}
	}()

	// Map buffers
}

func ioctl(fd, request uintptr, args ...uintptr) error {
	var (
		arg uintptr
		err error
	)
	if len(args) == 1 {
		arg = args[0]
	}
	_, _, errno := syscall.Syscall(syscall.SYS_IOCTL, fd, request, arg)
	if errno != 0 {
		err = errno
	}
	return err
}
