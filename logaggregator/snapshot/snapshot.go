package snapshot

import (
	"encoding/gob"
	"io"

	"github.com/flynn/flynn/pkg/syslog/rfc5424"
)

func Take(buffers [][]*rfc5424.Message, w io.Writer) error {
	enc := gob.NewEncoder(w)

	for _, buf := range buffers {
		for _, msg := range buf {
			if err := enc.Encode(msg); err != nil {
				return err
			}
		}
	}
	return nil
}

type Scanner struct {
	dec *gob.Decoder
	err error

	Message *rfc5424.Message
}

func NewScanner(r io.Reader) *Scanner {
	return &Scanner{dec: gob.NewDecoder(r)}
}

func (s *Scanner) Scan() bool {
	msg := &rfc5424.Message{}
	err := s.dec.Decode(msg)
	if err != nil {
		s.err = err
		return false
	}
	s.Message = msg
	return true
}

func (s *Scanner) Err() error {
	if s.err != io.EOF {
		return s.err
	}
	return nil
}
