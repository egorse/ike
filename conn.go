package ike

import (
	"io"
	"net"
	"os"
	"runtime"
	"syscall"

	"github.com/msgboxio/log"
	"github.com/pkg/errors"
	"golang.org/x/net/ipv4"
	"golang.org/x/net/ipv6"
)

type Conn interface {
	ReadPacket() (b []byte, remoteAddr net.Addr, localIP net.IP, err error)
	WritePacket(reply []byte, remoteAddr net.Addr) error
	LocalAddr() net.Addr
	Close() error
}

type pconnV4 ipv4.PacketConn

func (c *pconnV4) Close() error {
	return c.Conn.Close()
}

func (c *pconnV4) LocalAddr() net.Addr {
	return c.Conn.LocalAddr()
}

type pconnV6 ipv6.PacketConn

func (c *pconnV6) Close() error {
	return c.Conn.Close()
}

func (c *pconnV6) LocalAddr() net.Addr {
	return c.Conn.LocalAddr()
}

var ErrorUdpOnly = errors.New("only udp is supported for now")

// On Mac, dual stack bind for v4 addresses does not give us source IP addresses
func checkV4onX(address string) (bool, error) {
	if runtime.GOOS != "darwin" {
		return false, nil
	}
	v4Only := false
	addr, err := net.ResolveUDPAddr("udp4", address)
	if err != nil {
		return v4Only, err
	}
	if ip6 := addr.IP.To16(); ip6 == nil {
		v4Only = true
	}
	return v4Only, nil
}

func Listen(network, address string) (Conn, error) {
	isV4, err := checkV4onX(address)
	if err != nil {
		return nil, err
	}
	if isV4 {
		return listenUDP4(address)
	}
	switch network {
	case "udp4":
		return listenUDP4(address)
	case "udp6", "udp":
		return listenUDP6(address)
	}
	return nil, ErrorUdpOnly
}

func listenUDP4(localString string) (p4 *pconnV4, err error) {
	udp, err := net.ListenPacket("udp4", localString)
	if err != nil {
		return nil, errors.Wrap(err, "lsten")
	}
	p := ipv4.NewPacketConn(udp)
	// the interface could be set to any(0.0.0.0)
	// we need the exact address the packet came on
	cf := ipv4.FlagTTL | ipv4.FlagSrc | ipv4.FlagDst | ipv4.FlagInterface
	if err := p.SetControlMessage(cf, true); err != nil {
		if protocolNotSupported(err) {
			log.Warningf("udp source address detection not supported on %s", runtime.GOOS)
		} else {
			p.Close()
			return nil, err
		}
	}
	return (*pconnV4)(p), nil
}

func listenUDP6(localString string) (p6 *pconnV6, err error) {
	udp, err := net.ListenPacket("udp", localString)
	if err != nil {
		return nil, errors.Wrap(err, "lsten")
	}
	p := ipv6.NewPacketConn(udp)
	// the interface could be set to any(0.0.0.0)
	// we need the exact address the packet came on
	cf := ipv6.FlagSrc | ipv6.FlagDst | ipv6.FlagInterface
	if err := p.SetControlMessage(cf, true); err != nil {
		if protocolNotSupported(err) {
			log.Warningf("udp source address detection not supported on %s", runtime.GOOS)
		} else {
			p.Close()
			return nil, err
		}
	}
	return (*pconnV6)(p), nil
}

func (p *pconnV4) ReadPacket() (b []byte, remoteAddr net.Addr, localIP net.IP, err error) {
	b = make([]byte, 3000) // section 2
	n, cm, remoteAddr, err := p.ReadFrom(b)
	if err == nil {
		b = b[:n]
		localIP = cm.Dst
	}
	log.V(1).Infof("%d from %v", n, remoteAddr)
	return
}

func (p *pconnV6) ReadPacket() (b []byte, remoteAddr net.Addr, localIP net.IP, err error) {
	b = make([]byte, 3000) // section 2
	n, cm, remoteAddr, err := p.ReadFrom(b)
	if err == nil {
		b = b[:n]
		if cm != nil { // nil on mac
			localIP = cm.Dst
		}
	}
	log.V(1).Infof("%d from %v", n, remoteAddr)
	return
}

func (p *pconnV6) WritePacket(reply []byte, remoteAddr net.Addr) error {
	n, err := p.WriteTo(reply, nil, remoteAddr)
	if err != nil {
		return err
	} else if n != len(reply) {
		return io.ErrShortWrite
	}
	log.V(1).Infof("%d to %v", n, remoteAddr)
	return nil
}

func (p *pconnV4) WritePacket(reply []byte, remoteAddr net.Addr) error {
	n, err := p.WriteTo(reply, nil, remoteAddr)
	if err != nil {
		return err
	} else if n != len(reply) {
		return io.ErrShortWrite
	}
	log.V(1).Infof("%d to %v", n, remoteAddr)
	return nil
}

// ReadMessage reads an IKE message from connection
// Connection errors are returned, protocol errors are simply logged
func ReadMessage(conn Conn) (*Message, error) {
	var buf []byte
	for {
		b, remoteAddr, localIP, err := conn.ReadPacket()
		if err != nil {
			return nil, err
		}
		if buf != nil {
			b = append(buf, b...)
			buf = nil
		}
		msg, err := DecodeMessage(b)
		if err == io.ErrShortBuffer {
			buf = b
			continue
		}
		if err != nil {
			log.Error(err)
			continue
		}
		port := conn.LocalAddr().(*net.UDPAddr).Port
		msg.LocalAddr = &net.UDPAddr{
			IP:   localIP,
			Port: port,
		}
		msg.RemoteAddr = remoteAddr
		return msg, nil
	}
}

// InnerConn returns the conn buried within the conn used here
func InnerConn(p Conn) net.Conn {
	if p4Conn, ok := p.(*pconnV4); ok {
		return p4Conn.Conn
	} else if p6Conn, ok := p.(*pconnV6); ok {
		return p6Conn.Conn
	}
	return nil
}

// copied from golang.org/x/net/internal/nettest
func protocolNotSupported(err error) bool {
	switch err := err.(type) {
	case syscall.Errno:
		switch err {
		case syscall.EPROTONOSUPPORT, syscall.ENOPROTOOPT:
			return true
		}
	case *os.SyscallError:
		switch err := err.Err.(type) {
		case syscall.Errno:
			switch err {
			case syscall.EPROTONOSUPPORT, syscall.ENOPROTOOPT:
				return true
			}
		}
	}
	return false
}
