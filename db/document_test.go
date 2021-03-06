/*
   Copyright © 2019, 2020 M.Watermann, 10247 Berlin, Germany
                  All rights reserved
               EMail : <support@mwat.de>
*/

package db

//lint:file-ignore ST1017 - I prefer Yoda conditions

import (
	"path/filepath"
	"reflect"
	"testing"
)

func TestTDocument_Cover(t *testing.T) {
	SetCalibreLibraryPath("/var/opt/Calibre/")
	d1 := TDocument{
		ID:   7628,
		path: "/Spiegel/Der Spiegel (2019-06-01) 23_2019 (7628)",
	}
	w1 := "/cover/7628/cover.gif"
	d2 := TDocument{
		ID:   6730,
		path: "/John Scalzi/Zoe's Tale (6730)",
	}
	w2 := "/cover/6730/cover.gif"
	d3 := TDocument{
		ID:   4793,
		path: "/Gail Carriger/Soulless [1] (4793)",
	}
	w3 := "/cover/4793/cover.gif"
	tests := []struct {
		name   string
		fields TDocument
		want   string
	}{
		// TODO: Add test cases.
		{" 1", d1, w1},
		{" 2", d2, w2},
		{" 2", d3, w3},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			doc := &tt.fields
			if got := doc.Cover(); got != tt.want {
				t.Errorf("TDocument.Cover() = '%s',\nwant '%s'", got, tt.want)
			}
		})
	}
} // TestTDocument_Cover()

func TestTDocument_coverAbs(t *testing.T) {
	SetCalibreLibraryPath("/var/opt/Calibre/")
	d1 := TDocument{
		ID:   7628,
		path: "Spiegel/Der Spiegel (2019-06-01) 23_2019 (7628)",
	}
	w1 := d1.path + "/cover.jpg"

	d2 := TDocument{
		ID:   6730,
		path: "John Scalzi/Zoe's Tale (6730)",
	}
	w2 := d2.path + "/cover.jpg"
	w3 := filepath.Join(dbCalibreLibraryPath, w1)
	w4 := filepath.Join(dbCalibreLibraryPath, w2)
	d5 := TDocument{
		ID:   4793,
		path: "Gail Carriger/Soulless [1] (4793)",
	}
	w5 := d5.path + "/cover.jpg"
	type args struct {
		aRelative bool
	}
	w6 := filepath.Join(dbCalibreLibraryPath, w5)
	tests := []struct {
		name    string
		fields  TDocument
		args    args
		want    string
		wantErr bool
	}{
		// TODO: Add test cases.
		{" 1", d1, args{true}, w1, false},
		{" 2", d2, args{true}, w2, false},
		{" 3", d1, args{false}, w3, false},
		{" 4", d2, args{false}, w4, false},
		{" 5", d5, args{true}, w5, false},
		{" 6", d5, args{false}, w6, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			doc := &tt.fields
			got, err := doc.CoverAbs(tt.args.aRelative)
			if (err != nil) != tt.wantErr {
				t.Errorf("TDocument.coverAbs() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("TDocument.coverAbs() = '%s',\nwant '%s'", got, tt.want)
			}
		})
	}
} // TestTDocument_coverAbs()

func TestTDocument_Filename(t *testing.T) {
	SetCalibreLibraryPath("/var/opt/Calibre/")
	d1 := TDocument{
		formats: &tFormatList{
			TEntity{
				Name: "AZW3",
			},
			TEntity{
				Name: "EPUB",
			},
			TEntity{
				Name: "PDF",
			},
		},
		path: "John Scalzi/Zoe's Tale (6730)",
	}
	w1 := filepath.Join(d1.path, "Zoe's Tale - John Scalzi.azw3")
	type args struct {
		aFormat string
	}
	tests := []struct {
		name   string
		fields TDocument
		args   args
		want   string
	}{
		// TODO: Add test cases.
		{" 1", d1, args{"azw3"}, w1},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			doc := &tt.fields
			if got := doc.Filename(tt.args.aFormat); got != tt.want {
				t.Errorf("TDocument.Filename() = '%s',\nwant '%s'", got, tt.want)
			}
		})
	}
} // TestTDocument_Filename

func TestTDocument_filenames(t *testing.T) {
	SetCalibreLibraryPath("/var/opt/Calibre/")
	d1 := TDocument{
		formats: &tFormatList{
			TEntity{
				Name: "PDF",
			},
		},
		path: "Spiegel/Der Spiegel (2019-06-01) 23_2019 (7628)",
	}
	w1 := &tPathMap{
		"PDF": "/var/opt/Calibre/Spiegel/Der Spiegel (2019-06-01) 23_2019 (7628)/Der Spiegel (2019-06-01) 23_2019 - Spiegel.pdf",
	}
	d2 := TDocument{
		formats: &tFormatList{
			TEntity{
				Name: "AZW3",
			},
			TEntity{
				Name: "EPUB",
			},
			TEntity{
				Name: "PDF",
			},
		},
		path: "John Scalzi/Zoe's Tale (6730)",
	}
	w2 := &tPathMap{
		"AZW3": "/var/opt/Calibre/John Scalzi/Zoe's Tale (6730)/Zoe's Tale - John Scalzi.azw3",
		"EPUB": "/var/opt/Calibre/John Scalzi/Zoe's Tale (6730)/Zoe's Tale - John Scalzi.epub",
		"PDF":  "/var/opt/Calibre/John Scalzi/Zoe's Tale (6730)/Zoe's Tale - John Scalzi.pdf",
	}
	tests := []struct {
		name   string
		fields TDocument
		want   *tPathMap
	}{
		// TODO: Add test cases.
		{" 1", d1, w1},
		{" 2", d2, w2},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			doc := &tt.fields
			if got := doc.filenames(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("TDocument.Filenames() = %v,\nwant %v", got, tt.want)
			}
		})
	}
} // TestTDocument_filenames()

func TestTDocument_Files(t *testing.T) {
	SetCalibreLibraryPath("/var/opt/Calibre/")
	d1 := TDocument{
		ID: 1,
		formats: &tFormatList{
			TEntity{
				ID:   2,
				Name: "PDF",
			},
		},
		Title: "this is the document's title",
	}
	w1 := &TEntityList{
		TEntity{
			ID:   2,
			Name: "PDF",
			URL:  "/file/1/PDF/this_is_the_document%27s_title.pdf",
		},
	}
	d2 := TDocument{
		authors: &tAuthorList{
			TEntity{
				ID:   3,
				Name: "Spiegel",
			},
		},
		ID: 7628,
		formats: &tFormatList{
			TEntity{
				ID:   2,
				Name: "PDF",
			},
		},
		Title: "Der Spiegel (2019-06-01) 23/2019",
	}
	w2 := &TEntityList{
		TEntity{
			ID:   2,
			Name: `PDF`,
			URL:  `/file/7628/PDF/Spiegel_-_Der_Spiegel_%282019-06-01%29_23-2019.pdf`,
		},
	}
	tests := []struct {
		name   string
		fields TDocument
		want   *TEntityList
	}{
		// TODO: Add test cases.
		{" 1", d1, w1},
		{" 2", d2, w2},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			doc := &tt.fields
			if got := doc.Files(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("TDocument.Files() = %v,\nwant %v", got, tt.want)
			}
		})
	}
} // TestTDocument_Files()
