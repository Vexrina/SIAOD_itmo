bt_bench:
	go test ./btree -benchmem -count=100 -benchtime=60s -bench=.  > ./btree/benchs
	benchstat ./btree/benchs > ./btree/benchstat_res

bt_cpu_profiler:
	@mkdir -p ./profilers_pdf/bt/cpu
	@for bench in $$(go test -list=. ./btree | grep -v "ok"); do \
		go test -bench=^$$bench$$ -cpuprofile=./btree/cpu_$$bench.prof ./btree; \
		go tool pprof -pdf ./btree/cpu_$$bench.prof > ./profilers_pdf/bt/cpu/bt_cpu_$$bench.pdf; \
	done

bt_memory_profiler:
	@mkdir -p ./profilers_pdf/bt/memory
	@for bench in $$(go test -list=. ./btree | grep -v "ok"); do \
		go test -bench=^$$bench$$ -memprofile=./btree/mem_$$bench.prof ./btree; \
		go tool pprof -pdf ./btree/mem_$$bench.prof > ./profilers_pdf/bt/memory/bt_mem_$$bench.pdf; \
	done


kd_bench:
	go test ./kdtree -benchmem -count=10 -benchtime=60s -bench=.  > ./kdtree/benchs
	benchstat ./kdtree/benchs > ./kdtree/benchstat_res

kd_cpu_profiler:
	@mkdir -p ./profilers_pdf/kd/cpu
	@for bench in $$(go test -list=. ./kdtree | grep -v "ok"); do \
		go test -bench=^$$bench$$ -cpuprofile=./kdtree/cpu_$$bench.prof ./kdtree; \
		go tool pprof -pdf ./kdtree/cpu_$$bench.prof > ./profilers_pdf/kd/cpu/kd_cpu_$$bench.pdf; \
	done

kd_memory_profiler:
	@mkdir -p ./profilers_pdf/kd/memory
	@for bench in $$(go test -list=. ./kdtree | grep -v "ok"); do \
		go test -bench=^$$bench$$ -memprofile=./kdtree/mem_$$bench.prof ./kdtree; \
		go tool pprof -pdf ./kdtree/mem_$$bench.prof > ./profilers_pdf/kd/memory/kd_mem_$$bench.pdf; \
	done