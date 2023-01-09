ilog2_test: ilog2_test.c
	gcc -o $@ $^ -lmpfr -lgmp -lm

ilog2_ref: ilog2_ref.c
	gcc -o $@ $^ -lmpfr -lgmp

clzll_zero: clzll_zero.c
	gcc -o $@ $^

clean:
	@rm -f ilog2_test ilog2_ref clzll_zero
