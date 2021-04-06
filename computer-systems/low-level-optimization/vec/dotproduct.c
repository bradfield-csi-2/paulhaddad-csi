#include "vec.h"


data_t dotproduct(vec_ptr u, vec_ptr v) {
  data_t sum = 0, u_val, v_val;

  long length = vec_length(u); // we're assuming u and v are the same lengths

  for (long i = 0; i < length; i++) { // we can assume both vectors are same length

    u_val = u->data[i];
    v_val = v->data[i];

    sum += u_val * v_val;
  }
  return sum;
}
