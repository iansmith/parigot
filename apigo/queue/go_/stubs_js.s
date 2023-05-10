#include "textflag.h"

TEXT queuesvc·create_handler(SB), NOSPLIT, $0
  CallImport
  RET

TEXT queuesvc·delete_handler(SB), NOSPLIT, $0
  CallImport
  RET

TEXT queuesvc·mark_done_handler(SB), NOSPLIT, $0
  CallImport
  RET

TEXT queuesvc·length_handler(SB), NOSPLIT, $0
  CallImport
  RET

TEXT queuesvc·send_handler(SB), NOSPLIT, $0
  CallImport
  RET

TEXT queuesvc·receive_handler(SB), NOSPLIT, $0
  CallImport
  RET

TEXT queuesvc·locate_handler(SB), NOSPLIT, $0
  CallImport
  RET
