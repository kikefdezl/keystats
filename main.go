package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/MarinX/keylogger"
)

const (
	LogEveryMinutes = 10
	StateFile       = "/var/lib/keystats/keystats.json"
)

type KeyLog map[string]int

// log logs a key press
func (kl KeyLog) log(key string) {
	kl[key] += 1
}

// save saves the KeyLog to the JSON file. It overwrites any existing file.
func (kl *KeyLog) save() error {
	jsonString, err := json.MarshalIndent(kl, "", "    ")
	if err != nil {
		return err
	}
	os.WriteFile(StateFile, jsonString, 0o644)
	return nil
}

// loadKeyLogFromFile loads the KeyLog from the cache JSON file if it exists,
// or returns a new empty KeyLog if not.
func loadKeyLog() KeyLog {
	keylog := make(KeyLog)

	s, err := os.ReadFile(StateFile)
	if err != nil {
		return keylog
	}
	json.Unmarshal([]byte(s), &keylog)
	fmt.Println("Loading keystats from file", keylog)
	return keylog
}

// listenToDevice returns a channel of strings corresponsing to the key names
// whenever they are pressed
func listenToDevice(dev string) (chan string, error) {
	k, err := keylogger.New(dev)
	if err != nil {
		fmt.Println("Error setting up keylogger")
		return nil, err
	}

	ch := make(chan string)

	go func() {
		fmt.Println("Listening to device:", dev)
		defer k.Close()
		defer close(ch)

		events := k.Read()
		for e := range events {
			if e.Type == keylogger.EvKey && e.KeyPress() {
				ch <- e.KeyString()
			}
		}
	}()
	return ch, nil
}

// spawnTimer returns an empty struct every `minutes`.
func spawnTimer(minutes int) chan struct{} {
	duration := time.Duration(minutes) * time.Minute
	ch := make(chan struct{})
	go func() {
		defer close(ch)
		for {
			time.Sleep(duration)
			ch <- struct{}{}
		}
	}()
	return ch
}

func ensureStateDir() error {
	dir := filepath.Dir(StateFile)
	return os.MkdirAll(dir, 0o700)
}

func main() {
	ensureStateDir()

	device := keylogger.FindAllKeyboardDevices()
	fmt.Println("Found devices:", device)

	if len(device) == 0 {
		fmt.Println("No keyboard devices found.")
		return
	}

	keys, err := listenToDevice(device[1])
	if err != nil {
		fmt.Println("Error setting up keylogger. Are you running as root?")
		return
	}

	timer := spawnTimer(LogEveryMinutes)

	log := loadKeyLog()
	for {
		select {
		case key := <-keys:
			fmt.Println(key)
			log.log(key)
		case <-timer:
			log.save()
			fmt.Println("Saved to:", StateFile)
		}
	}
}
