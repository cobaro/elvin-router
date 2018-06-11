// Copyright 2018 Cobaro Pty Ltd. All Rights Reserved.

// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
// SOFTWARE.

package main

import (
	"bytes"
	_ "github.com/cobaro/elvin/elvin"
	"io"
	"sync/atomic"
	"testing"
	"time"
)

var xID uint32 = 0

func XID() uint32 {
	return atomic.AddUint32(&xID, 1)
}

func TestMockup(t *testing.T) {
    config = TestConfig()

	// Create a dummy connection using pipes, assigning roles for
	// reader,writer, and closer
	//
	// client write ----> server read
	//                        |
	// client read ----> server write
	//
	// closer needs to be assigned from the return from io.Pipe()
	// return although we could use a CloseReader
	var client, server Connection
	//var client elvin.Client
	cr, sw := io.Pipe()
	sr, cw := io.Pipe()
	client.reader = cr
	client.writer = cw
	client.closer = cr
	server.reader = sr
	server.writer = sw
	server.closer = sr

	server.SetState(StateNew)
	server.writeChannel = make(chan *bytes.Buffer, 4) // Some queuing allowed to smooth things out
	server.writeTerminate = make(chan int)

	go server.readHandler()
	go server.writeHandler()

	// FIXME: At this point we need to think about the client library
	client.SetState(StateNew)
	client.writeChannel = make(chan *bytes.Buffer, 4) // Some queuing allowed to smooth things out
	client.writeTerminate = make(chan int)

	go client.readHandler()  // Bogus
	go client.writeHandler() // Bogus

	// if err := client.Connect(); err != nil {
		// t.Errorf("Connect failed: %v", err)
    // }

	// And bail for now
	time.Sleep(1000 * 1000 * 100 * 1)

}
