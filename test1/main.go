// Опыт 1. «Выделить — не значит получить»
// Демонстрируем ленивое выделение памяти (demand paging) в Linux.
package main

import (
	"fmt"
	"os"
	"strings"
	"syscall"
	"time"
)

func printMem(stage string) {
	status, _ := os.ReadFile("/proc/self/status")
	rssKB := 0
	for _, line := range strings.Split(string(status), "\n") {
		if strings.HasPrefix(line, "VmRSS:") {
			fmt.Sscanf(strings.TrimSpace(strings.TrimPrefix(line, "VmRSS:")), "%d", &rssKB)
		}
	}

	stat, _ := os.ReadFile("/proc/self/stat")
	tail := ""
	if idx := strings.LastIndex(string(stat), ")"); idx >= 0 {
		tail = string(stat[idx+1:])
	}
	fields := strings.Fields(tail)
	minor := fields[7] // 0-based: state=0, ppid=1, pgrp=2, session=3, tty=4, tpgid=5, flags=6, min_flt=7

	fmt.Printf("%-28s RSS: %6.1f МБ | minor faults: %s\n", stage, float64(rssKB)/1024, minor)
}

func main() {
	const blockSize = 1 << 30
	pageSize := os.Getpagesize()

	fmt.Printf("Размер страницы: %d байт\n", pageSize)
	printMem("Начало")

	block := make([]byte, blockSize)
	if len(os.Args) > 1 && os.Args[1] == "madvise" {
		if err := syscall.Madvise(block, syscall.MADV_NOHUGEPAGE); err != nil {
			fmt.Println("madvise(MADV_NOHUGEPAGE): ошибка:", err)
		} else {
			fmt.Println("madvise(MADV_NOHUGEPAGE): выполнен успешно")
		}
	}
	printMem("После выделения")

	start := time.Now()
	for i := 0; i < blockSize; i += pageSize {
		block[i] = 1
	}
	elapsed1 := time.Since(start)
	printMem("После первой записи")
	fmt.Printf("Первая запись заняла: %v\n", elapsed1)

	start = time.Now()
	for i := 0; i < blockSize; i += pageSize {
		block[i] = 2
	}
	elapsed2 := time.Since(start)
	printMem("После второй записи")
	fmt.Printf("Вторая запись заняла: %v\n", elapsed2)

	fmt.Printf("Вторая запись быстрее в %.1f раз\n", float64(elapsed1)/float64(elapsed2))

	fmt.Println("block[12345] =", block[12345])
}
