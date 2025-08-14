package secenv

import (
	"encoding/json"

	"github.com/99designs/keyring"
)

// func OpenEnv(name string)

// type EnvConfig struct {
// 	DefaultEnv   string               `json:"defaultEnv" yaml:"defaultEnv"`
// 	Environments map[string]SecureEnv `json:"environments" yaml:"environments"`
// }

// type SecureEnv struct {
// 	Name    string            `json:"name" yaml:"name"`
// 	Desc    string            `json:"desc" yaml:"desc"`
// 	Control string            `json:"control" yaml:"control"`
// 	Vars    map[string][]byte `json:"vars" yaml:"vars"`
// 	Files   map[string][]byte `json:"files" yaml:"files"`
// }

func (s *SecureEnvironment) Init() {
	if s.config.ServiceName != "" && s.ring != nil {
		return
	}

	s.config = keyring.Config{
		ServiceName:             "VVVSecureEnv",
		KeychainName:            "VVVSecureEnv",
		LibSecretCollectionName: "VVVSecureEnv",
		WinCredPrefix:           "VVVSecureEnv",
		PassPrefix:              "VVVSecureEnv",
	}

	var err error
	s.ring, err = keyring.Open(s.config)
	if err != nil {
		panic(err)
	}

	keys, err := s.ring.Keys()
	if err != nil {
		panic(err)
	}

	if s.Name == "" {

	}

	for _, key := range keys {
		if key == s.Name {
			item, err := s.ring.Get(key)
			if err != nil {
				panic(err)
			}

			if item.Data == nil {
				continue
			}

			err = json.Unmarshal(item.Data, s)
			if err != nil {
				panic(err)
			}

			return
		}
	}

	s.Save()
}

func (s *SecureEnvironment) List() ([]string, error) {
	s.Init()

	return s.ring.Keys()
}

func (s *SecureEnvironment) Save() {
	data, err := json.Marshal(s)
	if err != nil {
		panic(err)
	}

	err = s.ring.Set(keyring.Item{
		Key:         s.Name,
		Data:        data,
		Label:       s.Name,
		Description: s.Desc,
	})

	if err != nil {
		panic(err)
	}
}

// ring, err := keyring.Open(keyring.Config{
// 	ServiceName:  "VVV",
// 	KeychainName: "SecureEnv",
// })
// if err != nil {
// 	panic(err)
// }

// err = ring.Set(keyring.Item{
// 	Key:         "testkey",
// 	Data:        []byte("testvalue"),
// 	Label:       "testlabel",
// 	Description: "testdescription",
// })

// if err != nil {
// 	panic(err)
// }

// item, err := ring.Get("testkey")
// if err != nil {
// 	panic(err)
// }
// println("Key:", item.Key)
// println("Label:", item.Label)
// println("Description:", item.Description)
// println("Data:", string(item.Data))

// // err := keyring.Set("vvv", "testkey", "testvalue")
// // if err != nil {
// // 	panic(err)
// // }
