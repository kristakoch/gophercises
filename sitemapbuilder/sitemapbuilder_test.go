// Package sitemapbuilder builds a sitemap of a website. WIP. Needs XML encoding.
package main

import (
	"reflect"
	"testing"
)

func Test_crawlSite(t *testing.T) {
	type args struct {
		rootURL   string
		unvisited map[string]empty
		found     map[string]empty
	}
	tests := []struct {
		name    string
		args    args
		want    []Page
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := crawlSite(tt.args.rootURL, tt.args.unvisited, tt.args.found)
			if (err != nil) != tt.wantErr {
				t.Errorf("crawlSite() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("crawlSite() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_clean(t *testing.T) {
	type args struct {
		s *string
	}
	tests := []struct {
		name string
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			clean(tt.args.s)
		})
	}
}
