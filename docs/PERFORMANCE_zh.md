## 性能提升

`静态模式`比`普通模式`提升性能`20%`以上。

使用两种模式对`routine`进行基准测试结果如下：

### :computer:普通模式

```text
> go test -a -bench .

*** GOOS: windows ***
*** GOARCH: amd64 ***
...
pkg: github.com/timandy/routine
cpu: AMD Ryzen 7 8845HS w/ Radeon 780M Graphics
...
BenchmarkGoid-4                                 382125445                3.088 ns/op           0 B/op          0 allocs/op
BenchmarkThreadLocal-4                          23455594                48.58 ns/op            7 B/op          0 allocs/op
BenchmarkThreadLocalWithInitial-4               22530049                50.08 ns/op            7 B/op          0 allocs/op
BenchmarkInheritableThreadLocal-4               21790962                49.64 ns/op            7 B/op          0 allocs/op
BenchmarkInheritableThreadLocalWithInitial-4    20939228                49.27 ns/op            7 B/op          0 allocs/op
BenchmarkGohack-4                               547616635                2.098 ns/op           0 B/op          0 allocs/op
```

### :rocket:静态模式

```text
> go test -toolexec="routinex" -a -bench .

*** GOOS: windows ***
*** GOARCH: amd64 ***
...
pkg: github.com/timandy/routine
cpu: AMD Ryzen 7 8845HS w/ Radeon 780M Graphics
...
BenchmarkGoid-4                                 599099253                1.979 ns/op           0 B/op          0 allocs/op
BenchmarkThreadLocal-4                          42951482                29.98 ns/op            7 B/op          0 allocs/op
BenchmarkThreadLocalWithInitial-4               41993280                30.54 ns/op            7 B/op          0 allocs/op
BenchmarkInheritableThreadLocal-4               45233858                28.50 ns/op            7 B/op          0 allocs/op
BenchmarkInheritableThreadLocalWithInitial-4    41149017                29.85 ns/op            7 B/op          0 allocs/op
BenchmarkGohack-4                               927123459                1.272 ns/op           0 B/op          0 allocs/op
```
