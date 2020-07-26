package ngin

import (
	"time"

	"zero/core/service"
)

type (
	PrivateKeyConf struct {
		Fingerprint string
		KeyFile     string
	}

	SignatureConf struct {
		Strict      bool          `json:",default=false"`
		Expiry      time.Duration `json:",default=1h"`
		PrivateKeys []PrivateKeyConf
	}

	NgConf struct {
		service.ServiceConf
		Host     string `json:",default=0.0.0.0"`
		Port     int
		Verbose  bool  `json:",optional"`
		MaxConns int   `json:",default=10000"`
		MaxBytes int64 `json:",default=1048576,range=[0:8388608]"`
		// milliseconds
		Timeout      int64         `json:",default=3000"`
		CpuThreshold int64         `json:",default=900,range=[0:1000]"`
		Signature    SignatureConf `json:",optional"`
	}
)
