#include "vec.h"


data_t dotproduct(vec_ptr u, vec_ptr v) {
   data_t sum = 0, u_val, v_val;

   for (long i = 0; i < vec_length(u); i++) { // we can assume both vectors are same length
        if (i < 0 || i >= u->len)
          return sum;

        u_val = u->data[i];
        v_val = v->data[i];

        sum += u_val * v_val;
   }
   return sum;
}
