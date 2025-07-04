package main

import (
    "crypto/sha256"
    "encoding/binary"
    "fmt"
    "math"
)

type BloomFilter struct {
    bitArray  []int
    size      int
    hashCount int
}

func NewBloomFilter(n, p float64) *BloomFilter {
    size := int(calculateSize(n, p))
    k := int(calculateHashCount(float64(size), n))
    bitArray := make([]int, size)
    return &BloomFilter{bitArray: bitArray, size: size, hashCount: k}
}


func calculateSize(n, p float64) float64 {
    return -1 * (n * math.Log(p)) / (math.Pow(math.Log(2), 2))
}


func calculateHashCount(m, n float64) float64 {
    return (m / n) * math.Log(2)
}


func (bf *BloomFilter) hashWithSeed(s string, seed int) uint32 {
    data := []byte(fmt.Sprintf("%d-%s", seed, s))
    hash := sha256.Sum256(data)
    return binary.BigEndian.Uint32(hash[:4])
}


func (bf *BloomFilter) Add(s string) {
    for i := 0; i < bf.hashCount; i++ {
        idx := int(bf.hashWithSeed(s, i)) % bf.size
        bf.bitArray[idx] = 1
    }
}

func (bf *BloomFilter) Contains(s string) bool {
    present := 0
    for i := 0; i < bf.hashCount; i++ {
        idx := int(bf.hashWithSeed(s, i)) % bf.size
        if bf.bitArray[idx] == 1 {
            present++
        }
    }
    return present == bf.hashCount
}

func main() {
    var n float64 = 1000000
    var p float64 = 0.01  

    bf := NewBloomFilter(n, p)

    for {
        fmt.Print("Enter a string: ")
        var input string
        fmt.Scanln(&input)

        present := bf.Contains(input)
        if present {
            fmt.Println("Possibly Present")
        } else {
            fmt.Println("Not Present")
            bf.Add(input)
        }
    }
}
