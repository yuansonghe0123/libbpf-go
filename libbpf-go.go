package main

/*
//#include <stdio.h>
//#include <errno.h>
#include <stdlib.h>
//#include <string.h>
//#include <stdarg.h>
//#include <sys/resource.h>
//#include <sys/syscall.h>
//#include <unistd.h>
#include <bpf/libbpf.h>
*/
import "C"
import (
	"bytes"
	"encoding/binary"
	"fmt"
)

func main() {
	module := NewModuleFromFile("./main.bpf.o")
	defer module.Close()
	module.AddElfFile("/root/main")
	err := module.LoadAllMap()
	if err != nil {
		panic(err)
	}
	err = module.LoadAllProgram()
	if err != nil {
		panic(err)
	}
	err = module.LoadAllIntoKernel()
	if err != nil {
		panic(err)
	}
	module.NewPerfBuffer("events")
	module.PerfStart("events")
	select {}
}

func perfOutput(cpu int, data []byte) {
	var a FuncEntryRecord
	if err := binary.Read(bytes.NewBuffer(data), binary.LittleEndian, &a); err != nil {
		panic(err)
	}
	fmt.Printf("GOID:%d TIME:%d FuncAddr:%x\n", a.GoId, a.KTime, a.FuncAddr)
}

func ringBufferOutput(data []byte) {

}

type FuncEntryRecord struct {
	GoId     uint64
	KTime    uint64
	FuncAddr uint32
}
