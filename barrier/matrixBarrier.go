package main

import (
	"fmt"
	"math/rand"
	"time"
)

const matrixSize = 3

var (
	matrixA      [matrixSize][matrixSize]int
	matrixB      [matrixSize][matrixSize]int
	result       [matrixSize][matrixSize]int
	workStart    = NewBarrier(matrixSize)
	workComplete = NewBarrier(matrixSize)
)

func generateMatrix(matrix *[matrixSize][matrixSize]int) {
	for i := 0; i < matrixSize; i++ {
		for j := 0; j < matrixSize; j++ {
			matrix[i][j] = rand.Intn(10) - 5
		}
	}
}

func workOut(row int) {

	for {
		workStart.Wait()
		for col := 0; col < matrixSize; col++ {
			for i := 0; i < matrixSize; i++ {
				result[row][col] += matrixA[row][i] * matrixB[i][col]
			}
		}
		workComplete.Wait()
	}

}

func main() {
	fmt.Println("working...")
	start := time.Now()

	for row := 0; row < matrixSize; row++ {
		go workOut(row)
	}
	for i := 0; i < 100; i++ {
		generateMatrix(&matrixA)
		generateMatrix(&matrixB)
		workStart.Wait()
		workComplete.Wait()
	}
	end := time.Since(start)
	fmt.Printf("time cost = %v\n", end)
	fmt.Println(result)

}
