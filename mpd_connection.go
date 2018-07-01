//
// get mpd log events, get corresponding metadata from mpd api, pass on metadata to elasticsearch
//

package main

import (
	"fmt"
	"time"
  // elastic "gopkg.in/olivere/elastic.v5"
  mpd "github.com/fhs/gompd/mpd"
)

type MpdClient struct {
  Conn  *mpd.Client
  proto string
  addr  string
}

// create new MPD client
func NewMpdClient(proto, addr string) (*MpdClient) {
  m := &MpdClient{
    proto: proto,
    addr:  addr,
  }

  return m
}

// need to be pinging MPD regularly so connection stays up
func (m *MpdClient) CreatePinger() (error) {
  for {
    select {
    case <- time.After(1000 * time.Millisecond):
      err := m.Conn.Ping()

      if err != nil {
        fmt.Printf("ping %s \n", err)
      } else {
        fmt.Printf("ping \n")
      }
    }
  }
}

// get or refresh mpd connection
func (m *MpdClient) MpdConn() (*MpdClient, error) {

  if m.Conn != nil {
    err := m.Conn.Ping()

    if err != nil {
      return m, nil
    }
  }

	for {
  	c, err := mpd.Dial(m.proto, m.addr)

	  if err != nil {
	    fmt.Println("cannot connect to MPD")
			time.Sleep(1000 * time.Millisecond)

		} else {
      go m.CreatePinger()

			m.Conn = c
		  return m, nil
		}
	}
}
