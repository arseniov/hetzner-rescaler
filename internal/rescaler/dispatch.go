// Package rescaler owns the async-rescale Manager that the HTTP API
// delegates to. The synchronous RescaleOnce/eventKindFor helpers that
// lived here previously were removed when Deps.Rescaler was replaced by
// Deps.Manager (Plan A.14).
package rescaler
