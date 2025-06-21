package crypto

type Crypto interface {
	Encrypt(data []byte)
	Decrypt(data []byte)
}

type crypto562 struct {
	constKey1   int
	constKey2   int
	dynamicKey  int
	dynamicKey1 byte
	dynamicKey2 byte
	constKeyEn  uint32
	constKeyDe  uint32
}

func NewCrypto562(dynamicKey int) Crypto {
	return &crypto562{
		constKey1:   0x241AE7,
		constKey2:   0x15DCB2,
		dynamicKey:  dynamicKey,
		dynamicKey1: 0x02,
		dynamicKey2: 0x01,
		constKeyEn:  0xA7F0753B,
		constKeyDe:  0xAAF29BF3,
	}
}

func (c *crypto562) Decrypt(data []byte) {
	bufferLen := len(data)
	sOffset := 0x0C
	for i := sOffset; i+4 <= bufferLen; i += 4 {
		DynamicKey := c.dynamicKey
		for j := i; j < i+4; j++ {
			pSrc := data[j]
			data[j] = pSrc ^ byte(DynamicKey>>8)
			DynamicKey = (int(pSrc)+DynamicKey)*c.constKey1 + c.constKey2
		}
	}
}

func (c *crypto562) Encrypt(data []byte) {
	bufferLen := len(data)
	sOffset := 0x0C
	for i := sOffset; i+4 <= bufferLen; i += 4 {
		DynamicKey := c.dynamicKey
		for j := i; j < i+4; j++ {
			data[j] = data[j] ^ byte(DynamicKey>>8)
			DynamicKey = (int(data[j])+DynamicKey)*c.constKey1 + c.constKey2
		}
	}
}
