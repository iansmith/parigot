#include "textflag.h"

TEXT filesvc·open(SB), NOSPLIT, $0
  CallImport
  RET

TEXT filesvc·load(SB), NOSPLIT, $0
  CallImport
  RET
