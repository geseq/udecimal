**Summary**

Based on https://github.com/robaho/fixed

A fixed place *unsigned* numeric library designed for performance.

All numbers have a fixed 8 decimal places, and the maximum permitted value is +- 99999999999,
or just under 100 billion.

The library is safe for concurrent use. It has built-in support for binary and json marshalling.

It is ideally suited for high performance trading financial systems. All common math operations are completed with 0 allocs.

**Design Goals**

Primarily developed to improve performance in [go-trader](https://github.com/robaho/go-trader).
Using Decimal rather than decimal.Decimal improves the performance by over 20%, and a lot less GC activity as well.

The decimal.Decimal API uses NaN for reporting errors in the common case, since often code is chained like:
```
result := someDecimal.Mul(NewS("123.50"))
```
and this would be a huge pain with error handling. Since all operations involving a NaN result in a NaN,
any errors quickly surface anyway.


**Performance**

<pre>
BenchmarkAddDecimal-8                   2000000000           0.87 ns/op        0 B/op          0 allocs/op
BenchmarkAddShopspringDecimal-8          5000000           247 ns/op         176 B/op          8 allocs/op
BenchmarkAddBigInt-8                    100000000           16.4 ns/op         0 B/op          0 allocs/op
BenchmarkAddBigFloat-8                  20000000            85.9 ns/op        48 B/op          1 allocs/op
BenchmarkMulDecimal-8                   300000000            4.07 ns/op        0 B/op          0 allocs/op
BenchmarkMulShopspringDecimal-8         20000000            75.4 ns/op        80 B/op          2 allocs/op
BenchmarkMulBigInt-8                    100000000           19.2 ns/op         0 B/op          0 allocs/op
BenchmarkMulBigFloat-8                  50000000            39.3 ns/op         0 B/op          0 allocs/op
BenchmarkDivDecimal-8                   300000000            5.22 ns/op        0 B/op          0 allocs/op
BenchmarkDivShopspringDecimal-8          2000000           792 ns/op         568 B/op         21 allocs/op
BenchmarkDivBigInt-8                    30000000            48.0 ns/op         8 B/op          1 allocs/op
BenchmarkDivBigFloat-8                  20000000           116 ns/op          24 B/op          2 allocs/op
BenchmarkCmpDecimal-8                   2000000000           0.42 ns/op        0 B/op          0 allocs/op
BenchmarkCmpShopspringDecimal-8         200000000            8.55 ns/op        0 B/op          0 allocs/op
BenchmarkCmpBigInt-8                    200000000            6.12 ns/op        0 B/op          0 allocs/op
BenchmarkCmpBigFloat-8                  300000000            5.33 ns/op        0 B/op          0 allocs/op
BenchmarkStringDecimal-8                20000000            61.5 ns/op        32 B/op          1 allocs/op
BenchmarkStringNDecimal-8               20000000            60.7 ns/op        32 B/op          1 allocs/op
BenchmarkStringShopspringDecimal-8       5000000           242 ns/op          64 B/op          5 allocs/op
BenchmarkStringBigInt-8                 10000000           139 ns/op          24 B/op          2 allocs/op
BenchmarkStringBigFloat-8                3000000           449 ns/op         192 B/op          8 allocs/op
BenchmarkWriteTo-8                      50000000            42.9 ns/op        21 B/op          0 allocs/op
</pre>

The "decimal" above is the common [shopspring decimal](https://github.com/shopspring/decimal) library
