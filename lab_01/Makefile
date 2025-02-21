style:
	go fmt ./...

ph_bench:
	go test -benchmem -count=10 -bench . ./perfect_hash > ./perfect_hash/benchs

ph_cover:
	go test -cover ./perfect_hash

ph_cpu_profiler:
	@mkdir -p ./profilers_pdf/ph/cpu
	@for bench in $$(go test -list=. ./perfect_hash | grep -v "ok"); do \
		go test -bench=^$$bench$$ -cpuprofile=./perfect_hash/cpu_$$bench.prof ./perfect_hash; \
		go tool pprof -pdf ./perfect_hash/cpu_$$bench.prof > ./profilers_pdf/ph/cpu/ph_cpu_$$bench.pdf; \
	done

ph_memory_profiler:
	@mkdir -p ./profilers_pdf/ph/memory
	@for bench in $$(go test -list=. ./perfect_hash | grep -v "ok"); do \
		go test -bench=^$$bench$$ -memprofile=./perfect_hash/mem_$$bench.prof ./perfect_hash; \
		go tool pprof -pdf ./perfect_hash/mem_$$bench.prof > ./profilers_pdf/ph/memory/ph_mem_$$bench.pdf; \
	done

eh_bench:
	go test -benchmem -count=7 -bench . ./extendible_hashing > ./extendible_hashing/benchs

eh_cover:
	go test -cover ./extendible_hashing

eh_cpu_profiler:
	@mkdir -p ./profilers_pdf/eh/cpu
	@for bench in $$(go test -list=. ./extendible_hashing | grep -v "ok"); do \
		go test -bench=^$$bench$$ -cpuprofile=./extendible_hashing/cpu_$$bench.prof ./extendible_hashing; \
		go tool pprof -pdf ./extendible_hashing/cpu_$$bench.prof > ./profilers_pdf/eh/cpu/eh_cpu_$$bench.pdf; \
	done

eh_memory_profiler:
	@mkdir -p ./profilers_pdf/eh/memory
	@for bench in $$(go test -list=. ./extendible_hashing | grep -v "ok"); do \
		go test -bench=^$$bench$$ -memprofile=./extendible_hashing/mem_$$bench.prof ./extendible_hashing; \
		go tool pprof -pdf ./extendible_hashing/mem_$$bench.prof > ./profilers_pdf/eh/memory/eh_mem_$$bench.pdf; \
	done

mh_bench:
	go test -benchmem -bench . ./min_hash > ./min_hash/benchs

mh_cover:
	go test -cover ./min_hash

mh_cpu_profiler:
	go test -bench=. -cpuprofile=./min_hash/cpu.prof ./min_hash
	go tool pprof -pdf ./min_hash/cpu.prof > ./profilers_pdf/mh_cpu_profiler.pdf

mh_memory_profiler:
	go test -bench=. -memprofile=./min_hash/mem.prof ./min_hash
	go tool pprof -pdf ./min_hash/mem.prof > ./profilers_pdf/mh_mem_profiler.pdf

mh_cpu_profiler:
	@mkdir -p ./profilers_pdf/min_hash/cpu
	@for bench in $$(go test -list=. ./min_hash | grep -v "ok"); do \
		go test -bench=^$$bench$$ -cpuprofile=./min_hash/cpu_$$bench.prof ./min_hash; \
		go tool pprof -pdf ./min_hash/cpu_$$bench.prof > ./profilers_pdf/min_hash/cpu/min_hash_cpu_$$bench.pdf; \
	done

mh_memory_profiler:
	@mkdir -p ./profilers_pdf/min_hash/memory
	@for bench in $$(go test -list=. ./min_hash | grep -v "ok"); do \
		go test -bench=^$$bench$$ -memprofile=./min_hash/mem_$$bench.prof ./min_hash; \
		go tool pprof -pdf ./min_hash/mem_$$bench.prof > ./profilers_pdf/min_hash/memory/min_hash_mem_$$bench.pdf; \
	done