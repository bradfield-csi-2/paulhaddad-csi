#include "vec.h"


data_t dotproduct(vec_ptr u, vec_ptr v) {
  data_t sum1 = 0, sum2 = 0, u_val1, v_val1, u_val2, v_val2;

  long length = vec_length(u); // we're assuming u and v are the same lengths
  int limit = length-1; // unrolling factor
  long i;

  data_t *u_start = u->data;
  data_t *v_start = v->data;

  for (i = 0; i < limit; i+=2) { // we can assume both vectors are same length

    u_val1 = u_start[i];
    v_val1 = v_start[i];

    u_val2 = u_start[i+1];
    v_val2 = v_start[i+1];

    sum1 += u_val1 * v_val1;
    sum2 += u_val2 * v_val2;
  }

  for (; i < length; i++) {
    u_val1 = u_start[i];
    v_val1 = v_start[i];

    sum1 += u_val1 * v_val1;
  }

  return sum1 + sum2;
}
