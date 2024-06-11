package imapclient

import (
	"fmt"

	"github.com/emersion/go-imap/v2"
	"github.com/emersion/go-imap/v2/internal/imapwire"
)

// Id sends an ID command.
// ID command is introduced in RFC 2971.
// An example ID command:
// ID ("name" "go-imap" "version" "1.0" "os" "Linux" "os-version" "7.9.4" "vendor" "Yahoo")
func (c *Client) Id(idData *imap.IdData) *IdCommand {
	cmd := &IdCommand{}
	enc := c.beginCommand("ID", cmd)

	if idData == nil {
		enc.SP().NIL()
	} else {
		enc.SP().Special('(')
		isFirstKey := true
		if idData.Name != "" {
			addIdKeyValue(enc, isFirstKey, "name", idData.Name)
			isFirstKey = false
		}
		if idData.Version != "" {
			addIdKeyValue(enc, isFirstKey, "version", idData.Version)
			isFirstKey = false
		}
		if idData.Os != "" {
			addIdKeyValue(enc, isFirstKey, "os", idData.Os)
			isFirstKey = false
		}
		if idData.OsVersion != "" {
			addIdKeyValue(enc, isFirstKey, "os-version", idData.OsVersion)
			isFirstKey = false
		}
		if idData.Vendor != "" {
			addIdKeyValue(enc, isFirstKey, "vendor", idData.Vendor)
			isFirstKey = false
		}
		if idData.SupportUrl != "" {
			addIdKeyValue(enc, isFirstKey, "support-url", idData.SupportUrl)
			isFirstKey = false
		}
		if idData.Address != "" {
			addIdKeyValue(enc, isFirstKey, "address", idData.Address)
			isFirstKey = false
		}
		if idData.Date != "" {
			addIdKeyValue(enc, isFirstKey, "date", idData.Date)
			isFirstKey = false
		}
		if idData.Command != "" {
			addIdKeyValue(enc, isFirstKey, "command", idData.Command)
			isFirstKey = false
		}
		if idData.Arguments != "" {
			addIdKeyValue(enc, isFirstKey, "arguments", idData.Arguments)
			isFirstKey = false
		}
		if idData.Environment != "" {
			addIdKeyValue(enc, isFirstKey, "environment", idData.Environment)
			isFirstKey = false
		}

		enc.Special(')')
	}
	enc.end()

	return cmd
}

func addIdKeyValue(enc *commandEncoder, isFirstKey bool, key, value string) {
	if !isFirstKey {
		enc.SP().Quoted(key).SP().Quoted(value)
	} else {
		enc.Quoted(key).SP().Quoted(value)
	}
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
			return nil
		}
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
			// Ignore unknown key
			// Yahoo server sends "host" and "remote-host" keys
			// which are not defined in RFC 2971
		}
		currKey = ""

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

func (r *IdCommand) Wait() (*imap.IdData, error) {
	return &r.data, r.cmd.Wait()
}
