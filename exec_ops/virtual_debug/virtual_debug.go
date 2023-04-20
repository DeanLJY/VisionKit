package virtual_debug

// An op for debugging that applies an identity function.
// It wraps items in the input datasets under a virtual provider that removes
// the filename reference. So it's useful for testing to make sure that all ops
// are properly handling cases where Item.Fname is not available.

import (
	"github.com/skyhookml/skyhookml/skyhook"
	"github.com/skyhookml/skyhookml/exec_ops"

	"fmt"
	urllib "net/url"
)

func init() {
	skyhook.ItemProviders["virtual_debug"] = skyhook.VirtualProvider(func(item skyho