package parser

import "github.com/gone-io/gone/v2"

func Load(loader gone.Loader) error {
	loader.
		MustLoad(&bodyNameParser{}).
		MustLoad(&cookeNameParser{}).
		MustLoad(&headerNameParser{}).
		MustLoad(&paramNameParser{}).
		MustLoad(&queryNameParser{}).
		MustLoad(&contextTypeParser{}).
		MustLoad(&originContextTypeParser{}).
		MustLoad(&httpRequestTypeParser{}).
		MustLoad(&httpHeaderTypeParser{}).
		MustLoad(&urlTypeParser{}).
		MustLoad(&responseTypeParser{})
	return nil
}
