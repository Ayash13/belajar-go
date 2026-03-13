# 8. Methods

## Apa itu Method?

Method adalah fungsi yang menempel pada sebuah type (biasanya struct). Go tidak memiliki class, jadi method menjadi cara untuk memberikan behaviour pada struct.

## Value Receiver

Menerima salinan (copy) dari struct. Perubahan di dalam method **tidak akan berpengaruh** pada struct aslinya.

```go
func (c Counter) Print() {
    fmt.Println(c.Value)
}
```

## Pointer Receiver

Menerima alamat memory struct. Perubahan di dalam method **akan mengubah** struct aslinya.

```go
func (c *Counter) Increment() {
    c.Value++
}
```

> **Kapan Pakai Pointer Receiver?**
>
> 1. Ingin mengubah state struct tersebut.
> 2. Menghindari copy variable yang ukurannya terlalu besar.
