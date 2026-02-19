package bot

import (
	"encoding/binary"
	"fmt"
	"io"
	"log"
	"math/rand"
	"os"
	"time"

	"github.com/bwmarrin/discordgo"
)

var buffer = make([][]byte, 0)

func getVoiceClip() error {
	files, err := os.ReadDir("media")
	if err != nil {
		log.Printf("Error reading media directory: %v", err)
		return err
	}
	var names []string
	for _, file := range files {
		if !file.IsDir() {
			names = append(names, file.Name())
		}
	}
	file, err := os.Open(files[rand.Intn(len(files))].Name())
	if err != nil {
		log.Printf("Error reading media file: %v", err)
		return err
	}
	defer file.Close()

	var opuslen int16

	for {
		// Read opus frame length from dca file.
		err = binary.Read(file, binary.LittleEndian, &opuslen)

		// If this is the end of the file, just return.
		if err == io.EOF || err == io.ErrUnexpectedEOF {
			err := file.Close()
			if err != nil {
				return err
			}
			return nil
		}

		if err != nil {
			fmt.Println("Error reading from dca file :", err)
			return err
		}

		// Read encoded pcm from dca file.
		InBuf := make([]byte, opuslen)
		err = binary.Read(file, binary.LittleEndian, &InBuf)

		// Should not be any end of file errors
		if err != nil {
			fmt.Println("Error reading from dca file :", err)
			return err
		}

		// Append encoded pcm data to the buffer.
		buffer = append(buffer, InBuf)
	}

}

func voiceKing(s *discordgo.Session, guildID, channelID string) {
	voice, err := s.ChannelVoiceJoin(guildID, channelID, false, false)
	if err != nil {
		log.Printf("Error joining voice channel: %v", err)
		return
	}
	getVoiceClip()

	time.Sleep(250 * time.Millisecond) // Wait for the bot to join the voice channel
	voice.Speaking(true)

	for _, buff := range buffer {
		voice.OpusSend <- buff
	}

	voice.Speaking(false)
	time.Sleep(250 * time.Millisecond)
	voice.Disconnect()
}
