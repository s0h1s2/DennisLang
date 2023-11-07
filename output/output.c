// Declerations
#include <stdbool.h>
#include <stdio.h>
#include <stdlib.h>
typedef char i8;
typedef short i16;
typedef int i32;
typedef long i64;
typedef char *string;
typedef struct Point {
  i8 x;
  i8 y;
} Point;
typedef struct Location {
  i32 line;
  i32 col;
} Location;
typedef struct Token {
  Location loc;
  string value;
} Token;
i32 main() {
  i32 a = 1;
  bool b = true;
  Point point = (Point){.x = 1, .y = 2};
  string str = "Hello,world!\n";
  Token token = (Token){.loc = (Location){.line = 1, .col = 1},
                        .value = "Hello,world!\n"};
  token = (Token){.loc = (Location){.line = 1, .col = 1},
                  .value = "Hello,world!\n"};
  printf(str);
  point.x = 1;
}
