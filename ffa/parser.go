// Copyright © 2020 Rodney Rodriguez
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package ffa

import (
	"strings"

	"github.com/asottile/dockerfile"
)

// ExtractAllCommandsFromDockerfile
func ExtractAllCommandsFromDockerfile(data string) ([]dockerfile.Command, error) {
	reader := strings.NewReader(data)
	commandList, err := dockerfile.ParseReader(reader)
	if err != nil {
		return nil, err
	}
	return commandList, nil
}

// ExtractRunCommandsFromDockerfile
func ExtractRunCommandsFromDockerfile(data string) ([]dockerfile.Command, error) {
	commandList, err := ExtractAllCommandsFromDockerfile(data)
	if err != nil {
		return nil, err
	}

	var commands []dockerfile.Command

	// Collect all commands in the Dockerfile
	for _, cmd := range commandList {
		if cmd.Cmd == "run" {
			commands = append(commands, cmd)
		}
	}
	return commands, nil
}
