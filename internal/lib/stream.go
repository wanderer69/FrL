package frl

import (
	"bufio"
	"fmt"
	"io"
	"net"
	"os"
	"regexp"
)

// поток - тип который включает в себя все способы получения потока информации
type ConnectorFuncs interface {
	ControlSet(cmd string, settings map[string]string) error
	ControlGet(cmd string) (map[string]string, error)
	Open(file string) error
	Read() (int, []byte, error)
	Write(data []byte) error
	Close() error
}

type Connector struct {
	ConnectorFuncs
	Type string
}

func (connector Connector) ControlSet(cmd string, settings map[string]string) error {
	return connector.ConnectorFuncs.ControlSet(cmd, settings)
}

func (connector Connector) ControlGet(cmd string) (map[string]string, error) {
	return connector.ConnectorFuncs.ControlGet(cmd)
}

func (connector Connector) Open(file string) error {
	return connector.ConnectorFuncs.Open(file)
}

func (connector Connector) Read() (int, []byte, error) {
	return connector.ConnectorFuncs.Read()
}

func (connector Connector) Write(data []byte) error {
	return connector.ConnectorFuncs.Write(data)
}

func (connector Connector) Close() error {
	return connector.ConnectorFuncs.Close()
}

type Stream struct {
	URI        string // источник либо приемник
	SourceType string // тип источника либо приемника
	Connector  *Connector
	tags       map[string]string
}

// <схема> <:> [<//>] [[<пользователь> <@>]] [<хост>]  [[<:> <порт>]] <путь> [[<?> <запрос>]] [[<#> <фрагмент>]]
func ParseURI(uri string) (map[string]string, error) {
	tags := make(map[string]string)
	re := regexp.MustCompile(`^((?P<schema>[^:/?#]+):)?(//(?P<source>[^/?#]*))?(?P<path>[^?#]*)(\?(?P<query>[^#]*))?(?P<fragment>#(.*))?`)
	matches := re.FindStringSubmatch(uri)
	ss := re.SubexpNames()
	for i := range ss {
		lastIndex := re.SubexpIndex(ss[i])
		if lastIndex >= 0 {
			tags[ss[i]] = matches[lastIndex]
		}
	}
	_, ok := tags["schema"]
	if !ok {
		return nil, fmt.Errorf("schema not found")
	}

	return tags, nil
}

func (s *Stream) ControlSet(cmd string, settings map[string]string) error {
	err := s.Connector.ControlSet(cmd, settings)
	return err
}

func (s *Stream) ControlGet(cmd string) (map[string]string, error) {
	ms, err := s.Connector.ControlGet(cmd)
	return ms, err
}

func (s *Stream) Open(file string) error {
	err := s.Connector.Open(file)
	return err
}
func (s *Stream) Read() (int, []byte, error) {
	cnt, data, err := s.Connector.Read()
	return cnt, data, err
}

func (s *Stream) Write(data []byte) error {
	err := s.Connector.Write(data)
	return err
}

func (s *Stream) Close() error {
	err := s.Connector.Close()
	return err
}

func NewStream(uri string) (*Stream, error) {
	dct, err := ParseURI(uri)
	if err != nil {
		return nil, err
	}
	s := Stream{}
	schema, ok := dct["schema"]
	if !ok {
		return nil, fmt.Errorf("NewStream: schema not found - bad uri %v", uri)
	}
	switch schema {
	case "file":
		s.SourceType = "file"
		s.Connector = &Connector{ConnectorFuncs: &FileConnector{stream: &s}}
	case "tcp":
		s.SourceType = "tcp"
		s.Connector = &Connector{ConnectorFuncs: &TCPConnector{stream: &s}}
	case "http":
	case "serial":
	case "tty":
	}
	s.tags = dct
	return &s, nil
}

// File interface
type FileConnector struct {
	mode      string
	file      *os.File
	file_name string
	rd        *bufio.Reader
	stream    *Stream
}

func (fc *FileConnector) ControlSet(cmd string, settings map[string]string) error {
	switch cmd {
	case "set":
	}
	for k, v := range settings {
		switch k {
		case "mode":
			fc.mode = v
		}
	}
	return nil
}

func (fc *FileConnector) ControlGet(cmd string) (map[string]string, error) {
	settings := make(map[string]string)
	return settings, nil
}

func (fc *FileConnector) Open(_ string) error {
	source := fc.stream.tags["source"]
	path := fc.stream.tags["path"]
	file_name := source + path
	switch fc.mode {
	case "bytes_packet":
		file, err := os.Open(file_name) // For read access.
		if err != nil {
			file, err = os.Create(file_name)
			if err != nil {
				return err
			}
		}
		fc.file = file
	case "full":
		fc.file_name = file_name
	case "by_lines":
		fc.file_name = file_name
		// file, err := os.Open(file_name)
		file, err := os.OpenFile(file_name, os.O_RDWR|os.O_CREATE, 0666)
		if err != nil {
			//file, err = os.Create(file_name)
			//if err != nil {
			return err
			//}
		}
		fc.file = file
		fc.rd = bufio.NewReader(fc.file)
	}
	return nil
}

func (fc *FileConnector) Read() (int, []byte, error) {
	switch fc.mode {
	case "bytes_packet":
		data := make([]byte, 100)
		count, err := fc.file.Read(data)
		if err != nil {
			return -1, nil, err
		}
		return count, data, nil
	case "full":
		data, err := os.ReadFile(fc.file_name)
		if err != nil {
			return -1, nil, err
		}
		return len(data), data, nil
	case "by_lines":
		line, err := fc.rd.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				break
			}
			return -1, nil, fmt.Errorf("read file line error: %v", err)
		}
		return len(line), []byte(line), nil
	}
	return -1, nil, nil
}

func (fc *FileConnector) Write(data []byte) error {
	switch fc.mode {
	case "bytes_packet":
		_, err := fc.file.Write(data)
		if err != nil {
			return fmt.Errorf("write file line error: %v", err)
		}
		return nil
	case "full":
		err := os.WriteFile(fc.file_name, data, 0644)
		if err != nil {
			return fmt.Errorf("write file line error: %v", err)
		}
		return nil
	case "by_lines":
		ss := fmt.Sprintf("%v\n", string(data))
		_, err := fc.file.WriteString(ss)
		if err != nil {
			return fmt.Errorf("write file line error: %v", err)
		}
		return nil
	}
	return nil
}

func (fc *FileConnector) Close() error {
	switch fc.mode {
	case "bytes_packet":
		fc.file.Close()
	case "full":
	case "by_lines":
		fc.file.Close()
	}
	return nil
}

// TCP interface
type TCPConnector struct {
	mode    string
	conn    *net.TCPConn
	tcpAddr *net.TCPAddr
	//	rd      *bufio.Reader
	stream *Stream
}

func (tcp_c *TCPConnector) ControlSet(cmd string, settings map[string]string) error {
	switch cmd {
	case "set":
	}
	for k, v := range settings {
		switch k {
		case "mode":
			tcp_c.mode = v
		}
	}
	return nil
}

func (tcp_c *TCPConnector) ControlGet(cmd string) (map[string]string, error) {
	settings := make(map[string]string)
	return settings, nil
}

func (tcp_c *TCPConnector) Open(source string) error {
	src := tcp_c.stream.tags["source"]
	tcpAddr, err := net.ResolveTCPAddr("tcp", src)
	if err != nil {
		return fmt.Errorf("ResolveTCPAddr failed: %v", err)
	}
	tcp_c.tcpAddr = tcpAddr

	conn, err := net.DialTCP("tcp", nil, tcpAddr)
	if err != nil {
		return fmt.Errorf("DialTCP failed: %v", err)
	}
	tcp_c.conn = conn
	return nil
}

func (tcp_c *TCPConnector) Read() (int, []byte, error) {
	data := make([]byte, 1024)
	count, err := tcp_c.conn.Read(data)
	if err != nil {
		return -1, nil, fmt.Errorf("Read from server failed: %v", err)
	}
	return count, data, nil
}

func (tcp_c *TCPConnector) Write(data []byte) error {
	_, err := tcp_c.conn.Write(data)
	if err != nil {
		return fmt.Errorf("Write to server failed: %v", err)
	}
	return nil
}

func (tcp_c *TCPConnector) Close() error {
	tcp_c.conn.Close()
	return nil
}
