package main

import (
    "crypto/sha256"
    "encoding/binary"
    "fmt"
    "sort"
)

func rightMostBitPos(x uint64) int {
    if x == 0 {
        return -1
    }
    pos := 0
    for (x & 1) == 0 {
        pos++
        x >>= 1
    }
    return pos
}

func hashItem(item string, seed int) uint64 {
    h := sha256.New()
    h.Write([]byte(fmt.Sprintf("%s-%d", item, seed)))
    bs := h.Sum(nil)
    return binary.BigEndian.Uint64(bs[:8])
}

func fmAlgo(items []string, numHashes int) int {
    estimates := make([]int, numHashes)

    for i := 0; i < numHashes; i++ {
        maxRho := 0
        for _, v := range items {
            hash := hashItem(v, i)
            rho := rightMostBitPos(hash)
            if rho > maxRho {
                maxRho = rho
            }
        }
        estimates[i] = 1 << maxRho
    }

    sort.Ints(estimates)
    median := estimates[numHashes/2]

    correctionFactor := 0.77351
    correctedEstimate := int(float64(median) / correctionFactor)

    return correctedEstimate
}

func main() {
    var items []string
    for i := 0; i < 1000; i++ {
        items = append(items, fmt.Sprintf("item_%d", i))
    }

    estimate := fmAlgo(items, 32)
    fmt.Printf("Estimated number of unique items: %d\n", estimate)
    fmt.Printf("Actual number of unique items: %d\n", len(items))
}
