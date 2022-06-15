// Copyright 2022 Redpanda Data, Inc.
//
// Use of this software is governed by the Business Source License
// included in the file https://github.com/redpanda-data/redpanda/blob/dev/licenses/bsl.md
//
// As of the Change Date specified in that file, in accordance with
// the Business Source License, use of this software will be governed
// by the Apache License, Version 2.0

package console

import (
	"flag"
	"fmt"
)

type Config struct {
	TopicDocumentation ConfigTopicDocumentation `yaml:"topicDocumentation" json:"topicDocumentation"`
}

func (c *Config) SetDefaults() {
	c.TopicDocumentation.SetDefaults()
}

func (c *Config) RegisterFlags(f *flag.FlagSet) {
	c.TopicDocumentation.RegisterFlags(f)
}

func (c *Config) Validate() error {
	err := c.TopicDocumentation.Validate()
	if err != nil {
		return fmt.Errorf("failed to validate topic documentation config: %w", err)
	}

	return nil
}
