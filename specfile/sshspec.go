/*
Copyright © 2020 Joseph Saylor <doug@saylorsolutions.com>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package specfile

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"golang.org/x/crypto/ssh"
)

func noOpBanner(message string) error { return nil }

func publicKey(path string) (ssh.AuthMethod, error) {
	key, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}
	signer, err := ssh.ParsePrivateKey(key)
	if err != nil {
		return nil, err
	}
	return ssh.PublicKeys(signer), nil
}

type ClientFilePair struct {
	Client  *ssh.Client
	HostTag string
	File    string
}

// setupClients validates the spec data and sets up ClientFilePair instances.
func setupClients(specData *SpecData) ([]*ClientFilePair, error) {
	var err error
	clientPairs := make([]*ClientFilePair, len(specData.Hosts))
	err = specData.Validate()
	if err != nil {
		return nil, fmt.Errorf("Invalid spec data: %v", err)
	}
	i := 0
	for k, v := range specData.Hosts {
		authMethod, err := publicKey(specData.Keys[k].Path)
		if err != nil {
			return nil, fmt.Errorf("Failed to load key from %s: %v", specData.Keys[k].Path, err)
		}
		config := &ssh.ClientConfig{
			User: v.Username,
			Auth: []ssh.AuthMethod{
				authMethod,
			},
			BannerCallback:  noOpBanner,
			HostKeyCallback: ssh.InsecureIgnoreHostKey(),
		}
		config.SetDefaults()
		hostPort := fmt.Sprintf("%s:%d", v.Hostname, v.Port)
		client, err := ssh.Dial("tcp", hostPort, config)
		if err != nil {
			return nil, fmt.Errorf("Failed to connect to %s: %v", hostPort, err)
		}

		clientPairs[i] = &ClientFilePair{client, k, v.File}
		i++
	}
	return clientPairs, nil
}

type TailChannelWriter struct {
	prefix string
	ch     chan<- string
}

func (t TailChannelWriter) Write(b []byte) (n int, err error) {
	t.ch <- fmt.Sprintf("[ %s ] %s", t.prefix, string(b))
	n = len(b)
	return
}

// TailSession represents
type TailSession struct {
	clientPair *ClientFilePair
	session    *ssh.Session
	closed     bool
	started    bool
}

// Closed returns whether the tail session has been previously closed. A closed tail session cannot be restarted.
func (s *TailSession) Closed() bool {
	return s.closed
}

// Started returns whether the tail session has already been started.
func (s *TailSession) Started() bool {
	return s.started
}

// Close stops the running tail session and disconnects the client.
func (s *TailSession) Close() (err error) {
	if !s.closed {
		s.closed = true
		sb := strings.Builder{}
		errorsOccurred := false
		e1 := s.session.Close()
		if e1 != nil {
			sb.WriteString(e1.Error())
			errorsOccurred = true
		}
		e2 := s.clientPair.Client.Close()
		if e2 != nil {
			sb.WriteString(e2.Error())
			errorsOccurred = true
		}
		if errorsOccurred {
			err = fmt.Errorf("Error(s) closing tail session: %s", sb.String())
		}
	}
	return
}

// Start the tail session using configured parameters
func (s *TailSession) start(ch chan<- string) error {
	if !s.closed {
		if s.started {
			return errors.New("Tail session is already started")
		}
		session, err := s.clientPair.Client.NewSession()
		if err != nil {
			return fmt.Errorf("Error establishing session: %v", err)
		}
		s.session = session
		session.Stdout = TailChannelWriter{s.clientPair.HostTag, ch}
		go func() {
			cmd := fmt.Sprintf("tail -n 0 -f %s", s.clientPair.File)
			err = session.Run(cmd)
			if err != nil {
				fmt.Printf("Error running '%s': %v", cmd, err)
			}
		}()
		s.started = true
	} else {
		return errors.New("Can't start a closed tail session")
	}
	return nil
}

// NewTailSession creates a new TailSession instance that is ready to be started.
func NewTailSession(client *ClientFilePair) (ts *TailSession, err error) {
	ts = &TailSession{client, nil, false, false}
	return
}

// ConsolidatedWriter receives messages from all of its tail session instances and writes them to its output stream.
type ConsolidatedWriter struct {
	ch       chan string
	sessions []*TailSession
	out      *os.File
	started  bool
	closed   bool
}

// NewConsolidatedWriter creates tail sessions that are ready to start and write to the provided writer.
func NewConsolidatedWriter(specData *SpecData, out *os.File) (*ConsolidatedWriter, error) {
	clientPairs, err := setupClients(specData)
	numHosts := len(specData.Hosts)
	var ch chan string = make(chan string, numHosts)
	var sessions []*TailSession = make([]*TailSession, numHosts)
	if err != nil {
		return nil, err
	}

	for i, pair := range clientPairs {
		ts, err := NewTailSession(pair)
		if err != nil {
			return nil, err
		}
		sessions[i] = ts
	}

	return &ConsolidatedWriter{ch, sessions, out, false, false}, nil
}

// Close closes all tail sessions as well as the connected clients.
func (c *ConsolidatedWriter) Close() error {
	for _, ts := range c.sessions {
		if ts.Started() && !ts.Closed() {
			ts.Close()
		}
	}
	return nil
}

// Start starts all tail sessions. In the event of an error, all already opened sessions are closed and an error is returned.
func (c *ConsolidatedWriter) Start() error {
	for _, ts := range c.sessions {
		if !ts.Started() && !ts.Closed() {
			err := ts.start(c.ch)
			if err != nil {
				fmt.Println("Failed to start consolidated writer. Closing sessions.")
				c.Close()
				return err
			}
		}
	}
	for line := range c.ch {
		c.out.WriteString(line)
	}
	return nil
}
