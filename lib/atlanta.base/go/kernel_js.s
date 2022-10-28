#include "textflag.h"

TEXT parigot·locate_(SB), NOSPLIT, $0
  CallImport
  RET

TEXT parigot·register_(SB), NOSPLIT, $0
    CallImport
    //CALL go·parigot·register(SB)
  RET

TEXT parigot·dispatch_(SB), NOSPLIT, $0
  CallImport
  RET

TEXT parigot·exit_(SB), NOSPLIT, $0
  CallImport
  RET

