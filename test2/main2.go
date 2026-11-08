// Опыт 2. «Матрица по строкам и по столбцам»
// Одномерный массив N×N int32, элемент (i, j) по адресу i*N + j.
// Сравниваем скорость суммирования по строкам и по столбцам.
package main

import (
	"fmt"
	"time"
)

const N = 8192

func main() {
	total := N * N
	a := make([]int32, total)
	for i := 0; i < total; i++ {
		a[i] = int32(i % 100)
	}

	var rowSum, colSum int64
	const runs = 3

	var rowTimes, colTimes time.Duration

	for r := 0; r < runs; r++ {
		start := time.Now()
		rowSum = 0
		for i := 0; i < N; i++ {
			for j := 0; j < N; j++ {
				rowSum += int64(a[i*N+j])
			}
		}
		rowTimes += time.Since(start)
	}

	for r := 0; r < runs; r++ {
		start := time.Now()
		colSum = 0
		for j := 0; j < N; j++ {
			for i := 0; i < N; i++ {
				colSum += int64(a[i*N+j])
			}
		}
		colTimes += time.Since(start)
	}

	fmt.Println("Сумма по строкам:  ", rowSum)
	fmt.Println("Сумма по столбцам: ", colSum)
	fmt.Println("Суммы совпадают:   ", rowSum == colSum)
	fmt.Printf("По строкам:  среднее %.1f мс\n", rowTimes.Seconds()*1000/float64(runs))
	fmt.Printf("По столбцам: среднее %.1f мс\n", colTimes.Seconds()*1000/float64(runs))
	fmt.Printf("По столбцам медленнее в %.1f раз\n",
		float64(colTimes)/float64(rowTimes))
}
