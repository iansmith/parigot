//go:build js && !browser

#include "textflag.h"

TEXT parigot·locate_(SB), NOSPLIT, $0
  CallImport
  RET


TEXT parigot·dispatch_(SB), NOSPLIT, $0
  CallImport
  RET

TEXT parigot·bind_method_(SB), NOSPLIT, $0
  CallImport
  RET

TEXT parigot·block_until_call_(SB), NOSPLIT, $0
  CallImport
  RET

TEXT parigot·return_value_(SB), NOSPLIT, $0
  CallImport
  RET

TEXT parigot·exit_(SB), NOSPLIT, $0
  CallImport
  RET

TEXT parigot·require_(SB), NOSPLIT, $0
  CallImport
  RET

TEXT parigot·export_(SB), NOSPLIT, $0
  CallImport
  RET

TEXT parigot·run_(SB), NOSPLIT, $0
  CallImport
  RET

TEXT parigot·backdoor_log_(SB), NOSPLIT, $0
  CallImport
  RET

