package items

import (
		"encoding/base64"
		"io"
)

type MobileCodeStrItem struct {
		Code   string
}

func (this *MobileCodeStrItem) WriteTo(w io.Writer) (n int64, err error) {
		nb, err := w.Write([]byte(this.Code))
		return int64(nb), err
}

func (this *MobileCodeStrItem) EncodeB64string() string {
		return base64.StdEncoding.EncodeToString([]byte(this.Code))
}
