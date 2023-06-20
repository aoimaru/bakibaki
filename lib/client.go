package lib

type Client struct {
	Root string
}

func (c *Client) GetIndexPath() string {
	return c.Root + "/index"
}
