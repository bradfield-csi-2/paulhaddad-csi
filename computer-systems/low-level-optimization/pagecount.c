#include <stdint.h>
#include <stdio.h>
#include <time.h>

#define TEST_LOOPS 10000000

uint64_t pagecount(uint64_t memory_size, int page_size_power) {
  return memory_size >> page_size_power;
}

int main (int argc, char** argv) {
  clock_t baseline_start, baseline_end, test_start, test_end;
  uint64_t memory_size, page_size;
  int page_size_power;
  double clocks_elapsed, time_elapsed;
  int i, ignore = 0;

  uint64_t msizes[] = {1L << 32, 1L << 40, 1L << 52};
  uint64_t psizes[] = {1L << 12, 1L << 16, 1L << 32};
  int p_powers[] = {12, 16, 32};

  baseline_start = clock();
  for (i = 0; i < TEST_LOOPS; i++) {
    int mod = i % 3;
    memory_size = msizes[mod];
    page_size = psizes[mod];
    ignore += 1 + memory_size +
              page_size; // so that this loop isn't just optimized away
  }
  baseline_end = clock();

  test_start = clock();
  for (i = 0; i < TEST_LOOPS; i++) {
    int mod = i % 3;
    memory_size = msizes[mod];
    page_size_power = p_powers[mod];
    ignore += pagecount(memory_size, page_size_power) + memory_size + page_size;
  }
  test_end = clock();

  clocks_elapsed = test_end - test_start - (baseline_end - baseline_start);
  time_elapsed = clocks_elapsed / CLOCKS_PER_SEC;

  printf("%.2fs to run %d tests (%.2fns per test)\n", time_elapsed, TEST_LOOPS,
         time_elapsed * 1e9 / TEST_LOOPS);
  return ignore;
}

/*
 * Expected instructions:
 *  Lines 22-23 and 31-32 will be memory reads
 *  The call to pagecount will have overhead, even though there are only
 *  9 possible results. We can calculate this upfront.
 *
 *  Line 8: The division operation will always be division since non-power of
 *  twos are possible. We know our data, so we can optimize this to a right
 *  shift.
 *
 *  The two loops must perform many operations in order to get the memory and
 *  page sizes.
 *
 *  Line 33: This loop has a function call within it.
 *
 *  The modulo operations are done through mult and shtl operations, followed by
 *  a memory read
 *
 *  Interestingly, getting page_size doesn't appear to have any operations.
 */
