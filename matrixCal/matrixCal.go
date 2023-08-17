package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

const matrixSize = 3

var (
	matrixA   [matrixSize][matrixSize]int
	matrixB   [matrixSize][matrixSize]int
	result    [matrixSize][matrixSize]int
	rwLock    = sync.RWMutex{}
	waitGroup = sync.WaitGroup{}
	condition = sync.NewCond(rwLock.RLocker())
)

func generateMatrix(matrix *[matrixSize][matrixSize]int) {
	for i := 0; i < matrixSize; i++ {
		for j := 0; j < matrixSize; j++ {
			matrix[i][j] = rand.Intn(10) - 5
		}
	}
}

func workOut(row int) {

	rwLock.RLock()
	for {
		waitGroup.Done()
		condition.Wait()
		for col := 0; col < matrixSize; col++ {
			for i := 0; i < matrixSize; i++ {
				result[row][col] += matrixA[row][i] * matrixB[i][col]
			}
		}
	}

}

func main() {
	fmt.Println("working...")
	start := time.Now()
	waitGroup.Add(matrixSize)
	for row := 0; row < matrixSize; row++ {
		go workOut(row)
	}
	for i := 0; i < 100; i++ {
		waitGroup.Wait()
		rwLock.Lock()
		generateMatrix(&matrixA)
		generateMatrix(&matrixB)
		waitGroup.Add(matrixSize)
		rwLock.Unlock()
		condition.Broadcast()
	}
	end := time.Since(start)
	fmt.Printf("time cost = %v\n", end)
	fmt.Println(result)

}
