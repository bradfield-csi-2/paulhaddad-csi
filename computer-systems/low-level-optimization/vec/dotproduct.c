#include "vec.h"


data_t dotproduct(vec_ptr u, vec_ptr v) {
  data_t sum1 = 0, sum2 = 0, sum3 = 0, sum4 = 0, sum5 = 0;
  data_t u_val1, v_val1, u_val2, v_val2, u_val3, v_val3, u_val4, v_val4, u_val5, v_val5;

  long length = vec_length(u); // we're assuming u and v are the same lengths
  int limit = length-4; // unrolling factor
  long i;

  data_t *u_start = u->data;
  data_t *v_start = v->data;

  for (i = 0; i < limit; i+=5) { // we can assume both vectors are same length

    u_val1 = u_start[i];
    v_val1 = v_start[i];

    u_val2 = u_start[i+1];
    v_val2 = v_start[i+1];

    u_val3 = u_start[i+2];
    v_val3 = v_start[i+2];

    u_val4 = u_start[i+3];
    v_val4 = v_start[i+3];

    u_val5 = u_start[i+4];
    v_val5 = v_start[i+4];

    sum1 += u_val1 * v_val1;
    sum2 += u_val2 * v_val2;
    sum3 += u_val3 * v_val3;
    sum4 += u_val4 * v_val4;
    sum5 += u_val5 * v_val5;
  }

  for (; i < length; i++) {
    u_val1 = u_start[i];
    v_val1 = v_start[i];

    sum1 += u_val1 * v_val1;
  }

  return sum1 + sum2 + sum3 + sum4 + sum5;
}
