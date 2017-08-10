package ccruncher

import "github.com/cloudfoundry-incubator/candiedyaml"

func (e LogEntry) Render() ([]byte, error) {
	return candiedyaml.Marshal(e)
}

// func (c *CCLog) Render() ([]byte, error) {
//   var
// }
