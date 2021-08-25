#include <stdio.h>
#include <stdlib.h>
#include <signal.h>
#include <unistd.h>

void sigint_handler(int sig) {
  printf("received signal\n");
  exit(0);
}

int main(int arc, char *argv[]) {
  if (signal(SIGINT, sigint_handler) == SIG_ERR) {
    printf("error installing signal\n");
    exit(1);
  }

  pause();

  printf("signal\n");
  return EXIT_SUCCESS;
}
