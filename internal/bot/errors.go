package bot

import "errors"

var ErrMessageNotModified = errors.New("telegram: Bad Request: message is not modified: specified new message content and reply markup are exactly the same as a current content and reply markup of the message (400)")
