package imapclient

import (
	"fmt"

	"github.com/emersion/go-imap/v2"
	"github.com/emersion/go-imap/v2/internal/imapwire"
)

func (c *Client) Id(keysAndValues ...string) *IdCommand {
	if len(keysAndValues)%2 != 0 {
		panic("imapclient: the length of keys and values is odd")
	}

	cmd := &IdCommand{}
	enc := c.beginCommand("ID", cmd)

	if len(keysAndValues) == 0 {
		enc.SP().NIL()
	} else {
		enc.SP().Special('(')

		for i, keyOrValue := range keysAndValues {
			enc.Quoted(keyOrValue)

			if i != len(keysAndValues)-1 {
				enc.SP()
			}
		}

		enc.Special(')')
	}
	enc.end()

	return cmd
}

func (c *Client) handleId() error {
	data, err := c.readId(c.dec)
	if err != nil {
		return fmt.Errorf("in id: %v", err)
	}

	if cmd := findPendingCmdByType[*IdCommand](c); cmd != nil {
		cmd.data = *data
	}

	return nil
}

func (c *Client) readId(dec *imapwire.Decoder) (*imap.IdData, error) {
	var data = imap.IdData{}

	if !dec.ExpectSP() {
		return nil, dec.Err()
	}

	if dec.ExpectNIL() {
		return &data, nil
	}

	currKey := ""
	err := dec.ExpectList(func() error {
		var keyOrValue string
		if !dec.String(&keyOrValue) {
			return fmt.Errorf("in id key-val list: %v", dec.Err())
		}

		if currKey == "" {
			currKey = keyOrValue
		} else {
			switch currKey {
			case "name":
				data.Name = keyOrValue
			case "version":
				data.Version = keyOrValue
			case "os":
				data.Os = keyOrValue
			case "os-version":
				data.OsVersion = keyOrValue
			case "vendor":
				data.Vendor = keyOrValue
			case "support-url":
				data.SupportUrl = keyOrValue
			case "address":
				data.Address = keyOrValue
			case "date":
				data.Date = keyOrValue
			case "command":
				data.Command = keyOrValue
			case "arguments":
				data.Arguments = keyOrValue
			case "environment":
				data.Environment = keyOrValue
			default:
				return fmt.Errorf("unknown key: %v", currKey)
			}
			currKey = ""
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	return &data, nil
}

type IdCommand struct {
	cmd
	data imap.IdData
}

func (r *IdCommand) Wait() imap.IdData {
	return r.data
}
