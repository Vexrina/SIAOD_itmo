style:
	go fmt ./...

ph_bench:
	go test -benchmem -bench . ./perfect_hash > ./perfect_hash/benchs

ph_cover:
	go test -cover ./perfect_hash

ph_cpu_profiler:
	go test -bench=. -cpuprofile=./perfect_hash/cpu.prof ./perfect_hash
	go tool pprof -pdf ./perfect_hash/cpu.prof > ph_cpu_profiler.pdf

ph_memory_profiler:
	go test -bench=. -memprofile=./perfect_hash/mem.prof ./perfect_hash
	go tool pprof -pdf ./perfect_hash/mem.prof > ph_mem_profiler.pdf

eh_bench:
	go test -benchmem -bench . ./extendible_hashing > ./extendible_hashing/benchs

eh_cover:
	go test -cover ./extendible_hashing

eh_cpu_profiler:
	go test -bench=. -cpuprofile=./extendible_hashing/cpu.prof ./extendible_hashing
	go tool pprof -pdf ./extendible_hashing/cpu.prof > eh_cpu_profiler.pdf

eh_memory_profiler:
	go test -bench=. -memprofile=./extendible_hashing/mem.prof ./extendible_hashing
	go tool pprof -pdf ./extendible_hashing/mem.prof > eh_mem_profiler.pdf

