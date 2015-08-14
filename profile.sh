go test -run=none -bench=BenchmarkVerto_Param -c
go test -run=none -bench=BenchmarkVerto_Param -memprofile=verto.mprof
go test -run=none -bench=BenchmarkVerto_Param -cpuprofile=verto.prof
go tool pprof --dot verto.test verto.prof > verto.dot
dot -Tpng verto.dot -o verto.png
