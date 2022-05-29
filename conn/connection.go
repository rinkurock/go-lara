package conn

import (
	c "app/config"
	"time"
)


func ServiceConnection(){
	ConnectSample(c.GetConfig().Services.Sample.Timeout * time.Millisecond)
}