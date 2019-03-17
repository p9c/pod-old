package sub

// Reed Solomon 9/3 forward error correction, intended to be sent as 9 pieces where 3 uncorrupted parts allows assembly of the message
import (
	"encoding/binary"
	"hash/crc32"
	"log"

	"github.com/vivint/infectious"
)

var (
	rsTotal    = 9
	rsRequired = 3
	rsFEC      = func() *infectious.FEC {

		fec, err := infectious.NewFEC(3, 9)

		if err != nil {

			panic(err)
		}
		return fec
	}()
)

// padData appends a 2 byte length prefix, and pads to a multiple of rsTotal. An empty slice will be returned if the total length is greater than maxMessageSize.
func padData(
	data []byte) (out []byte) {

	dataLen := len(data)
	prefixBytes := make([]byte, 2)
	binary.LittleEndian.PutUint16(prefixBytes, uint16(dataLen))
	data = append(prefixBytes, data...)
	dataLen = len(data)

	if dataLen > maxMessageSize {

		return []byte{}
	}
	chunkLen := (dataLen) / rsTotal
	chunkMod := (dataLen) % rsTotal

	if chunkMod != 0 {

		chunkLen++
	}
	padLen := rsTotal*chunkLen - dataLen
	out = append(data, make([]byte, padLen)...)
	return
}
func rsEncode(
	data []byte) (chunks [][]byte) {

	// First we must pad the data
	data = padData(data)
	shares := make([]infectious.Share, rsTotal)
	output := func(s infectious.Share) {

		shares[s.Number] = s.DeepCopy()
	}
	err := rsFEC.Encode(data, output)

	if err != nil {

		panic(err)
	}

	for i := range shares {

		// Append the chunk number to the front of the chunk
		chunk := append([]byte{byte(shares[i].Number)}, shares[i].Data...)

		// Checksum includes chunk number byte so we know if its checksum is incorrect so could the chunk number be
		checksum := crc32.Checksum(chunk, crc32.MakeTable(crc32.Castagnoli))
		checkbytes := make([]byte, 4)
		binary.LittleEndian.PutUint32(checkbytes, checksum)
		chunk = append(chunk, checkbytes...)
		chunks = append(chunks, chunk)
	}
	return
}
func rsDecode(
	chunks [][]byte) (data []byte, err error) {

	defer func() {

		if r := recover(); r != nil {

			log.Print("Recovered in f", r)
		}
	}()
	var shares []infectious.Share

	for i := range chunks {

		bodyLen := len(chunks[i])
		body := chunks[i][:bodyLen-4]
		share := infectious.Share{
			Number: int(body[0]),
			Data:   body[1:],
		}
		shares = append(shares, share)
	}
	data, err = rsFEC.Decode(nil, shares)
	return
}
