package skillgetmanager

import (
	"context"
	"fmt"
	"net/url"
	"strconv"
)

// SearchSkillsOptions mirrors the TS SearchSkillsOptions.
type SearchSkillsOptions struct {
	Query  string
	Limit  int
	Offset int
}

// SearchSkills lists or searches skills on the registry.
func SearchSkills(ctx context.Context, opts SearchSkillsOptions) (*ListSkillsResponse, error) {
	limit := opts.Limit
	if limit == 0 {
		limit = 20
	}
	offset := opts.Offset
	q := url.Values{}
	q.Set("limit", strconv.Itoa(limit))
	q.Set("offset", strconv.Itoa(offset))
	if opts.Query != "" {
		q.Set("q", opts.Query)
	}
	path := "/skills?" + q.Encode()
	var out ListSkillsResponse
	if err := FetchJSON(ctx, path, &out); err != nil {
		return nil, fmt.Errorf("search skills: %w", err)
	}
	return &out, nil
}
