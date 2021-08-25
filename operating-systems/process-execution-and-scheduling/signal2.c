#include <stdio.h>
#include <stdlib.h>
#include <unistd.h>

int main(int argc, char *argv[]) {
  int rc, status;

  if ((rc = fork()) == -1) {
    printf("error forking\n");
  } else if (rc == 0) {
    printf("in child process\n");
    pause();
    exit(0);
  } else {
    printf("in parent process, waiting on child: %d\n", rc);
    waitpid(rc, &status, 0);

    if WIFEXITED(status) {
      printf("child exited normally\n");
    }

    if WIFSIGNALED(status) {
      printf("child exited due to signal\n");
    }
  }

  return EXIT_SUCCESS;
}
