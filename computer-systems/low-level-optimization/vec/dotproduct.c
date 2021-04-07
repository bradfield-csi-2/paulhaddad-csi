#include "vec.h"


data_t dotproduct(vec_ptr u, vec_ptr v) {
  data_t sum = 0, u_val1, v_val1, u_val2, v_val2;

  long length = vec_length(u); // we're assuming u and v are the same lengths
  int limit = length-1; // unrolling factor
  long i;

  for (i = 0; i < limit; i+=2) { // we can assume both vectors are same length

    u_val1 = u->data[i];
    v_val1 = v->data[i];

    u_val2 = u->data[i+1];
    v_val2 = v->data[i+1];

    sum += u_val1 * v_val1;
    sum += u_val2 * v_val2;
  }

  for (; i < length; i++) {
    u_val1 = u->data[i];
    v_val1 = v->data[i];

    sum += u_val1 * v_val1;
  }

  return sum;
}
