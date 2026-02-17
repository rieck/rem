package rem

import "embed"

// EmbeddedSkills contains the agent skills files (SKILL.md + references/)
// baked into the binary at build time.
//
//go:embed skills/rem-cli
var EmbeddedSkills embed.FS
