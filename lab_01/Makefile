.PHONY: style ph_bench

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

