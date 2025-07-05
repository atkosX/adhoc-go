package main

import (
    "fmt"
    "hash/crc32"
    "strconv"
    "sync"
)

type ConsistentHash struct {
    replicas int
    hashes   []uint32
    nodes    map[uint32]string
    mu       sync.RWMutex
}

func New(replicas int) *ConsistentHash {
    return &ConsistentHash{
        replicas: replicas,
        nodes:    make(map[uint32]string),
    }
}

func (ch *ConsistentHash) Add(node string) {
    ch.mu.Lock()
    defer ch.mu.Unlock()
    for i := 0; i < ch.replicas; i++ {
        key := node + ":" + strconv.Itoa(i)
        h := crc32.ChecksumIEEE([]byte(key))
        ch.hashes = append(ch.hashes, h)
        ch.nodes[h] = node
    }
   
    for i := 0; i < len(ch.hashes); i++ {
        for j := i + 1; j < len(ch.hashes); j++ {
            if ch.hashes[i] > ch.hashes[j] {
                ch.hashes[i], ch.hashes[j] = ch.hashes[j], ch.hashes[i]
            }
        }
    }
}

func (ch *ConsistentHash) Remove(node string) {
    ch.mu.Lock()
    defer ch.mu.Unlock()
    filtered := ch.hashes[:0]
    for _, h := range ch.hashes {
        if ch.nodes[h] == node {
            delete(ch.nodes, h)
        } else {
            filtered = append(filtered, h)
        }
    }
    ch.hashes = filtered
}

func (ch *ConsistentHash) Get(key string) (string, bool) {
    ch.mu.RLock()
    defer ch.mu.RUnlock()
    if len(ch.hashes) == 0 {
        return "", false
    }

    target := crc32.ChecksumIEEE([]byte(key))
    i := ch.binarySearch(target)
    return ch.nodes[ch.hashes[i]], true
}

func (ch *ConsistentHash) binarySearch(target uint32) int {
    lo, hi := 0, len(ch.hashes)-1
    for lo <= hi {
        mid := lo + (hi-lo)/2
        if ch.hashes[mid] == target {
            return mid
        }
        if ch.hashes[mid] < target {
            lo = mid + 1
        } else {
            hi = mid - 1
        }
    }
    if lo >= len(ch.hashes) {
        return 0
    }
    return lo
}

func main() {
    ch := New(3)

    ch.Add("A")
    ch.Add("B")
    ch.Add("C")

    keys := []string{"apple", "banana", "cherry", "date"}
    fmt.Println("Initial assignments:")
    for _, k := range keys {
        node, _ := ch.Get(k)
        fmt.Printf("  %s -> %s\n", k, node)
    }

    fmt.Println("\nRemove B…")
    ch.Remove("B")
    for _, k := range keys {
        node, _ := ch.Get(k)
        fmt.Printf("  %s -> %s\n", k, node)
    }

    fmt.Println("\nAdd D…")
    ch.Add("D")
    for _, k := range keys {
        node, _ := ch.Get(k)
        fmt.Printf("  %s -> %s\n", k, node)
    }
}
