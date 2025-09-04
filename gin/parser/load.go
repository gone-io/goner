package parser

import "github.com/gone-io/gone/v2"

func Load(loader gone.Loader) error {
	loader.
		MustLoad(&bodyNameParser{}).
		MustLoad(&cookeNameParser{}).
		MustLoad(&headerNameParser{}).
		MustLoad(&paramNameParser{}).
		MustLoad(&queryNameParser{}).
		MustLoad(&ginContextTypeParser{}).
		MustLoad(&httpRequestTypeParser{}).
		MustLoad(&httpHeaderTypeParser{}).
		MustLoad(&urlTypeParser{}).
		MustLoad(&responseTypeParser{}).
		MustLoad(&httpResponseTypeParser{})
	return nil
}
