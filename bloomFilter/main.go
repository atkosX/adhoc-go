package main

import (
    "crypto/sha256"
    "encoding/binary"
    "fmt"
    "math"
)

func calculateSize(n float64, p float64) float64 {
    return -1 * (n * math.Log(p)) / (math.Pow(math.Log(2), 2))
}

func calculateHashCount(m float64, n float64) float64 {
    return (m / n) * math.Log(2)
}


func hashWithSeed(s string, seed int) uint32 {
    data := []byte(fmt.Sprintf("%d-%s", seed, s))
    hash := sha256.Sum256(data)
    return binary.BigEndian.Uint32(hash[:4])
}

func main() {
    var n float64 = 1000000 
    var p float64 = 0.01  
    size := int(calculateSize(n, p))
    k := int(calculateHashCount(float64(size), n))

    bitArray := make([]int, size)

    for {
        fmt.Print("Enter a string: ")
        var input string
        fmt.Scan(&input)

        present := true
        for i := 0; i < k; i++ {
            index := int(hashWithSeed(input, i)) % size
            if bitArray[index] == 0 {
                present = false
            }
        }

        if present {
            fmt.Println("Possibly Present")
        } else {
            fmt.Println("Not Present")
            for i := 0; i < k; i++ {
                index := int(hashWithSeed(input, i)) % size
                bitArray[index] = 1
            }
        }
    }
}
