/*
Naive code for multiplying two matrices together.

There must be a better way!
*/

#include <stdio.h>
#include <stdlib.h>

/*
  A naive implementation of matrix multiplication.

  DO NOT MODIFY THIS FUNCTION, the tests assume it works correctly, which it
  currently does
*/
void matrix_multiply(double **C, double **A, double **B, int a_rows, int a_cols,
                     int b_cols) {
  for (int i = 0; i < a_rows; i++) {
    for (int j = 0; j < b_cols; j++) {
      C[i][j] = 0;
      for (int k = 0; k < a_cols; k++)
        C[i][j] += A[i][k] * B[k][j];
    }
  }
}

/*
  This loop sequence is A rows, A cols, B cols. The inner loop proceeds along b in row-order,
  avoiding cache misses that are present in column-order.
*/
void fast_matrix_multiply(double **c, double **a, double **b, int a_rows,
                          int a_cols, int b_cols) {
  for (int i = 0; i < a_rows; i++) {
    for (int k = 0; k < a_cols; k++) {
      double a_val = a[i][k]; // avoid additional loads since this is the same on each inner loop
      for (int j = 0; j < b_cols; j++)
        c[i][j] += a_val * b[k][j];
    }
  }
}
