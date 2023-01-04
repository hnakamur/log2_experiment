#include <stdint.h>
#include <stdio.h>
#include <string.h>

#include <gmp.h>
#include <mpfr.h>

int ilog2_ref(uint64_t x) {
  mpfr_t hx, hy;
  int y;

  mpfr_init2(hx, 20);
  mpfr_init2(hy, 20);

  mpfr_set_ui(hx, (unsigned long int)x, MPFR_RNDD);
  mpfr_log2(hy, hx, MPFR_RNDD);

  y = (int)mpfr_get_si(hy, MPFR_RNDD);
  mpfr_clear(hx);
  mpfr_clear(hy);
  mpfr_free_cache();
  return y;
}

int main() {
  for (int x = 1; x < 100000; x++) {
    printf("%u %d\n", x, ilog2_ref(x));
  }

  for (int i = 1; i < 64; i++) {
    uint64_t x = (uint64_t)1 << i;
    printf("%lu %d\n", x - 1, ilog2_ref(x - 1));
    printf("%lu %d\n", x, ilog2_ref(x));
    printf("%lu %d\n", x + 1, ilog2_ref(x + 1));
  }
  printf("%lu %d\n", 0xffffffffffffffff, ilog2_ref(0xffffffffffffffff));
}
