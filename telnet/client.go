package telnet

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"net"
	"os"
	"sync"
	"time"
)

//Client represents telnet client
type Client struct {
	Network     string
	Address     string
	Timeout     time.Duration
	Input       io.Reader
	Output      io.Writer
	isConnected bool
}

//NewClient create new telnet client for specified network and address 
func NewClient(network string, address string, timeout time.Duration, input io.Reader, output io.Writer) *Client {
	if network == "" {
		network = "tcp"
	}

	if address == "" {
		address = "127.0.0.1:3302"
	}

	if input == nil {
		input = os.Stdin
	}

	if output == nil {
		output = os.Stdout
	}

	return &Client{Network: network, Address: address, Timeout: timeout, Input: input, Output: output}
}

//Connect establishes connection with remote host using the provided context.
//Transfers text from connection and write to output stream
//Transfers text from input stream and write to net connection
func (c *Client) Connect(ctx context.Context) error {
	if ctx == nil {
		ctx = context.Background()
	}

	dialer := &net.Dialer{}
	ctx, cancelFunc := context.WithTimeout(ctx, c.Timeout)
	connection, err := dialer.DialContext(ctx, c.Network, c.Address)
	if err != nil {
		cancelFunc()
		return fmt.Errorf("Could not connect to remote host: %w", err)
	}

	c.isConnected = true

	var outputErr error
	var inputErr error
	
	wg := sync.WaitGroup{}

	wg.Add(1)
	go func() {
		//Read from net connection and write to output stream
		outputErr = c.process(ctx, cancelFunc, connection, c.Output)
		wg.Done()
	}()

	wg.Add(1)
	go func() {
		//Read from input stream and write to net connection
		inputErr = c.process(ctx, cancelFunc, c.Input, connection)
		wg.Done()
	}()

	wg.Wait()
	c.isConnected = false

	cancelFunc()
	connection.Close()

	if inputErr != nil && outputErr != nil {
		return fmt.Errorf(
			"error has occurred while process input (%w); error has occurred while process output (%w)", 
			inputErr, 
			outputErr,
		)
	} else if inputErr != nil {
		return fmt.Errorf("error has occurred while process input (%w)", inputErr)
	} else if outputErr != nil {
		return fmt.Errorf("error has occurred while process output (%w)", outputErr)
	}

	return nil
}

//Process transfers text data from input stream to output stream.
func (c *Client) process(ctx context.Context, cancelFunc context.CancelFunc, input io.Reader, output io.Writer) error {
	scanner := bufio.NewScanner(input)

	messageChan := make(chan string)
	errorChan := make(chan error)

	defer close(messageChan)
	defer close(errorChan)

	var err error
loop:
	for {
		go c.scan(scanner, messageChan, errorChan)

		select {
		case <-ctx.Done():
			break loop
		case msg := <-messageChan:
			output.Write([]byte(fmt.Sprintf("%s\n", msg)))
		case err = <-errorChan:
			cancelFunc()
			break loop
		}
	}

	return err
}

//Scan read message to message channel.
func (c *Client) scan(scanner *bufio.Scanner, messageChan chan<- string, errorChan chan<- error) {
	ok := scanner.Scan()
	if c.isConnected {
		if ok {
			messageChan <- scanner.Text()
		} else {
			errorChan <- scanner.Err()
		}
	}
}