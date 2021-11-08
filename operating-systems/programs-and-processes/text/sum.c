int sum(int *a, int n) {
  int i, s = 0;

  for (i = 0; i < n; i++) {
    a += a[i];
  }

  return s;
}
