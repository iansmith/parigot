//go:build js && !browser

#include "textflag.h"

TEXT ·open(SB), NOSPLIT, $0
  CallImport
  RET

TEXT ·close(SB), NOSPLIT, $0
  CallImport
  RET


TEXT ·load_test_data(SB), NOSPLIT, $0
  CallImport
  RET
