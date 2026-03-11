package client

import (
	"fmt"
	"html/template"
	"net/url"
	"strings"
)

type urlOptions struct {
	offset      *uint
	limit       *uint
	pathParams  map[string]string
	queryParams map[string]string
}

type URLOption func(*urlOptions)

func WithOffset(offset uint) URLOption {
	return func(o *urlOptions) {
		o.offset = &offset
	}
}

func WithLimit(limit uint) URLOption {
	return func(o *urlOptions) {
		o.limit = &limit
	}
}

func WithPathParams(pathParams map[string]string) URLOption {
	return func(o *urlOptions) {
		o.pathParams = pathParams
	}
}

func WithQueryParams(queryParams map[string]string) URLOption {
	return func(o *urlOptions) {
		if o.queryParams == nil {
			o.queryParams = make(map[string]string)
		}
		for k, v := range queryParams {
			o.queryParams[k] = v
		}
	}
}

// constructURL builds the full URL for an API request.
func (c *FusionClient) constructURL(basePath, path string, opts ...URLOption) (*url.URL, error) {
	options := &urlOptions{}
	for _, opt := range opts {
		opt(options)
	}

	if c.instanceURL == nil {
		return nil, fmt.Errorf("instance URL is not configured")
	}
	u := *c.instanceURL
	u.Path = basePath

	// Add the path with template parameter replacement
	if path != "" {
		tmpl, err := template.New("path").Parse(path)
		if err != nil {
			return nil, fmt.Errorf("failed to parse path template: %w", err)
		}
		var buf strings.Builder
		pathData := options.pathParams
		if pathData == nil {
			pathData = make(map[string]string)
		}
		if err := tmpl.Execute(&buf, pathData); err != nil {
			return nil, fmt.Errorf("failed to execute path template: %w", err)
		}
		u.Path += buf.String()
	}

	// Add query parameters
	q := u.Query()

	// Oracle Fusion Cloud uses offset-based pagination
	if options.offset != nil {
		q.Set("offset", fmt.Sprintf("%d", *options.offset))
	}
	if options.limit != nil {
		q.Set("limit", fmt.Sprintf("%d", *options.limit))
	}

	if options.queryParams != nil {
		for k, v := range options.queryParams {
			q.Set(k, v)
		}
	}
	u.RawQuery = q.Encode()

	return &u, nil
}

// GetNextOffset returns the next offset for pagination.
// Returns -1 if there are no more pages.
func GetNextOffset(hasMore bool, currentOffset, limit uint) int {
	if !hasMore {
		return -1
	}
	return int(currentOffset + limit)
}
