#include <math.h>
#include <stdint.h>
#include <stdio.h>
#include <stdlib.h>
#include <string.h>

#include <gmp.h>
#include <mpfr.h>

int ilog2(uint64_t x) { return 63 - __builtin_clzll(x); }

int ilog2b(uint64_t x) {
  double f = x;
  uint64_t v;
  memcpy(&v, &f, 8);
  return (v >> 52) - 1023;
}

int ilog2c(uint64_t x) {
  static const uint64_t debruijn_magic = 0x022fdd63cc95386dULL;

  static const uint64_t magic_table[] = {
      0,  1,  2,  53, 3,  7,  54, 27, 4,  38, 41, 8,  34, 55, 48, 28,
      62, 5,  39, 46, 44, 42, 22, 9,  24, 35, 59, 56, 49, 18, 29, 11,
      63, 52, 6,  26, 37, 40, 33, 47, 61, 45, 43, 21, 23, 58, 17, 10,
      51, 25, 36, 32, 60, 20, 57, 16, 50, 31, 19, 15, 30, 14, 13, 12,
  };

  x |= (x >> 1);
  x |= (x >> 2);
  x |= (x >> 4);
  x |= (x >> 8);
  x |= (x >> 16);
  x |= (x >> 32);
  return (magic_table[((x & ~(x >> 1)) * debruijn_magic) >> 58]);
}

int ilog2d(uint64_t x) {
  union {
    double d;
    uint32_t u[2];
  } ff;

  ff.d = x | 1;
  return (ff.u[1] >> 20) - 1023;
}

int ilog2_double(uint64_t x) { return (int)(log2(x)); }

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

void check(uint64_t x, int (*f)(uint64_t), int (*ref)(uint64_t)) {
  int got, want;
  got = f(x);
  want = ref(x);
  if (got != want) {
    printf("x=%lx, got=%d, want=%d, (uint64_t)(double)x=%lx\n", x, got, want,
           (uint64_t)(double)x);
  }
}

void check_fn(int (*f)(uint64_t)) {
  for (int x = 1; x < 100000; x++) {
    check(x, f, ilog2_ref);
  }

  for (int i = 1; i < 64; i++) {
    uint64_t x = (uint64_t)1 << i;
    check(x - 1, f, ilog2_ref);
    check(x, f, ilog2_ref);
    check(x + 1, f, ilog2_ref);
  }
  check(0xffffffffffffffff, f, ilog2_ref);
}

int main() {
  printf("ilog2 -----------\n");
  check_fn(ilog2);

  printf("ilog2b -----------\n");
  check_fn(ilog2b);

  printf("ilog2c -----------\n");
  check_fn(ilog2c);

  printf("ilog2d -----------\n");
  check_fn(ilog2d);

  printf("ilog2_double -----------\n");
  check_fn(ilog2_double);
}
