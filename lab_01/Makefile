style:
	go fmt ./...

ph_bench:
	go test -benchmem -bench . ./perfect_hash > ./perfect_hash/benchs

ph_cover:
	go test -cover ./perfect_hash

ph_cpu_profiler:
	go test -cpuprofile=./perfect_hash/cpu.prof ./perfect_hash

ph_memory_profiler:
	go test -memprofile=./perfect_hash/mem.prof ./perfect_hash

eh_bench:
	go test -benchmem -bench . ./extendible_hashing > ./extendible_hashing/benchs

eh_cover:
	go test -cover ./extendible_hashing

eh_cpu_profiler:
	go test -cpuprofile=./extendible_hashing/cpu.prof ./extendible_hashing

eh_memory_profiler:
	go test -memprofile=./extendible_hashing/mem.prof ./extendible_hashing
