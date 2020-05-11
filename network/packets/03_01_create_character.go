package packets

import (
	"bytes"

	. "../structures"
)

type CreateCharacterRequest struct {
	CharInfo CharacterInfo
}

func (p *CreateCharacterRequest) unmarshal(buffer *bytes.Buffer) error {
	p.CharInfo = CharacterInfo{}
	p.CharInfo.Read(buffer)

	return nil
}
