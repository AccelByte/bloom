# Bloom

This project bloom filter implementation using Murmur3 hash by [github.com/spaolacci/murmur3](github.com/spaolacci/murmur3).

## Usage

### Importing package

```go
import "github.com/AccelByte/bloom"
```

### Creating new bloom filter

```go
// create new filter with size of 100
// with default Murmur3 hashing strategy
// and 1.e-5 expected false positive probability
b := bloom.New(100)
```

### Putting item into bloom filter

```go
b.Put([]byte("an_item"))
```

### Checking if an item exists

```gp
b.MightContain([]byte("an_item"))
```

### Exporting bloom filter to JSON

```go
exported, _ := b.MarshalJSON()
```

### Constructing bloom filter from exported JSON

```go
bloomFilterJSON := &bloom.FilterJSON{}
json.Unmarshal(exported, bloomFilterJSON)
newB := bloom.From(bloomFilterJSON.B, bloomFilterJSON.K)
```