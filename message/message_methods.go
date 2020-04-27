package message

import "fmt"

func (c *AppConnectionDeviceMessage) MakeHostAddress() string {
	return fmt.Sprintf("%s:%s", c.IP, c.Port)
}