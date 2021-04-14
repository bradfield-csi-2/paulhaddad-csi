  /*

Two different ways to loop over an array of arrays.

Spottei at:
http://stackoverflow.com/questions/9936132/why-does-the-order-of-the-loops-affect-performance-when-iterating-over-a-2d-arra

*/

void option_one() {
  int i, j;
  static int x[10000][10000];
  for (i = 0; i < 10000; i++) {
    for (j = 0; j < 10000; j++) {
      x[i][j] = i + j;
    }
  }
}

void option_two() {
  int i, j;
  static int x[10000][10000];
  for (i = 0; i < 10000; i++) {
    for (j = 0; j < 10000; j++) {
      x[j][i] = i + j;
    }
  }
}

int main() {
  option_one();
  option_two();
  return 0;
}
