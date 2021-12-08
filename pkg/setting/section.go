package setting

import "time"

type ServerSettingS struct {
	RunMode      string
	HttpPort     string
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
	Addr         string
}

type BitCaskSettingS struct {
	MountDir string
}

type LoggerSettingS struct {
	FileName        string
	LogSavePath     string
	LogFileExt      string
	MaxPageSize     int
	DefaultPageSize int
}

func (s *Setting) ReadSection(k string, v interface{}) error {
	err := s.vp.UnmarshalKey(k, v)
	if err != nil {
		return err
	}
	return nil
}
