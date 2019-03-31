#include <stdio.h>
#include <string.h>
#include <stdlib.h>
#include "libtupass.h"

int main(int argc, char **argv) {
    char *password;
    GoFloat64 result;

    if (argc < 2) {
        fprintf(stderr, "missing argument\n");
        return 1;
    }

    password = strdup(argv[1]);
    result = CalculateStrength(password);


    printf("%s: %f\n", argv[1], result);

    return 0;
}